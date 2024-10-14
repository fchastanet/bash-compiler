// Package log allowing to load logger configuration
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fchastanet/bash-compiler/internal/utils/files"
)

const (
	LogFieldFunc                   string = "func"
	LogFieldFilePath               string = "file"
	LogFieldFilePathExpanded       string = "filePathExpanded"
	LogFieldDirPath                string = "dirPath"
	LogFieldDirPathExpanded        string = "dirPathExpanded"
	LogFieldErr                    string = "err"
	LogFieldLineNumber             string = "lineNumber"
	LogFieldLineContent            string = "line"
	LogFieldTemplateDirs           string = "templateDirs"
	LogFieldTemplateName           string = "templateName"
	LogFieldTemplateData           string = "templateData"
	LogFieldAvailableTemplateFiles string = "availableTemplateFiles"
	LogFieldVariableName           string = "variableName"
	LogFieldVariableValue          string = "variableValue"
)

// No way to mock os.Exit function easily
//
//nolint:gochecknoglobals // no other solution for now
var myPrivateExitFunction = os.Exit

// InitLogger initializes the logger in slog instance
func InitLogger(level int) {
	slogLevel := slog.Level(level)
	opts := &slog.HandlerOptions{
		AddSource:   slogLevel == slog.LevelDebug,
		Level:       slogLevel,
		ReplaceAttr: nil,
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
		slog.Error(
			"Error",
			LogFieldFilePath, filename,
			LogFieldLineNumber, line,
			LogFieldErr, e,
		)
		myPrivateExitFunction(1)
	}
}

func getTargetFile(targetDir string, basename string, suffix string) string {
	return filepath.Join(
		targetDir,
		fmt.Sprintf("%s%s", basename, suffix),
	)
}

func DebugCopyIntermediateFile(
	targetDir string, basename string, suffix string, tempYamlFile string,
) (err error) {
	targetFile := getTargetFile(targetDir, basename, suffix)
	err = files.Copy(tempYamlFile, targetFile)
	if err != nil {
		return err
	}
	slog.Debug(
		"KeepIntermediateFiles - merged config file",
		LogFieldFilePath, targetFile,
	)

	return nil
}

func DebugSaveIntermediateFile(
	targetDir string, basename string, suffix string, content string,
) (err error) {
	targetFile := getTargetFile(targetDir, basename, suffix)
	err = os.WriteFile(targetFile, []byte(content), files.UserReadWriteExecutePerm)
	slog.Debug(
		"KeepIntermediateFiles - merged config file",
		LogFieldFilePath, targetFile,
	)

	return err
}

// this logs the function name as well.
func FancyHandleError(err error) bool {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)

		slog.Error(
			"error",
			LogFieldFunc, runtime.FuncForPC(pc).Name(),
			LogFieldFilePath, filename,
			LogFieldLineNumber, line,
			LogFieldErr, err,
		)

		return true
	}

	return false
}
