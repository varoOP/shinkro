package server

import (
	"log"
	"net/http"

	"github.com/varoOP/shinkuro/internal/config"
)

func StartHttp(cfg *config.Config, a *AnimeUpdate) {

	log.Println("Started listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, a))
}
