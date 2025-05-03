package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"

	"github.com/varoOP/shinkro/internal/anime"
	"github.com/varoOP/shinkro/internal/animeupdate"
	"github.com/varoOP/shinkro/internal/api"
	"github.com/varoOP/shinkro/internal/auth"
	"github.com/varoOP/shinkro/internal/config"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/filesystem"
	"github.com/varoOP/shinkro/internal/http"
	"github.com/varoOP/shinkro/internal/logger"
	"github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/mapping"
	"github.com/varoOP/shinkro/internal/notification"
	"github.com/varoOP/shinkro/internal/plex"
	"github.com/varoOP/shinkro/internal/plexsettings"
	"github.com/varoOP/shinkro/internal/server"
	"github.com/varoOP/shinkro/internal/user"
)

const usage = `shinkro
Sync your Anime watch status in Plex to myanimelist.net!
Usage:
    shinkro --config <path to shinkro configuration directory> Run shinkro
    shinkro version                                            Print version info
    shinkro help                                               Show this help message
`

func init() {
	pflag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}
}

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	var configPath string

	d, err := homedir.Dir()
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: Unable to get home directory")
		os.Exit(1)
	}

	d = filepath.Join(d, ".config", "shinkro")
	pflag.StringVar(&configPath, "config", d, "path to configuration")
	pflag.Parse()
	configPath, err = homedir.Expand(configPath)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: Unable to expand configuration path")
		os.Exit(1)
	}

	switch cmd := pflag.Arg(0); cmd {
	case "":
		cfg := config.NewConfig(configPath, version)
		log := logger.NewLogger(cfg.Config).Logger
		go cfg.DynamicReload(log)
		db := database.NewDB(configPath, &log)

		log.Info().Msg("Starting shinkro")
		log.Info().Msgf("Version: %s", version)
		log.Info().Msgf("Commit: %s", commit)
		log.Info().Msgf("Build date: %s", date)
		log.Info().Msgf("Base URL: %s", cfg.Config.BaseUrl)
		log.Info().Msgf("Log-level: %s", cfg.Config.LogLevel)

		err = db.Migrate()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		var animeRepo = database.NewAnimeRepo(log, db)
		var animeUpdateRepo = database.NewAnimeUpdateRepo(log, db)
		var plexRepo = database.NewPlexRepo(log, db)
		var plexSettingsRepo = database.NewPlexSettingsRepo(log, db)
		var malauthRepo = database.NewMalAuthRepo(log, db)
		var userRepo = database.NewUserRepo(log, db)
		var apiRepo = database.NewAPIRepo(log, db)
		var mappingRepo = database.NewMappingRepo(log, db)
		var notificationRepo = database.NewNotificationRepo(log, db)

		var animeService = anime.NewService(log, animeRepo)
		var malauthService = malauth.NewService(cfg.Config, log, malauthRepo)
		var mapService = mapping.NewService(log, mappingRepo)
		var animeUpdateService = animeupdate.NewService(log, animeUpdateRepo, animeService, mapService, malauthService)
		var plexSettingsService = plexsettings.NewService(cfg.Config, log, plexSettingsRepo)
		var notificationService = notification.NewService(log, notificationRepo)
		var plexService = plex.NewService(log, plexSettingsService, plexRepo, animeService, mapService, malauthService, animeUpdateService, notificationService)
		var userService = user.NewService(userRepo, log)
		var authService = auth.NewService(log, userService)
		var apiService = api.NewService(log, apiRepo)
		var fsService = filesystem.NewService(cfg.Config, log)

		srv := server.NewServer(log, cfg.Config, animeService, mapService)
		if err := srv.Start(); err != nil {
			log.Fatal().Stack().Err(err).Msg("could not start server")
			return
		}

		errorChannel := make(chan error)

		go func() {
			httpServer := http.NewServer(
				log,
				cfg,
				db,
				version,
				commit,
				date,
				plexService,
				plexSettingsService,
				malauthService,
				apiService,
				authService,
				mapService,
				fsService,
				notificationService,
			)
			errorChannel <- httpServer.Open()
		}()

		sigchnl := make(chan os.Signal, 1)
		signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

		for sig := range sigchnl {
			log.Info().Msgf("received signal: %v, shutting down server.", sig)
			if err := db.Close(); err != nil {
				log.Error().Err(err).Msg("failed to close the database connection properly")
				os.Exit(1)
			}
			os.Exit(0)
		}

	case "version":
		fmt.Fprintln(flag.CommandLine.Output(), "Version:", version)
		fmt.Fprintln(flag.CommandLine.Output(), "Commit:", commit)
		fmt.Fprintln(flag.CommandLine.Output(), "Build Date:", date)

	default:
		pflag.Usage()
		if cmd != "help" {
			os.Exit(0)
		}
	}
}
