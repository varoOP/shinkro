package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToNullString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected sql.Null[string]
	}{
		{
			name:  "non-empty string",
			input: "test",
			expected: sql.Null[string]{
				V:     "test",
				Valid: true,
			},
		},
		{
			name:  "empty string",
			input: "",
			expected: sql.Null[string]{
				V:     "",
				Valid: false,
			},
		},
		{
			name:  "whitespace string",
			input: "   ",
			expected: sql.Null[string]{
				V:     "   ",
				Valid: true,
			},
		},
		{
			name:  "unicode string",
			input: "こんにちは",
			expected: sql.Null[string]{
				V:     "こんにちは",
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toNullString(tt.input)
			assert.Equal(t, tt.expected.V, result.V)
			assert.Equal(t, tt.expected.Valid, result.Valid)
		})
	}
}

