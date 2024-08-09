package errors

import "fmt"

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
