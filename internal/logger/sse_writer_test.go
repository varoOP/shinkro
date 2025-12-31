package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeedsQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "simple string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "string with space",
			input:    "hello world",
			expected: true,
		},
		{
			name:     "string with backslash",
			input:    "hello\\world",
			expected: true,
		},
		{
			name:     "string with quote",
			input:    "hello\"world",
			expected: true,
		},
		{
			name:     "string with control character (newline)",
			input:    "hello\nworld",
			expected: true,
		},
		{
			name:     "string with tab",
			input:    "hello\tworld",
			expected: true,
		},
		{
			name:     "string with null byte",
			input:    "hello\x00world",
			expected: true,
		},
		{
			name:     "string with high ASCII",
			input:    "hello\x7fworld",
			expected: true,
		},
		{
			name:     "string with printable ASCII only",
			input:    "hello123",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "string with unicode (above 0x7e)",
			input:    "hello世界",
			expected: true,
		},
		{
			name:     "string with extended ASCII",
			input:    "hello\x80world",
			expected: true,
		},
		{
			name:     "string with carriage return",
			input:    "hello\rworld",
			expected: true,
		},
		{
			name:     "string with bell character",
			input:    "hello\x07world",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := needsQuote(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultPartsOrder(t *testing.T) {
	order := defaultPartsOrder()
	assert.NotNil(t, order)
	assert.Greater(t, len(order), 0)
	// Verify it contains expected fields
	assert.Contains(t, order, "caller")
	assert.Contains(t, order, "message")
}

func TestDefaultFormatMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "test message",
			expected: "test message",
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "",
		},
		{
			name:     "number input",
			input:    123,
			expected: "%!s(int=123)", // fmt.Sprintf("%s", int) produces this format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultFormatMessage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFormatFieldName(t *testing.T) {
	formatter := defaultFormatFieldName()

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string field name",
			input:    "field",
			expected: "field=",
		},
		{
			name:     "number field name",
			input:    123,
			expected: "%!s(int=123)=", // fmt.Sprintf("%s", int) produces this format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFormatFieldValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string value",
			input:    "value",
			expected: "value",
		},
		{
			name:     "number value",
			input:    123,
			expected: "%!s(int=123)", // fmt.Sprintf("%s", int) produces this format
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultFormatFieldValue(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFormatErrFieldName(t *testing.T) {
	formatter := defaultFormatErrFieldName()

	result := formatter("error")
	assert.Equal(t, "error=", result)
}

func TestDefaultFormatErrFieldValue(t *testing.T) {
	formatter := defaultFormatErrFieldValue()

	result := formatter("error message")
	assert.Equal(t, "error message=", result)
}

func TestDecodeIfBinaryToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "regular bytes",
			input:    []byte("test"),
			expected: []byte("test"),
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "nil bytes",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decodeIfBinaryToBytes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogMessage_Bytes(t *testing.T) {
	tests := []struct {
		name     string
		message  LogMessage
		validate func(*testing.T, []byte)
	}{
		{
			name: "complete message",
			message: LogMessage{
				Time:    "2024-01-01T00:00:00Z",
				Level:   "INFO",
				Message: "test message",
			},
			validate: func(t *testing.T, data []byte) {
				assert.NotNil(t, data)
				assert.Contains(t, string(data), "2024-01-01T00:00:00Z")
				assert.Contains(t, string(data), "INFO")
				assert.Contains(t, string(data), "test message")
			},
		},
		{
			name: "empty message",
			message: LogMessage{
				Time:    "",
				Level:   "",
				Message: "",
			},
			validate: func(t *testing.T, data []byte) {
				assert.NotNil(t, data)
				// Should still be valid JSON
				assert.Contains(t, string(data), "time")
				assert.Contains(t, string(data), "level")
				assert.Contains(t, string(data), "message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.message.Bytes()
			assert.NoError(t, err)
			assert.NotNil(t, result)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
