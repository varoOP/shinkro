package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/nstratos/go-myanimelist/mal"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"

	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/server"
)

func main() {

	var config_path string

	k := koanf.New(".")

	pflag.StringVar(&config_path, "config", ".", "Absolute path to shinkuro's configuration directory")

	pflag.Parse()

	cfg := config.NewConfigPath(config_path)

	if err := k.Load(file.Provider(cfg.Config), toml.Parser()); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	logger, err := os.OpenFile(cfg.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatalln(err)
	}

	mw := io.MultiWriter(os.Stdout, logger)

	log.SetOutput(mw)

	db := animedb.NewDB(cfg.Dsn)
	animedb.UpdateDB(db)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() { animedb.UpdateDB(db) })
	c.Start()

	oauth_client := server.NewOauth2Client(context.Background(), k.String("mal_client_id"), k.String("mal_client_secret"), cfg.Token)
	client := mal.NewClient(oauth_client)

	go server.StartHttp(db, client, k.String("host"), k.Int("port"))

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	sig := <-sigchnl
	log.Println("Caught", sig, "shutting down")
	db.Close()
	logger.Close()
	os.Exit(1)
}
