package sse

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	server := New()
	assert.NotNil(t, server)
	assert.NotNil(t, server.Server)
	assert.NotNil(t, server.buffers)
	assert.NotNil(t, server.opts)
	assert.Equal(t, 0, len(server.buffers))
	assert.Equal(t, 0, len(server.opts))
}

func TestCreateStreamWithOpts(t *testing.T) {
	server := New()

	tests := []struct {
		name        string
		streamName  string
		opts        StreamOpts
		expectBuffer bool
	}{
		{
			name:        "create stream with buffer",
			streamName:  "test-stream",
			opts:        StreamOpts{MaxEntries: 10, AutoReplay: true},
			expectBuffer: true,
		},
		{
			name:        "create stream without buffer",
			streamName:  "test-stream-2",
			opts:        StreamOpts{MaxEntries: 0, AutoReplay: false},
			expectBuffer: false,
		},
		{
			name:        "create stream with buffer but no replay",
			streamName:  "test-stream-3",
			opts:        StreamOpts{MaxEntries: 5, AutoReplay: false},
			expectBuffer: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.CreateStreamWithOpts(tt.streamName, tt.opts)

			// Verify buffer was created if MaxEntries > 0
			server.mu.RLock()
			_, hasBuffer := server.buffers[tt.streamName]
			opts, hasOpts := server.opts[tt.streamName]
			server.mu.RUnlock()

			assert.Equal(t, tt.expectBuffer, hasBuffer)
			assert.True(t, hasOpts)
			assert.Equal(t, tt.opts.MaxEntries, opts.MaxEntries)
			assert.Equal(t, tt.opts.AutoReplay, opts.AutoReplay)
		})
	}
}

func TestPublish(t *testing.T) {
	server := New()
	server.CreateStreamWithOpts("test-stream", StreamOpts{MaxEntries: 3, AutoReplay: true})

	event1 := &sse.Event{
		ID:    []byte("1"),
		Event: []byte("test"),
		Data:  []byte("data1"),
	}

	event2 := &sse.Event{
		ID:    []byte("2"),
		Event: []byte("test"),
		Data:  []byte("data2"),
	}

	// Publish events
	server.Publish("test-stream", event1)
	server.Publish("test-stream", event2)

	// Verify events were buffered
	server.mu.RLock()
	buf := server.buffers["test-stream"]
	server.mu.RUnlock()

	assert.NotNil(t, buf)
	events := buf.getAll()
	assert.Equal(t, 2, len(events))
	assert.Equal(t, "data1", string(events[0].Data))
	assert.Equal(t, "data2", string(events[1].Data))
}

func TestPublish_BufferOverflow(t *testing.T) {
	server := New()
	server.CreateStreamWithOpts("test-stream", StreamOpts{MaxEntries: 2, AutoReplay: true})

	// Publish 3 events to a buffer of size 2
	event1 := &sse.Event{ID: []byte("1"), Data: []byte("data1")}
	event2 := &sse.Event{ID: []byte("2"), Data: []byte("data2")}
	event3 := &sse.Event{ID: []byte("3"), Data: []byte("data3")}

	server.Publish("test-stream", event1)
	server.Publish("test-stream", event2)
	server.Publish("test-stream", event3)

	// Verify only last 2 events are in buffer (circular buffer)
	server.mu.RLock()
	buf := server.buffers["test-stream"]
	server.mu.RUnlock()

	events := buf.getAll()
	assert.Equal(t, 2, len(events))
	// Should contain event2 and event3 (event1 was overwritten)
	assert.Contains(t, []string{string(events[0].Data), string(events[1].Data)}, "data2")
	assert.Contains(t, []string{string(events[0].Data), string(events[1].Data)}, "data3")
}

func TestPublish_NoBuffer(t *testing.T) {
	server := New()
	server.CreateStreamWithOpts("test-stream", StreamOpts{MaxEntries: 0, AutoReplay: false})

	event := &sse.Event{
		ID:    []byte("1"),
		Event: []byte("test"),
		Data:  []byte("data"),
	}

	// Publish should not panic even without buffer
	assert.NotPanics(t, func() {
		server.Publish("test-stream", event)
	})

	// Verify no buffer was created
	server.mu.RLock()
	_, hasBuffer := server.buffers["test-stream"]
	server.mu.RUnlock()

	assert.False(t, hasBuffer)
}

