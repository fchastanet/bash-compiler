// Package log allowing to load logger configuration
package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fchastanet/bash-compiler/internal/files"
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

func DebugSaveGeneratedFile(
	targetDir string, basename string, suffix string, tempYamlFile string,
) (err error) {
	targetFile := filepath.Join(
		targetDir,
		fmt.Sprintf("%s%s", basename, suffix),
	)
	err = files.Copy(tempYamlFile, targetFile)
	if err != nil {
		return err
	}
	slog.Info("KeepIntermediateFiles", "merged config file", targetFile)
	return nil
}

func DebugCopyGeneratedFile(
	targetDir string, basename string, suffix string, code string,
) (err error) {
	targetFile := filepath.Join(
		targetDir,
		fmt.Sprintf("%s%s", basename, suffix),
	)
	err = os.WriteFile(targetFile, []byte(code), files.UserReadWriteExecutePerm)
	slog.Info("KeepIntermediateFiles", "merged config file", targetFile)
	return err
}
