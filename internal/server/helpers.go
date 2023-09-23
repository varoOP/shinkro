package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/varoOP/shinkro/internal/domain"
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
