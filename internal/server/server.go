package server

import (
	"log"
	"net/http"
)

func StartHttp(addr, baseUrl string, a *AnimeUpdate) {
	http.Handle(baseUrl, a)
	log.Println("Started listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
