package http

import (
	"net/http"
)

// LivenessHandler handles the /healthz/liveness endpoint
// This is a simple health check endpoint that returns 200 OK if the server is running
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
