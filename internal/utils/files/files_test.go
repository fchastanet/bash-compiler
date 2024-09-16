package files

import (
	"testing"

	"gotest.tools/v3/assert"
)

const (
	expectedFileFilesGo  = "files.go"
	expectedFileFiles2Go = "files2.go"
)

func TestFilePathExists(t *testing.T) {
	assert.Equal(t, FilePathExists(expectedFileFilesGo), nil)
}

func TestFilePathNotExists(t *testing.T) {
	assert.Error(t, FilePathExists(expectedFileFiles2Go), "file path does not exist: files2.go")
}

func TestFilePathDirExists(t *testing.T) {
	assert.Equal(t, FilePathExists(".."), nil)
}

func TestFileExists(t *testing.T) {
	assert.Equal(t, FileExists(expectedFileFilesGo), nil)
}

func TestFileExistsButDir(t *testing.T) {
	assert.Error(t, FileExists("testsData"), "a file was expected: testsData")
}

func TestFileNotExists(t *testing.T) {
	assert.Error(t, FileExists(expectedFileFiles2Go), "file path does not exist: files2.go")
}

func TestDirExists(t *testing.T) {
	assert.Equal(t, DirExists(".."), nil)
}

func TestDirNotExists(t *testing.T) {
	assert.Error(t, DirExists("dirNotExists"), "directory path does not exist: dirNotExists")
}

func TestDirExistsNotADirectory(t *testing.T) {
	assert.Error(t, DirExists(expectedFileFilesGo), "a directory was expected: files.go")
}

func TestSha256FromFileButDir(t *testing.T) {
	sha256, err := ChecksumFromFile("testsData")
	assert.Equal(t, sha256, "")
	assert.Error(t, err, "read testsData: is a directory")
}

func TestSha256FromFileOk(t *testing.T) {
	sha256, err := ChecksumFromFile("testsData/testMd5.txt")
	assert.Equal(t, sha256, "af99a79c936e4625c10bc2d3b9e4adf14a67f2d8a1ae27453a77fc5a59bb1b4b")
	assert.NilError(t, err, "error should have been nil")
}

func TestSha256FromFileNotExists(t *testing.T) {
	sha256, err := ChecksumFromFile("testsData/notExists.txt")
	assert.Equal(t, sha256, "")
	assert.Error(t, err, "open testsData/notExists.txt: no such file or directory")
}
