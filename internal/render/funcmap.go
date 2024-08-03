// Package render
package render

import (
	"fmt"
	"reflect"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/fchastanet/bash-compiler/internal/utils/bash"
)

type notSupportedTypeError struct {
	error
	ObjectType string
}

func (e *notSupportedTypeError) Error() string {
	return "type not supported : " + e.ObjectType
}

func format(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func stringLength(list interface{}) (int, error) {
	tp := reflect.TypeOf(list)
	//nolint:exhaustive // no need to be more exhaustive
	switch tp.Kind() {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)
		return l2.Len(), nil
	case reflect.String:
		return stringLength(list.(string))
	default:
		return 0, &notSupportedTypeError{nil, tp.String()}
	}
}

func FuncMap() map[string]interface{} {
	funcMap := sprig.FuncMap()
	// string functions
	funcMap["len"] = stringLength
	funcMap["format"] = format
	// templates functions
	funcMap["include"] = Include
	funcMap["includeFile"] = includeFile
	funcMap["includeFileAsTemplate"] = includeFileAsTemplate
	funcMap["dynamicFile"] = dynamicFile
	funcMap["removeFirstShebangLineIfAny"] = bash.RemoveFirstShebangLineIfAny
	funcMap["firstCharacterTitle"] = FirstCharacterTitle
	funcMap["snakeCase"] = ToSnakeCase

	return funcMap
}
