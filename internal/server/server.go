package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/domain"
)

func StartHttp(cfg *domain.Config, a *domain.AnimeUpdate, log *zerolog.Logger) {
	http.Handle(cfg.BaseUrl, a)
	addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	serverLog := log.With().Str("module", "server").Logger()
	log = &serverLog
	log.Info().Msgf("started listening on %v", addr)
	log.Fatal().Err(http.ListenAndServe(addr, nil))
}
