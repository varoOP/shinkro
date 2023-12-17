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

func NewLogger(path string, c *domain.Config) zerolog.Logger {
	logPath := filepath.Join(path, "shinkro.log")
	lumberlog := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    c.LogMaxSize,
		MaxBackups: c.LogMaxBackups,
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	mw := io.MultiWriter(
		zerolog.ConsoleWriter{
			TimeFormat: time.DateTime,
			Out:        os.Stdout,
		},
		zerolog.ConsoleWriter{
			TimeFormat: time.DateTime,
			Out:        lumberlog,
		},
	)

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
