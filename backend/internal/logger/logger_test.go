package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestInit_TextFormat(t *testing.T) {
	// Act
	Init("info", "text")

	// Assert
	if Logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestInit_JSONFormat(t *testing.T) {
	// Act
	Init("debug", "json")

	// Assert
	if Logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"INFO", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"unknown", slog.LevelInfo}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	// Arrange
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	Logger = slog.New(handler)

	// Act
	Debug("test debug message")
	Info("test info message")
	Warn("test warn message")
	Error("test error message")

	// Assert
	output := buf.String()
	if !strings.Contains(output, "test debug message") {
		t.Error("Debug message not found in output")
	}
	if !strings.Contains(output, "test info message") {
		t.Error("Info message not found in output")
	}
	if !strings.Contains(output, "test warn message") {
		t.Error("Warn message not found in output")
	}
	if !strings.Contains(output, "test error message") {
		t.Error("Error message not found in output")
	}
}
