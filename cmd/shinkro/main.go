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
	"github.com/varoOP/shinkro/internal/config"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/http"
	"github.com/varoOP/shinkro/internal/logger"
	"github.com/varoOP/shinkro/internal/plex"
	"github.com/varoOP/shinkro/internal/server"

	// "github.com/varoOP/shinkro/internal/malauth"
	"github.com/varoOP/shinkro/internal/notification"
)

const usage = `shinkro
Sync your Anime watch status in Plex to myanimelist.net!
Usage:
    shinkro --config <path to shinkro configuration directory> Run shinkro
    shinkro genkey                                             Generate an API key
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
		cfg := config.NewConfig(configPath).Config
		log := logger.NewLogger(configPath, cfg)
		db := database.NewDB(configPath, &log)

		log.Info().Msg("Starting shinkro")
		log.Info().Msgf("Version: %s", version)
		log.Info().Msgf("Commit: %s", commit)
		log.Info().Msgf("Build date: %s", date)
		log.Info().Msgf("Base URL: %s", cfg.BaseUrl)
		log.Info().Msgf("Log-level: %s", cfg.LogLevel)

		err, mapLoaded := domain.ChecklocalMaps(cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to load local custom mapping")
		}

		if mapLoaded {
			log.Info().Msg("Loaded local custom mapping")
		}

		err = db.Migrate()
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}

		var animeRepo = database.NewAnimeRepo(log, db)
		var animeService = anime.NewService(log, animeRepo)
		// c := cron.New(cron.WithLocation(time.UTC))
		// c.AddFunc("0 1 * * MON", func() {
		// 	// db.UpdateAnime()
		// 	// malauth.NewOauth2Client(context.Background(), db)
		// })
		// c.Start()

		n := notification.NewAppNotification(cfg.DiscordWebHookURL, &log)
		go n.ListenforNotification()
		// s := server.NewServer(cfg, n.Notification, db, log)
		// go s.Start()

		var plexRepo = database.NewPlexRepo(log, db)
		var plexService = plex.NewService(log, cfg, plexRepo, animeService)

		srv := server.NewServer(log, cfg, animeService)
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

	case "genkey":
		fmt.Fprintln(os.Stdout, config.GenApikey())

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
