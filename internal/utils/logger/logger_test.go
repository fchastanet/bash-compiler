package logger

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type handler struct {
	Record *slog.Record
	Attrs  []slog.Attr
}

func (*handler) Enabled(context.Context, slog.Level) bool { return true }

func (h *handler) Handle(
	_ context.Context,
	//nolint:gocritic // hugeparam: test
	record slog.Record,
) error {
	h.Record = &record
	return nil
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.Attrs = attrs
	return &handler{} //nolint:exhaustruct // test
}

func (*handler) WithGroup(_ string) slog.Handler {
	return &handler{} //nolint:exhaustruct // test
}

func TestInitLogger(t *testing.T) {
	slogHandler := mockSlog()
	exitTriggered := false

	// Prepare testing
	myPrivateExitFunction = func(_ int) {
		exitTriggered = true
	}

	t.Run("Check nil", func(t *testing.T) {
		Check(nil)
		assert.Equal(t, exitTriggered, false)
	})

	t.Run("Check error", func(t *testing.T) {
		Check(&fs.PathError{Op: "", Path: "", Err: nil})
		assert.Equal(t, exitTriggered, true)
		assert.IsType(t, &slog.Record{}, slogHandler.Record) //nolint:exhaustruct // test
		assert.Equal(t, slogHandler.Record.Level, slog.LevelError)
		assert.Equal(t, slogHandler.Record.Message, "Error")
		var attrsSeen []string
		slogHandler.Record.Attrs(
			func(a slog.Attr) bool {
				attrsSeen = append(attrsSeen, a.Key)
				switch a.Key {
				case LogFieldFilePath:
					assert.Contains(t, a.Value.String(), "internal/utils/logger/logger_test.go")
				case LogFieldLineNumber:
					var i int64 = 1
					assert.IsTypef(t, i, a.Value.Int64(), "test")
				case LogFieldErr:
					assert.IsType(t, &fs.PathError{Op: "", Path: "", Err: nil}, a.Value.Any())
				}
				return true
			},
		)
		assert.Equal(t, attrsSeen, []string{LogFieldFilePath, LogFieldLineNumber, LogFieldErr})
	})

	// Restore if need
	myPrivateExitFunction = os.Exit
}

func mockSlog() *handler {
	slogHandler := handler{} //nolint:exhaustruct // test
	logger := slog.New(&slogHandler)
	slog.SetDefault(logger)
	return &slogHandler
}

func TestDebugSaveIntermediateFile(t *testing.T) {
	slogHandler := mockSlog()
	targetFile := getTargetFile(os.TempDir(), "testBaseName", ".log")
	defer os.Remove(targetFile)

	assertResult := func() (bytes []byte) {
		assert.FileExists(t, targetFile)
		bytes, _ = os.ReadFile(targetFile)
		assert.IsType(t, &slog.Record{}, slogHandler.Record) //nolint:exhaustruct // test
		assert.Equal(t, slogHandler.Record.Level, slog.LevelDebug)
		assert.Equal(t, slogHandler.Record.Message, "KeepIntermediateFiles - merged config file")
		var attrsSeen []string
		slogHandler.Record.Attrs(
			func(a slog.Attr) bool {
				attrsSeen = append(attrsSeen, a.Key)
				if a.Key == LogFieldFilePath {
					assert.Equal(t, targetFile, a.Value.String())
				}
				return true
			},
		)
		assert.Equal(t, attrsSeen, []string{LogFieldFilePath})
		return bytes
	}

	t.Run("TestDebugSaveIntermediateFile", func(t *testing.T) {
		DebugSaveIntermediateFile(os.TempDir(), "testBaseName", ".log", "content")
		bytes := assertResult()
		assert.Equal(t, "content", string(bytes))
	})

	t.Run("TestDebugCopyIntermediateFile", func(t *testing.T) {
		DebugCopyIntermediateFile(
			os.TempDir(), "testBaseName", ".log",
			"testsData/DebugCopyIntermediateFile.txt",
		)
		bytes := assertResult()
		assert.Equal(t, "content\n", string(bytes))
	})
}

func TestFancyHandleError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		res := FancyHandleError(nil)
		assert.Equal(t, false, res)
	})

	t.Run("error raised", func(t *testing.T) {
		res := FancyHandleError(&fs.PathError{Op: "", Path: "", Err: nil})
		assert.Equal(t, true, res)
	})
}
