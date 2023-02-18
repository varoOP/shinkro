package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/varoOP/shinkuro/internal/animedb"
	"github.com/varoOP/shinkuro/internal/server"
)

func main() {

	db := animedb.NewDB("file:./anime.db?cache=shared&mode=rwc&_journal_mode=WAL")
	animedb.CreateDB(db)

	c := cron.New()
	c.AddFunc("0 0 * * *", func() { animedb.UpdateDB(db) })
	c.Start()

	go server.StartHttp(db)

	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl, syscall.SIGINT)
	<-sigchnl
	db.Close()
	fmt.Println("\nExited after closing db.")
	os.Exit(0)

}
