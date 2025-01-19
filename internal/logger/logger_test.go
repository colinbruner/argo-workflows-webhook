package logger

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestDebug(t *testing.T) {
	tests := []struct {
		level      slog.Level
		output     string
		wantOutput bool
	}{
		{
			level:      slog.LevelDebug,
			output:     "test debug message",
			wantOutput: true,
		},
		{
			level:      slog.LevelInfo,
			output:     "test no debug message",
			wantOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			var buf bytes.Buffer
			opts := &slog.HandlerOptions{
				Level: tt.level,
			}
			handler := slog.NewTextHandler(&buf, opts)
			log = slog.New(handler)

			Debug(tt.output)

			if tt.wantOutput && !bytes.Contains(buf.Bytes(), []byte(tt.output)) {
				t.Errorf("Expected 'test info message' to be logged, but it wasn't")
			}
		})
	}
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&buf, opts)
	log = slog.New(handler)

	Info("test info message")

	if !bytes.Contains(buf.Bytes(), []byte("test info message")) {
		t.Errorf("Expected 'test info message' to be logged, but it wasn't")
	}
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&buf, opts)
	log = slog.New(handler)

	Warn("test warn message")

	if !bytes.Contains(buf.Bytes(), []byte("test warn message")) {
		t.Errorf("Expected 'test warn message' to be logged, but it wasn't")
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(&buf, opts)
	log = slog.New(handler)

	Error("test error message")

	if !bytes.Contains(buf.Bytes(), []byte("test error message")) {
		t.Errorf("Expected 'test error message' to be logged, but it wasn't")
	}
}
