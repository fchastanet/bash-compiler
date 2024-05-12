// Package functions
package functions

import (
	"errors"
	"log/slog"

	sprig "github.com/Masterminds/sprig/v3"
)

var errorIfEmptyError = errors.New("Value cannot be empty")

func errorIfEmpty(value string) (string, error) {
	if value == "" {
		return "", errorIfEmptyError
	}
	return value, nil
}

func logWarn(message string, args ...any) string {
	slog.Warn(message, args...)
	return ""
}

func FuncMap() map[string]interface{} {
	funcMap := sprig.FuncMap()
	funcMap["errorIfEmpty"] = errorIfEmpty
	funcMap["logWarn"] = logWarn
	// templates functions
	funcMap["include"] = include
	// YAML functions
	funcMap["fromYAMLFile"] = FromYAMLFile

	return funcMap
}
