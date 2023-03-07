package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/robfig/cron/v3"

	"github.com/varoOP/shinkuro/internal/config"
	"github.com/varoOP/shinkuro/internal/database"
	"github.com/varoOP/shinkuro/internal/server"
)

func main() {

	cfg := config.NewConfig()

	db := database.NewDB(cfg.Dsn)
	database.UpdateDB(db)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() { database.UpdateDB(db) })
	c.Start()

	oauth_client := server.NewOauth2Client(context.Background(), cfg.K.String("mal_client_id"), cfg.K.String("mal_client_secret"), cfg.Token)
	client := mal.NewClient(oauth_client)

	ac := server.NewAnimeCon(client, db)

	go server.StartHttp(cfg, ac)

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	sig := <-sigchnl

	log.Println("Caught", sig, "shutting down")
	db.Close()
	cfg.Logger.Close()
	os.Exit(1)
}
