package plex

import (
	"strings"

	"github.com/varoOP/shinkro/internal/domain"
)

// func isMetadataAgent(p *domain.Plex) (bool, string) {
// 	if strings.Contains(p.Metadata.GUID.GUID, "agents.hama") {
// 		return true, "hama"
// 	}

// 	if strings.Contains(p.Metadata.GUID.GUID, "myanimelist") {
// 		return true, "mal"
// 	}

// 	if strings.Contains(p.Metadata.GUID.GUID, "plex://") {
// 		return true, "plex"
// 	}

// 	return false, ""
// }

func isPlexUser(p *domain.Plex, c *domain.Config) bool {
	return p.Account.Title == c.PlexUser
}

func isEvent(p *domain.Plex) bool {
	return p.Event == "media.rate" || p.Event == "media.scrobble"
}

func isAnimeLibrary(p *domain.Plex, c *domain.Config) bool {
	l := strings.Join(c.AnimeLibraries, ",")
	return strings.Contains(l, p.Metadata.LibrarySectionTitle)
}

func mediaType(p *domain.Plex) bool {
	return p.Metadata.Type == "episode" || p.Metadata.Type == "movie"
}
