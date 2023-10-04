package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"

	"github.com/varoOP/shinkro/internal/config"
	"github.com/varoOP/shinkro/internal/database"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/logger"
	"github.com/varoOP/shinkro/internal/notification"
	"github.com/varoOP/shinkro/internal/server"
)

const usage = `shinkro
Sync your Anime watch status in Plex to myanimelist.net!
Usage:
    shinkro --config <path to shinkro configuration>	Run shinkro
    shinkro malauth --config <path to shinkro configuration> Set up your MAL account for use with shinkro
    shinkro version	Print version info
    shinkro help	Show this help message
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
		db := database.NewDB(configPath, log)

		log.Info().Msg("Starting shinkro")
		log.Info().Msgf("Version: %s", version)
		log.Info().Msgf("Commit: %s", commit)
		log.Info().Msgf("Build date: %s", date)
		log.Info().Msgf("Log-level: %s", cfg.LogLevel)

		err, mapLoaded := domain.ChecklocalMaps(cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to load local custom mapping")
		}

		if mapLoaded {
			log.Info().Msg("Loaded local custom mapping")
		}

		db.CreateDB()
		db.UpdateAnime()

		c := cron.New(cron.WithLocation(time.UTC))
		c.AddFunc("0 1 * * MON", func() { db.UpdateAnime() })
		c.Start()

		n := notification.NewAppNotification(cfg.DiscordWebHookURL, log)
		go n.ListenforNotification()
		s := server.NewServer(cfg, n.Notification, db, log)
		go s.Start()

		sigchnl := make(chan os.Signal, 1)
		signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-sigchnl

		db.Close()
		log.Fatal().Msgf("Caught signal %v, Shutting Down", sig)

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
