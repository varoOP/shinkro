package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/varoOP/shinkro/internal/domain"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(c *domain.Config) zerolog.Logger {

	var mw io.Writer
	var defaultWriter io.Writer = os.Stderr

	if c.Version == "dev" {
		defaultWriter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}

	logPath := filepath.Join(c.ConfigPath, c.LogPath)
	lumberlog := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    c.LogMaxSize,
		MaxBackups: c.LogMaxBackups,
	}

	mw = io.MultiWriter(
		defaultWriter,
		lumberlog,
	)

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log := zerolog.New(mw).With().Timestamp().Logger()
	switch c.LogLevel {
	case "TRACE":
		log = log.Level(zerolog.TraceLevel)
	case "DEBUG":
		log = log.Level(zerolog.DebugLevel)
	case "ERROR":
		log = log.Level(zerolog.ErrorLevel)
	case "INFO":
		log = log.Level(zerolog.InfoLevel)
	}

	return log
}
