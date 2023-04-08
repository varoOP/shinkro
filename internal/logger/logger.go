package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/domain"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(path string, c *domain.Config) *zerolog.Logger {
	logPath := filepath.Join(path, "shinkuro.log")
	logLevel := zerolog.InfoLevel
	switch c.LogLevel {
	case "TRACE":
		logLevel = zerolog.TraceLevel
	case "DEBUG":
		logLevel = zerolog.DebugLevel
	case "Error":
		logLevel = zerolog.ErrorLevel
	case "INFO":
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	lumberlog := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    c.LogMaxSize,
		MaxBackups: c.LogMaxBackups,
	}

	mw := io.MultiWriter(os.Stdout, lumberlog)
	log := zerolog.New(mw).With().Timestamp().Logger()
	return &log
}
