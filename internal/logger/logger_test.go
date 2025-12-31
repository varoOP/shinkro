package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected zerolog.Level
	}{
		{
			name:     "TRACE level",
			level:    "TRACE",
			expected: zerolog.TraceLevel,
		},
		{
			name:     "DEBUG level",
			level:    "DEBUG",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "INFO level",
			level:    "INFO",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "ERROR level",
			level:    "ERROR",
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "lowercase info",
			level:    "info",
			expected: zerolog.InfoLevel, // Default fallback
		},
		{
			name:     "empty string",
			level:    "",
			expected: zerolog.InfoLevel, // Default fallback
		},
		{
			name:     "invalid level",
			level:    "INVALID",
			expected: zerolog.InfoLevel, // Default fallback
		},
		{
			name:     "WARN level (not explicitly handled)",
			level:    "WARN",
			expected: zerolog.InfoLevel, // Default fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

