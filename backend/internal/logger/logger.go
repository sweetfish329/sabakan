package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is the global application logger.
var Logger *slog.Logger

// Init initializes the global logger with the specified level and format.
// Level can be: debug, info, warn, error
// Format can be: json, text
func Init(level, format string) {
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// parseLevel converts a string level to slog.Level.
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Debug logs a debug message.
func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

// Info logs an info message.
func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

// Warn logs a warning message.
func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

// Error logs an error message.
func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}
