package files

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestBaseNameWithoutExtensionJustFilename(t *testing.T) {
	files := BaseNameWithoutExtension("files.go")
	assert.Equal(t, files, "files")
}

func TestBaseNameWithoutExtensionWithDir(t *testing.T) {
	files := BaseNameWithoutExtension("internal/myDir/files.go")
	assert.Equal(t, files, "files")
}

func TestSortFilesByPath(t *testing.T) {
	files := []string{
		"dir1/subDir/file",
		"dir2/file2",
		"dir1/file2",
		"dir1/file1",
	}
	expectedFiles := []string{
		"dir1/file1",
		"dir1/file2",
		"dir1/subDir/file",
		"dir2/file2",
	}
	actualFiles := make([]string, len(files))
	copy(actualFiles, files)
	SortFilesByPath(actualFiles)
	assert.DeepEqual(t, actualFiles, expectedFiles)
}
