package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLivenessHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz/liveness", nil)
	w := httptest.NewRecorder()

	LivenessHandler(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	
	body := w.Body.String()
	assert.Equal(t, `{"status":"ok"}`, body)
}

func TestLivenessHandler_MultipleCalls(t *testing.T) {
	// Test that the handler works consistently across multiple calls
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthz/liveness", nil)
		w := httptest.NewRecorder()

		LivenessHandler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"status":"ok"}`, w.Body.String())
	}
}

