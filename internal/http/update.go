package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/varoOP/shinkro/internal/domain"
	"github.com/varoOP/shinkro/internal/update"
)

func GetLatestReleaseHandler(w http.ResponseWriter, r *http.Request) {
	tag, err := update.LatestTag(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch latest release"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(domain.UpdateRelease{TagName: tag})
}
