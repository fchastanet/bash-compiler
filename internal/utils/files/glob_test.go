package files

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMatchFullDirectoryTestsData(t *testing.T) {
	files, err := MatchFullDirectory("testsData")
	expectedFiles := []string{"testsData/testMd5.txt"}
	assert.DeepEqual(t, files, expectedFiles)
	assert.NilError(t, err, "error should have been nil")
}

func MatchFullDirectoryEmptyDir(t *testing.T) {
	files, err := MatchFullDirectory("testsData")
	assert.Equal(t, files, nil)
	assert.NilError(t, err, "error should have been nil")
}
