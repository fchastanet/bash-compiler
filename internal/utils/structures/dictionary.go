package structures

import (
	"fmt"
	"os"
)

type missingKeyError struct {
	error
	Key string
}

func (e *missingKeyError) Error() string {
	return "missing key: " + e.Key
}

type invalidValueTypeError struct {
	error
	Value any
}

func (e *invalidValueTypeError) Error() string {
	return fmt.Sprintf("invalid type: %v", e.Value)
}

type Dictionary map[string]any

func (dic Dictionary) GetStringValue(key string) (value string, err error) {
	val, ok := dic[key]
	if !ok {
		return "", &missingKeyError{nil, key}
	}
	if value, ok := val.(string); ok {
		return ExpandStringValue(value), nil
	}

	return "", &invalidValueTypeError{nil, val}
}

func ExpandStringValue(value string) string {
	return os.ExpandEnv(value)
}

func (dic Dictionary) GetStringList(key string) (values []string, err error) {
	val, ok := dic[key]
	if !ok {
		return nil, &missingKeyError{nil, key}
	}
	if values, ok := val.([]string); ok {
		return ExpandStringList(values), nil
	}

	return nil, &invalidValueTypeError{nil, val}
}

func ExpandStringList(values []string) []string {
	slice := make([]string, len(values))
	for i := len(values) - 1; i >= 0; i-- {
		slice[i] = os.ExpandEnv(values[i])
	}

	return slice
}
