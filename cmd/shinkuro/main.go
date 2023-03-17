package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"

	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/logger"
	"github.com/varoOP/shinkuro/internal/malauth"
	"github.com/varoOP/shinkuro/internal/mapping"
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

	pflag.StringVar(&configPath, "config", "", "path to configuration")
	pflag.Parse()
	logger := logger.NewLogger(filepath.Join(configPath, "shinkuro.log"))

	switch cmd := pflag.Arg(0); cmd {
	case "":
		var (
			cfg = config.NewConfig(configPath)
			db  = database.NewDB(cfg.Dsn)
		)

		if cfg.CustomMap != "" {
			mapping.ChecklocalMap(cfg.CustomMap)
		}

		database.UpdateAnime(db)

		c := cron.New()
		c.AddFunc("0 0 * * *", func() { database.UpdateAnime(db) })
		c.Start()

		cc := malauth.NewOauth2Client(context.Background(), db)
		client := mal.NewClient(cc)

		a := server.NewAnimeUpdate(db, client, cfg)

		go server.StartHttp(cfg.Addr, cfg.BaseUrl, a)

		sigchnl := make(chan os.Signal, 1)
		signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-sigchnl

		log.Println("Caught", sig, "shutting down")
		db.Close()
		logger.Close()
		os.Exit(1)

	case "malauth":
		var (
			cfg = config.NewConfig(configPath)
			db  = database.NewDB(cfg.Dsn)
		)
		database.CreateDB(db)
		malauth.NewMalAuth(db)
		fmt.Fprintf(flag.CommandLine.Output(), "MAL API credentials saved.\nTesting client..\n")
		cc := malauth.NewOauth2Client(context.Background(), db)
		client := mal.NewClient(cc)
		_, _, err := client.User.MyInfo(context.Background())
		if err != nil {
			fmt.Fprintln(flag.CommandLine.Output(), "Unabled to load user info from MAL. Retry MAL authentication.")
			db.Close()
			os.Exit(1)
		}
		db.Close()
		fmt.Fprintln(flag.CommandLine.Output(), "Test successful.")

	case "version":
		fmt.Fprintln(flag.CommandLine.Output(), "0.0")

	default:
		pflag.Usage()
		if cmd != "help" {
			os.Exit(0)
		}
	}
}
