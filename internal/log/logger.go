// Package log allowing to load logger configuration
package log

import (
	"log/slog"
	"os"
)

// InitLogger initializes the logger in slog instance
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
