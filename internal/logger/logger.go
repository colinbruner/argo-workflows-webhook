package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	logLevel := &slog.LevelVar{}
	logLevel.Set(slog.LevelDebug)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if os.Getenv("ENVIRONMENT") == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
		logLevel.Set(slog.LevelInfo)
	}

	log = slog.New(handler)
}

// Debug only prints messages when log level is set to debug
func Debug(msg string, v ...interface{}) {
	log.Debug(msg, v...)
}

// Info ...
func Info(msg string, v ...interface{}) {
	log.Info(msg, v...)
}

// Warn ...
func Warn(msg string, v ...interface{}) {
	log.Warn(msg, v...)
}

// Error ...
func Error(msg string, v ...interface{}) {
	log.Error(msg, v...)
}
