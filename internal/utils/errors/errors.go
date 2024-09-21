package errors

import (
	"fmt"
)

type ValidationError struct {
	InnerError error
	Context    string
	FieldName  string
	FieldValue any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation failed invalid value : context %s field %s value %v",
		e.Context, e.FieldName, e.FieldValue,
	)
}

type closeInterface interface {
	Close() error
}

func SafeCloseDeferCallback(file closeInterface, err *error) {
	// Report the error, if any, from Close, but do so
	// only if there isn't already an outgoing error.
	if c := file.Close(); *err == nil {
		*err = c
	}
}
