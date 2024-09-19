package errors

import (
	"fmt"
	"io/fs"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRenderError(t *testing.T) {
	err := &ValidationError{
		InnerError: &fs.PathError{Op: "", Path: "", Err: nil},
		Context:    "compiler",
		FieldName:  "fieldName",
		FieldValue: "fieldValue",
	}

	errStr := fmt.Sprintf("%v", err)
	assert.Equal(
		t,
		"validation failed invalid value : context compiler field fieldName value fieldValue",
		errStr,
	)
}
