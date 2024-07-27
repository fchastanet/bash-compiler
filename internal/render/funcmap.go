// Package render
package render

import (
	"errors"
	"fmt"
	"reflect"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/fchastanet/bash-compiler/internal/utils/bash"
)

var errorNotSupportedType = errors.New("Type not supported")

func format(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func stringLength(list interface{}) (int, error) {
	tp := reflect.TypeOf(list).Kind()
	//nolint:exhaustive
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)
		return l2.Len(), nil
	case reflect.String:
		return stringLength(list.(string))
	default:
		return 0, errorNotSupportedType
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
