package structures

import (
	"errors"
	"fmt"
	"os"
)

var errMissingKey = errors.New("Invalid key")

func ErrMissingKey(key string) error {
	return fmt.Errorf("%w: %s", errMissingKey, key)
}

var errInvalidType = errors.New("Invalid type")

func ErrInvalidType(myVar interface{}) error {
	return fmt.Errorf("%w: %T", errInvalidType, myVar)
}

type Dictionary map[string]interface{}

func (dic Dictionary) GetStringValue(key string) (value string, err error) {
	val, ok := dic[key]
	if !ok {
		return "", ErrMissingKey(key)
	}
	if value, ok := val.(string); ok {
		return ExpandStringValue(value), nil
	}

	return "", ErrInvalidType(value)
}

func ExpandStringValue(value string) string {
	return os.ExpandEnv(value)
}

func (dic Dictionary) GetStringList(key string) (values []string, err error) {
	val, ok := dic[key]
	if !ok {
		return nil, ErrMissingKey(key)
	}
	if values, ok := val.([]string); ok {
		return ExpandStringList(values), nil
	}

	return nil, ErrInvalidType(values)
}

func ExpandStringList(values []string) []string {
	slice := make([]string, len(values))
	for i := len(values) - 1; i >= 0; i-- {
		slice[i] = os.ExpandEnv(values[i])
	}
	return slice
}
