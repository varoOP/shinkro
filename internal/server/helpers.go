package server

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/hlog"
	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/tautulli"
	"github.com/varoOP/shinkro/pkg/plex"
)

const InternalServerError string = "internal server error"

func isMetadataAgent(p *plex.PlexWebhook) (bool, string) {
	if strings.Contains(p.Metadata.GUID.GUID, "agents.hama") {
		return true, "hama"
	}

	if strings.Contains(p.Metadata.GUID.GUID, "myanimelist") {
		return true, "mal"
	}

	if strings.Contains(p.Metadata.GUID.GUID, "plex://") {
		return true, "plex"
	}

	return false, ""
}

func isPlexUser(p *plex.PlexWebhook, c *domain.Config) bool {
	return p.Account.Title == c.PlexUser
}

func isEvent(p *plex.PlexWebhook) bool {
	return p.Event == "media.rate" || p.Event == "media.scrobble"
}

func isAnimeLibrary(p *plex.PlexWebhook, c *domain.Config) bool {
	l := strings.Join(c.AnimeLibraries, ",")
	return strings.Contains(l, p.Metadata.LibrarySectionTitle)
}

func mediaType(p *plex.PlexWebhook) bool {
	if p.Metadata.Type == "episode" {
		return true
	}

	if p.Metadata.Type == "movie" {
		return true
	}

	return false
}

func notify(a *domain.AnimeUpdate, err error) {
	if a.Notify.Url == "" {
		return
	}

	if err != nil {
		a.Notify.Error <- err
		return
	}

	a.Notify.Anime <- *a
}

func isAuthorized(apiKey string, in map[string][]string) bool {
	if keys, ok := in["apiKey"]; ok {
		for _, vv := range keys {
			if vv == apiKey {
				return true
			}
		}
	}

	if keys, ok := in["Shinkro-Api-Key"]; ok {
		for _, vv := range keys {
			if vv == apiKey {
				return true
			}
		}
	}

	return false
}

func contentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return "plexWebhook"
	}

	if strings.Contains(contentType, "application/json") {
		return "tautulli"
	}

	return contentType
}

func readRequest(r *http.Request) (string, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()
	return string(b), nil
}

func joinUrlPath(base, extra string) string {
	u, err := url.JoinPath(base, extra)
	if err != nil {
		return extra
	}

	return u
}

func parsePayloadBySourceType(w http.ResponseWriter, r *http.Request, sourceType string) (*plex.PlexWebhook, error) {
	log := hlog.FromRequest(r)
	switch sourceType {
	case "plexWebhook":
		return handlePlexWebhook(w, r)

	case "tautulli":
		return handleTautulli(w, r)

	default:
		log.Error().Str("sourceType", sourceType).Msg("sourceType not supported")
		return nil, errors.New("unsupported source type")
	}
}

func handlePlexWebhook(w http.ResponseWriter, r *http.Request) (*plex.PlexWebhook, error) {
	log := hlog.FromRequest(r)
	if err := r.ParseMultipartForm(0); err != nil {
		http.Error(w, "received bad request", http.StatusBadRequest)
		log.Trace().Err(err).Msg("received bad request")
		return nil, err
	}

	ps := r.PostFormValue("payload")
	if ps == "" {
		log.Info().Msg("Received empty payload from Plex, webhook added successfully.")
		return nil, nil
	}

	log.Trace().RawJSON("rawPlexPayload", []byte(ps)).Msg("")
	return plex.NewPlexWebhook([]byte(ps))
}

func handleTautulli(w http.ResponseWriter, r *http.Request) (*plex.PlexWebhook, error) {
	log := hlog.FromRequest(r)
	ps, err := readRequest(r)
	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		log.Trace().Err(err).Msg(InternalServerError)
		return nil, err
	}

	log.Trace().RawJSON("rawPlexPayload", []byte(ps)).Msg("")
	return tautulli.ToPlex([]byte(ps))
}
