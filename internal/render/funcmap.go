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

func sortByKeys(myMap map[string]any) (list []any, err error) {
	keys := make([]int, 0, len(myMap))
	list = make([]any, 0, len(myMap))

	for k := range myMap {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}
		keys = append(keys, kInt)
	}
	sort.Ints(keys)

	for _, k := range keys {
		list = append(list, myMap[strconv.Itoa(k)])
	}
	return list, nil
}

func arrayDefaultValue(list []string, i int, defaultValue int) int {
	if len(list) >= i+1 {
		if number, err := strconv.Atoi(list[i]); err == nil {
			return number
		}
	}
	return defaultValue
}

// chunkBase64 splits a base64 string into chunks of 76 characters with bash line continuation
func chunkBase64(s string) string {
	chunkSize := 76 // Standard base64 line length

	var result strings.Builder
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}

		// Add the chunk
		result.WriteString(s[i:end])

		// Add bash line continuation if not the last chunk
		if end < len(s) {
			result.WriteString(" \\\n")
		}
	}

	return result.String()
}

func FuncMap() map[string]any {
	funcMap := sprig.FuncMap()
	// string functions
	funcMap["len"] = stringLength
	funcMap["format"] = format
	funcMap["chunkBase64"] = chunkBase64
	// callbacks
	funcMap["sortCallbacks"] = sortCallbacks
	funcMap["sortByKeys"] = sortByKeys
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
