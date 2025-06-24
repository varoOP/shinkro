package http

import (
	"encoding/json"
	"net/http"

	"github.com/varoOP/shinkro/internal/domain"
)

func GetLatestReleaseHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.github.com/repos/varoOP/shinkro/releases/latest")
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch latest release"})
		return
	}
	defer resp.Body.Close()

	var release domain.UpdateRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to decode release response"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(release)
}
