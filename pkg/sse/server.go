package sse

import (
	"io"
	"net/http"
	"sync"

	"github.com/r3labs/sse/v2"
)

// StreamOpts configures stream behavior
type StreamOpts struct {
	MaxEntries int  // Maximum number of entries to buffer
	AutoReplay bool // Automatically replay buffered entries to new clients
}

// Server wraps the standard SSE server with buffering and replay capabilities
type Server struct {
	*sse.Server
	mu       sync.RWMutex
	buffers  map[string]*buffer // stream name -> buffer
	opts     map[string]StreamOpts
}

// buffer maintains a circular buffer of events
type buffer struct {
	mu      sync.RWMutex
	entries []*sse.Event
	maxSize int
	head    int // current write position
	size    int // current number of entries
}

// New creates a new buffered SSE server
func New() *Server {
	return &Server{
		Server:  sse.New(),
		buffers: make(map[string]*buffer),
		opts:    make(map[string]StreamOpts),
	}
}

// CreateStreamWithOpts creates a stream with buffering options
func (s *Server) CreateStreamWithOpts(streamName string, opts StreamOpts) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create the underlying stream
	s.Server.CreateStream(streamName)

	// Create buffer if MaxEntries > 0
	if opts.MaxEntries > 0 {
		s.buffers[streamName] = &buffer{
			entries: make([]*sse.Event, opts.MaxEntries),
			maxSize: opts.MaxEntries,
		}
	}

	s.opts[streamName] = opts
}

// Publish publishes an event to a stream and buffers it if configured
func (s *Server) Publish(streamName string, event *sse.Event) {
	s.mu.RLock()
	buf, hasBuffer := s.buffers[streamName]
	s.mu.RUnlock()

	// Add to buffer if configured
	if hasBuffer && buf != nil {
		buf.add(event)
	}

	// Publish to SSE server
	s.Server.Publish(streamName, event)
}

// ServeHTTP handles HTTP requests with replay support
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		http.Error(w, "stream parameter required", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	opts, hasOpts := s.opts[streamName]
	buf, hasBuffer := s.buffers[streamName]
	s.mu.RUnlock()

	// Replay buffered entries if AutoReplay is enabled
	if hasOpts && opts.AutoReplay && hasBuffer && buf != nil {
		replayWriter := &replayResponseWriter{
			ResponseWriter: w,
			flusher:       w.(http.Flusher),
		}

		// Replay all buffered entries
		events := buf.getAll()
		for _, event := range events {
			if event != nil {
				writeSSEEvent(replayWriter, event)
				replayWriter.Flush()
			}
		}

		// Continue with normal SSE handling
		s.Server.ServeHTTP(w, r)
		return
	}

	// No replay, use standard SSE handling
	s.Server.ServeHTTP(w, r)
}

// replayResponseWriter wraps http.ResponseWriter to write SSE events
type replayResponseWriter struct {
	http.ResponseWriter
	flusher http.Flusher
}

func (w *replayResponseWriter) Flush() {
	if w.flusher != nil {
		w.flusher.Flush()
	}
}

// writeSSEEvent writes an SSE event in the proper format
func writeSSEEvent(w io.Writer, event *sse.Event) {
	if event.ID != nil {
		_, _ = io.WriteString(w, "id: "+string(event.ID)+"\n")
	}
	if event.Event != nil {
		_, _ = io.WriteString(w, "event: "+string(event.Event)+"\n")
	}
	if event.Data != nil {
		_, _ = io.WriteString(w, "data: ")
		_, _ = w.Write(event.Data)
		_, _ = io.WriteString(w, "\n")
	}
	_, _ = io.WriteString(w, "\n")
}

// add adds an event to the circular buffer
func (b *buffer) add(event *sse.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Copy the event to avoid reference issues
	eventCopy := &sse.Event{
		ID:    event.ID,
		Event: event.Event,
		Data:  make([]byte, len(event.Data)),
	}
	copy(eventCopy.Data, event.Data)

	b.entries[b.head] = eventCopy
	b.head = (b.head + 1) % b.maxSize
	if b.size < b.maxSize {
		b.size++
	}
}

// getAll returns all buffered entries in order
func (b *buffer) getAll() []*sse.Event {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.size == 0 {
		return nil
	}

	result := make([]*sse.Event, b.size)
	if b.size == b.maxSize {
		// Buffer is full, start from head (oldest)
		for i := 0; i < b.size; i++ {
			idx := (b.head + i) % b.maxSize
			if b.entries[idx] != nil {
				result[i] = b.entries[idx]
			}
		}
	} else {
		// Buffer not full, start from beginning
		for i := 0; i < b.size; i++ {
			if b.entries[i] != nil {
				result[i] = b.entries[i]
			}
		}
	}

	return result
}

