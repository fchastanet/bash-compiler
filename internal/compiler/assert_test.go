package compiler

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsFunctionDirectiveTrue(t *testing.T) {
	var line = []byte("# FUNCTIONS")
	assert.Equal(t, IsFunctionDirective(line), true)
}

func TestIsFunctionDirectiveNoMatch(t *testing.T) {
	var line = []byte("# FUNC")
	assert.Equal(t, IsFunctionDirective(line), false)
}

func TestIsCommentLineCase1(t *testing.T) {
	var line = []byte("# FUNCTIONS")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineCase2(t *testing.T) {
	var line = []byte("### comment")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineCase3(t *testing.T) {
	var line = []byte(" \t# comment with spaces before")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineNoMatch(t *testing.T) {
	var line = []byte("cmd # FUNC")
	assert.Equal(t, IsCommentLine(line), false)
}

func TestIsBashFrameworkFunctionNoMatch1(t *testing.T) {
	var line = []byte("TEST")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatch2(t *testing.T) {
	var line = []byte("Log:fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchComments1(t *testing.T) {
	var line = []byte("# Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchComments2(t *testing.T) {
	var line = []byte("  \t # Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchAccents(t *testing.T) {
	var line = []byte("Log::fatal::François")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionValidSimple(t *testing.T) {
	var line = []byte("Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), true)
}

func TestIsBashFrameworkFunctionValidMultiple(t *testing.T) {
	var line = []byte("Namespace1::Namespace2::function")
	assert.Equal(t, IsBashFrameworkFunction(line), true)
}
