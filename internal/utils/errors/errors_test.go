package errors

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"

	"gotest.tools/v3/assert"
)

var errNotImplemented = errors.New("not implemented")

func TestRenderError(t *testing.T) {
	err := &ValidationError{
		InnerError: &fs.PathError{Op: "", Path: "", Err: errNotImplemented},
		Context:    "compiler",
		FieldName:  "fieldName",
		FieldValue: "fieldValue",
	}

	errStr := fmt.Sprintf("%v", err)
	assert.Equal(
		t,
		"validation failed invalid value : context compiler field fieldName value fieldValue inner error  : not implemented",
		errStr,
	)
}
