package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"

	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/domain"
	"github.com/varoOP/shinkuro/internal/logger"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/server"
)

const usage = `shinkuro
Sync your Anime watch status in Plex to myanimelist.net!
Usage:
    shinkuro --config <path to shinkuro configuration>	Run shinkuro
    shinkuro malauth --config <path to shinkuro configuration> Set up your MAL account for use with shinkuro
    shinkuro version	Print version info
    shinkuro help	Show this help message
`

func init() {
	pflag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), usage)
	}
}

func main() {
	var configPath string

	d, err := homedir.Dir()
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: unable to get home directory")
		os.Exit(1)
	}

	d = filepath.Join(d, ".config", "shinkuro")
	pflag.StringVar(&configPath, "config", d, "path to configuration")
	pflag.Parse()
	configPath, err = homedir.Expand(configPath)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "FATAL: unable to expand configuration path")
		os.Exit(1)
	}

	cfg := config.NewConfig(configPath).Config
	log := logger.NewLogger(configPath, cfg)
	db := database.NewDB(configPath, log)

	switch cmd := pflag.Arg(0); cmd {
	case "":
		if cfg.CustomMapPath != "" {
			domain.ChecklocalMap(cfg.CustomMapPath)
		}

		db.UpdateAnime()

		c := cron.New()
		c.AddFunc("0 0 * * *", func() { db.UpdateAnime() })
		c.Start()

		a := domain.NewAnimeUpdate(db, cfg)

		go server.StartHttp(cfg, a, log)

		sigchnl := make(chan os.Signal, 1)
		signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-sigchnl

		log.Info().Msgf("caught signal %v, shutting down", sig)
		db.Close()
		os.Exit(1)

	case "malauth":
		db.CreateDB()
		malauth.NewMalAuth(db)
		fmt.Fprintf(flag.CommandLine.Output(), "MAL API credentials saved. Testing client..\n")

		cc := malauth.NewOauth2Client(context.Background(), db)
		client := mal.NewClient(cc)
		_, _, err := client.User.MyInfo(context.Background())
		if err != nil {
			fmt.Fprintln(flag.CommandLine.Output(), "Unabled to load user info from MAL. Retry MAL authentication.")
			db.Close()
			os.Exit(1)
		}

		db.Close()
		fmt.Fprintln(flag.CommandLine.Output(), "Test successful! Run shinkuro now.")

	case "version":
		fmt.Fprintln(flag.CommandLine.Output(), "0.0")

	default:
		pflag.Usage()
		if cmd != "help" {
			os.Exit(0)
		}
	}
}
