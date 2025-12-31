package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_StatusResponse(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		response       interface{}
		expectedStatus int
		validate       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "with response data",
			status:         http.StatusOK,
			response:       map[string]string{"message": "success"},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "success", result["message"])
			},
		},
		{
			name:           "with nil response",
			status:         http.StatusNoContent,
			response:       nil,
			expectedStatus: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, w.Code)
				assert.Empty(t, w.Body.Bytes())
			},
		},
		{
			name:           "with struct response",
			status:         http.StatusOK,
			response:       struct{ ID int }{ID: 123},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result struct{ ID int }
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, 123, result.ID)
			},
		},
		{
			name:           "with array response",
			status:         http.StatusOK,
			response:       []string{"a", "b", "c"},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result []string
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, []string{"a", "b", "c"}, result)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.StatusResponse(w, tt.status, tt.response)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestEncoder_StatusResponseMessage(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		message        string
		expectedStatus int
		validate       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "with message",
			status:         http.StatusOK,
			message:        "test message",
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "test message", result.Message)
			},
		},
		{
			name:           "empty message",
			status:         http.StatusNoContent,
			message:        "",
			expectedStatus: http.StatusNoContent,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNoContent, w.Code)
				assert.Empty(t, w.Body.Bytes())
			},
		},
		{
			name:           "special characters in message",
			status:         http.StatusOK,
			message:        "test & message <script>",
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "test & message <script>", result.Message)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.StatusResponseMessage(w, tt.status, tt.message)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestEncoder_StatusCreated(t *testing.T) {
	enc := encoder{}
	w := httptest.NewRecorder()
	enc.StatusCreated(w)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestEncoder_StatusCreatedData(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "with data",
			data: map[string]int{"id": 123},
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, w.Code)
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result map[string]int
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, 123, result["id"])
			},
		},
		{
			name: "with struct",
			data: struct{ Name string }{Name: "test"},
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, w.Code)
				var result struct{ Name string }
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "test", result.Name)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.StatusCreatedData(w, tt.data)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestEncoder_NoContent(t *testing.T) {
	enc := encoder{}
	w := httptest.NewRecorder()
	enc.NoContent(w)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestEncoder_StatusNotFound(t *testing.T) {
	enc := encoder{}
	w := httptest.NewRecorder()
	enc.StatusNotFound(w)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEncoder_NotFoundErr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "with error",
			err:  assert.AnError,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, w.Code)
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Message)
			},
		},
		{
			name: "with custom error message",
			err:  errors.New("custom not found error"),
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, w.Code)
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "custom not found error", result.Message)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.NotFoundErr(w, tt.err)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestEncoder_StatusInternalError(t *testing.T) {
	enc := encoder{}
	w := httptest.NewRecorder()
	enc.StatusInternalError(w)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEncoder_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "with error",
			err:  assert.AnError,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Message)
			},
		},
		{
			name: "with custom error message",
			err:  errors.New("custom internal error"),
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.Equal(t, "custom internal error", result.Message)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.Error(w, tt.err)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestEncoder_StatusError(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		err            error
		expectedStatus int
		validate       func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "bad request",
			status:         http.StatusBadRequest,
			err:            assert.AnError,
			expectedStatus: http.StatusBadRequest,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
				var result struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err)
				assert.NotEmpty(t, result.Message)
			},
		},
		{
			name:           "forbidden",
			status:         http.StatusForbidden,
			err:            assert.AnError,
			expectedStatus: http.StatusForbidden,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			name:           "unauthorized",
			status:         http.StatusUnauthorized,
			err:            assert.AnError,
			expectedStatus: http.StatusUnauthorized,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
	}

	enc := encoder{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			enc.StatusError(w, tt.status, tt.err)
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}
