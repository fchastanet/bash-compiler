// Package functions
package functions

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	sprig "github.com/Masterminds/sprig/v3"
)

var errorIfEmptyError = errors.New("Value cannot be empty")
var errorNotSupportedType = errors.New("Type not supported")

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

func bashVariableRef(variableName string) string {
	return fmt.Sprintf("${%s}", variableName)
}

func len(list interface{}) (int, error) {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)
		return l2.Len(), nil
	case reflect.String:
		return len(list.(string))
	default:
		return 0, errorNotSupportedType
	}
}

func FuncMap() map[string]interface{} {
	funcMap := sprig.FuncMap()
	funcMap["errorIfEmpty"] = errorIfEmpty
	funcMap["logWarn"] = logWarn
	// string functions
	funcMap["len"] = len
	funcMap["bashVariableRef"] = bashVariableRef
	// templates functions
	funcMap["include"] = include
	// YAML functions
	funcMap["fromYAMLFile"] = FromYAMLFile

	return funcMap
}
