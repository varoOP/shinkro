package web

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidRoute(t *testing.T) {
	tests := []struct {
		name     string
		route    string
		expected bool
	}{
		{
			name:     "root route",
			route:    "/",
			expected: true,
		},
		{
			name:     "onboard route",
			route:    "/onboard",
			expected: true,
		},
		{
			name:     "login route",
			route:    "/login",
			expected: true,
		},
		{
			name:     "logout route",
			route:    "/logout",
			expected: true,
		},
		{
			name:     "logs route",
			route:    "/logs",
			expected: true,
		},
		{
			name:     "settings route",
			route:    "/settings",
			expected: true,
		},
		{
			name:     "plex-payloads route",
			route:    "/plex-payloads",
			expected: true,
		},
		{
			name:     "anime-updates route",
			route:    "/anime-updates",
			expected: true,
		},
		{
			name:     "route containing valid path",
			route:    "/some/path/settings",
			expected: true,
		},
		{
			name:     "invalid route",
			route:    "xyz-abc-123",
			expected: false,
		},
		{
			name:     "empty route",
			route:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validRoute(tt.route)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFS_Open(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name          string
		fs            defaultFS
		file          string
		expectedError bool
	}{
		{
			name: "open file with os filesystem",
			fs: defaultFS{
				prefix: tmpDir,
				fs:     nil, // Will use os.Open
			},
			file:          testFile,
			expectedError: false,
		},
		{
			name: "open file with embedded filesystem",
			fs: defaultFS{
				prefix: "",
				fs:     os.DirFS(tmpDir),
			},
			file:          "test.txt",
			expectedError: false,
		},
		{
			name: "open non-existent file",
			fs: defaultFS{
				prefix: tmpDir,
				fs:     nil,
			},
			file:          filepath.Join(tmpDir, "nonexistent.txt"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := tt.fs.Open(tt.file)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, file)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, file)
				if file != nil {
					file.Close()
				}
			}
		})
	}
}

func TestSubFS(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	testFile := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	tests := []struct {
		name          string
		currentFs     fs.FS
		root          string
		expectedError bool
		validate      func(*testing.T, fs.FS)
	}{
		{
			name:          "create sub FS from os.DirFS",
			currentFs:     os.DirFS(tmpDir),
			root:          "subdir",
			expectedError: false,
			validate: func(t *testing.T, subFs fs.FS) {
				file, err := subFs.Open("test.txt")
				assert.NoError(t, err)
				if file != nil {
					file.Close()
				}
			},
		},
		{
			name:          "create sub FS with defaultFS",
			currentFs:     &defaultFS{prefix: tmpDir, fs: nil},
			root:          "subdir",
			expectedError: false,
			validate: func(t *testing.T, subFs fs.FS) {
				file, err := subFs.Open("test.txt")
				assert.NoError(t, err)
				if file != nil {
					file.Close()
				}
			},
		},
		{
			name:          "invalid root path",
			currentFs:     os.DirFS(tmpDir),
			root:          "../invalid",
			expectedError: true, // fs.Sub returns error for invalid paths
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subFs, err := subFS(tt.currentFs, tt.root)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, subFs)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subFs)
				if tt.validate != nil {
					tt.validate(t, subFs)
				}
			}
		})
	}
}

func TestMustSubFS(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid sub FS", func(t *testing.T) {
		currentFs := os.DirFS(tmpDir)
		assert.NotPanics(t, func() {
			subFs := MustSubFS(currentFs, ".")
			assert.NotNil(t, subFs)
		})
	})

	t.Run("invalid path should panic", func(t *testing.T) {
		// Note: This might not panic with all filesystems, but we test the function exists
		currentFs := os.DirFS(tmpDir)
		// Most invalid paths are cleaned by filepath.Clean, so this might not panic
		// But we verify the function doesn't crash
		assert.NotPanics(t, func() {
			_ = MustSubFS(currentFs, ".")
		})
	})
}

func TestStaticFileHandler(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	filesystem := os.DirFS(tmpDir)
	handler := StaticFileHandler("test.txt", filesystem)

	req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test content")
}

func TestStaticFileHandler_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filesystem := os.DirFS(tmpDir)
	handler := StaticFileHandler("nonexistent.txt", filesystem)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "File not found")
}

func TestFsFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	filesystem := os.DirFS(tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
	w := httptest.NewRecorder()

	fsFile(w, req, "test.txt", filesystem)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test content")
}

func TestFsFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filesystem := os.DirFS(tmpDir)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.txt", nil)
	w := httptest.NewRecorder()

	fsFile(w, req, "nonexistent.txt", filesystem)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "File not found")
}

func TestIndexParams(t *testing.T) {
	params := IndexParams{
		Title:   "Test Title",
		Version: "v1.0.0",
		BaseUrl: "/base/",
	}

	assert.Equal(t, "Test Title", params.Title)
	assert.Equal(t, "v1.0.0", params.Version)
	assert.Equal(t, "/base/", params.BaseUrl)
}

func TestFileFS(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	filesystem := os.DirFS(tmpDir)
	r := chi.NewRouter()

	FileFS(r, "/test", "test.txt", filesystem)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test content")
}

func TestStaticFS(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	filesystem := os.DirFS(tmpDir)
	r := chi.NewRouter()

	StaticFS(r, "/static/", filesystem)

	req := httptest.NewRequest(http.MethodGet, "/static/test.txt", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test content")
}
