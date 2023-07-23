package server

import (
	"strings"

	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/pkg/plex"
)

func isMetadataAgent(p *plex.PlexWebhook) bool {
	return strings.Contains(p.Metadata.GUID.GUID, "com.plexapp.agents.hama") || strings.Contains(p.Metadata.GUID.GUID, "net.fribbtastic.coding.plex.myanimelist")
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

func mediaType(p *plex.PlexWebhook) (bool, string) {
	if p.Metadata.Type == "episode" {
		return true, p.Metadata.GrandparentTitle
	}

	if p.Metadata.Type == "movie" {
		return true, p.Metadata.Title
	}

	return false, ""
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
	for key, v := range in {
		if key == "apiKey" || key == "Shinkro-Api-Key" {
			for _, vv := range v {
				if vv == apiKey {
					return true
				}
			}
		}
	}

	return false
}
