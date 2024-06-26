// Package log allowing to load logger configuration
package logger

import (
	"log"
	"log/slog"
	"os"
	"runtime"
)

// InitLogger initializes the logger in slog instance
func InitLogger(level int) {
	opts := &slog.HandlerOptions{
		Level: slog.Level(level),
	}
	handler := slog.NewTextHandler(os.Stderr, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Check(e error) {
	if e != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", filename, line, e)
	}
}
