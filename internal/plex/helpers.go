package plex

import (
	"strings"

	"github.com/varoOP/shinkro/internal/domain"
)

func isMetadataAgent(p *domain.Plex) (bool, domain.PlexSupportedAgents) {
	if strings.Contains(p.Metadata.GUID.GUID, "agents.hama") {
		return true, domain.HAMA
	}

	if strings.Contains(p.Metadata.GUID.GUID, "myanimelist") {
		return true, domain.MALAgent
	}

	if strings.Contains(p.Metadata.GUID.GUID, "plex://") {
		return true, domain.PlexAgent
	}

	return false, ""
}

func isPlexUser(p *domain.Plex, c *domain.Config) bool {
	return p.Account.Title == c.PlexUser
}

func isEvent(p *domain.Plex) bool {
	return p.Event == domain.PlexRateEvent || p.Event == domain.PlexScrobbleEvent
}

func isAnimeLibrary(p *domain.Plex, c *domain.Config) bool {
	l := strings.Join(c.AnimeLibraries, ",")
	return strings.Contains(l, p.Metadata.LibrarySectionTitle)
}

func mediaType(p *domain.Plex) bool {
	return p.Metadata.Type == domain.PlexEpisode || p.Metadata.Type == domain.PlexMovie
}

func isPlexClient(cfg *domain.Config) bool {
	return cfg.PlexToken != ""
}
