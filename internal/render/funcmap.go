// Package render
package render

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

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

func stringLength(list any) (int, error) {
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

func sortCallbacks(list []any) (any, error) {
	sort.SliceStable(list, func(i int, j int) bool {
		priorityA := i
		if valA, okA := list[i].(string); okA {
			priorityA = arrayDefaultValue(strings.Split(valA, "@"), 1, i)
		}
		priorityB := j
		if valB, okB := list[j].(string); okB {
			priorityB = arrayDefaultValue(strings.Split(valB, "@"), 1, j)
		}
		return priorityA < priorityB
	})
	list2 := make([]string, len(list))
	for _, elem := range list {
		file, ok := elem.(string)
		if !ok {
			continue
		}
		split := strings.Split(file, "@")
		list2 = append(list2, split[0])
	}
	return list2, nil
}

func arrayDefaultValue(list []string, i int, defaultValue int) int {
	if len(list) >= i+1 {
		if number, err := strconv.Atoi(list[i]); err == nil {
			return number
		}
	}
	return defaultValue
}

func FuncMap() map[string]any {
	funcMap := sprig.FuncMap()
	// string functions
	funcMap["len"] = stringLength
	funcMap["format"] = format
	// callbacks
	funcMap["sortCallbacks"] = sortCallbacks
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