func TestServeHTTP_MissingStream(t *testing.T) {
	server := New()

	req := httptest.NewRequest(http.MethodGet, "/?stream=", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "stream parameter required")
}

func TestServeHTTP_StreamNotFound(t *testing.T) {
	server := New()

	req := httptest.NewRequest(http.MethodGet, "/?stream=nonexistent", nil)
	w := httptest.NewRecorder()

	// This will call the underlying SSE server which may handle it differently
	// We just verify it doesn't panic
	assert.NotPanics(t, func() {
		server.ServeHTTP(w, req)
	})
}

func TestBuffer_Add(t *testing.T) {
	buf := &buffer{
		entries: make([]*sse.Event, 3),
		maxSize: 3,
	}

	event1 := &sse.Event{ID: []byte("1"), Data: []byte("data1")}
	event2 := &sse.Event{ID: []byte("2"), Data: []byte("data2")}

	buf.add(event1)
	buf.add(event2)

	assert.Equal(t, 2, buf.size)
	assert.Equal(t, 2, buf.head)
}

func TestBuffer_GetAll(t *testing.T) {
	buf := &buffer{
		entries: make([]*sse.Event, 3),
		maxSize: 3,
	}

	event1 := &sse.Event{ID: []byte("1"), Data: []byte("data1")}
	event2 := &sse.Event{ID: []byte("2"), Data: []byte("data2")}
	event3 := &sse.Event{ID: []byte("3"), Data: []byte("data3")}

	buf.add(event1)
	buf.add(event2)
	buf.add(event3)

	events := buf.getAll()
	assert.Equal(t, 3, len(events))
	assert.Equal(t, "data1", string(events[0].Data))
	assert.Equal(t, "data2", string(events[1].Data))
	assert.Equal(t, "data3", string(events[2].Data))
}

func TestBuffer_GetAll_Empty(t *testing.T) {
	buf := &buffer{
		entries: make([]*sse.Event, 3),
		maxSize: 3,
	}

	events := buf.getAll()
	assert.Nil(t, events)
}

func TestBuffer_GetAll_Circular(t *testing.T) {
	buf := &buffer{
		entries: make([]*sse.Event, 3),
		maxSize: 3,
	}

	// Fill buffer completely
	event1 := &sse.Event{ID: []byte("1"), Data: []byte("data1")}
	event2 := &sse.Event{ID: []byte("2"), Data: []byte("data2")}
	event3 := &sse.Event{ID: []byte("3"), Data: []byte("data3")}
	event4 := &sse.Event{ID: []byte("4"), Data: []byte("data4")} // This will overwrite event1

	buf.add(event1)
	buf.add(event2)
	buf.add(event3)
	buf.add(event4) // Overwrites event1

	events := buf.getAll()
	assert.Equal(t, 3, len(events))
	// Should contain event2, event3, event4 (event1 was overwritten)
	dataValues := []string{
		string(events[0].Data),
		string(events[1].Data),
		string(events[2].Data),
	}
	assert.Contains(t, dataValues, "data2")
	assert.Contains(t, dataValues, "data3")
	assert.Contains(t, dataValues, "data4")
	assert.NotContains(t, dataValues, "data1")
}

func TestWriteSSEEvent(t *testing.T) {
	tests := []struct {
		name  string
		event *sse.Event
		validate func(*testing.T, string)
	}{
		{
			name: "complete event",
			event: &sse.Event{
				ID:    []byte("123"),
				Event: []byte("test"),
				Data:  []byte("test data"),
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: 123")
				assert.Contains(t, output, "event: test")
				assert.Contains(t, output, "data: test data")
			},
		},
		{
			name: "event without ID",
			event: &sse.Event{
				Event: []byte("test"),
				Data:  []byte("test data"),
			},
			validate: func(t *testing.T, output string) {
				assert.NotContains(t, output, "id:")
				assert.Contains(t, output, "event: test")
				assert.Contains(t, output, "data: test data")
			},
		},
		{
			name: "event without event type",
			event: &sse.Event{
				ID:   []byte("123"),
				Data: []byte("test data"),
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: 123")
				assert.NotContains(t, output, "event:")
				assert.Contains(t, output, "data: test data")
			},
		},
		{
			name: "event without data",
			event: &sse.Event{
				ID:    []byte("123"),
				Event: []byte("test"),
			},
			validate: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: 123")
				assert.Contains(t, output, "event: test")
				assert.NotContains(t, output, "data:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writeSSEEvent(&buf, tt.event)
			output := buf.String()
			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

