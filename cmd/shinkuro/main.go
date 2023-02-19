package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/internal/server"
)

func main() {

	db := animedb.NewDB("file:./anime.db?cache=shared&mode=rwc&_journal_mode=WAL")
	animedb.UpdateDB(db)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() { animedb.UpdateDB(db) })
	c.Start()

	client := server.NewOauth2Client(context.Background())

	go server.StartHttp(db, client)

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	for range sigchnl {
		db.Close()
		log.Println("Exited after closing db.")
		os.Exit(1)
	}

}
