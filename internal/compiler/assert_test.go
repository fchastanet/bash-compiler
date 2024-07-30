package compiler

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsFunctionDirectiveTrue(t *testing.T) {
	line := []byte("# FUNCTIONS")
	assert.Equal(t, IsFunctionDirective(line), true)
}

func TestIsFunctionDirectiveNoMatch(t *testing.T) {
	line := []byte("# FUNC")
	assert.Equal(t, IsFunctionDirective(line), false)
}

func TestIsCommentLineCase1(t *testing.T) {
	line := []byte("# FUNCTIONS")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineCase2(t *testing.T) {
	line := []byte("### comment")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineCase3(t *testing.T) {
	line := []byte(" \t# comment with spaces before")
	assert.Equal(t, IsCommentLine(line), true)
}

func TestIsCommentLineNoMatch(t *testing.T) {
	line := []byte("cmd # FUNC")
	assert.Equal(t, IsCommentLine(line), false)
}

func TestIsBashFrameworkFunctionNoMatch1(t *testing.T) {
	line := []byte("TEST")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatch2(t *testing.T) {
	line := []byte("Log:fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchComments1(t *testing.T) {
	line := []byte("# Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchComments2(t *testing.T) {
	line := []byte("  \t # Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionNoMatchAccents(t *testing.T) {
	line := []byte("Log::fatal::François")
	assert.Equal(t, IsBashFrameworkFunction(line), false)
}

func TestIsBashFrameworkFunctionValidSimple(t *testing.T) {
	line := []byte("Log::fatal")
	assert.Equal(t, IsBashFrameworkFunction(line), true)
}

func TestIsBashFrameworkFunctionValidMultiple(t *testing.T) {
	line := []byte("Namespace1::Namespace2::function")
	assert.Equal(t, IsBashFrameworkFunction(line), true)
}
