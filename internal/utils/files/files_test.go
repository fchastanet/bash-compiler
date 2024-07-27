package files

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestFilePathExists(t *testing.T) {
	assert.Equal(t, FilePathExists("files.go"), nil)
}

func TestFilePathNotExists(t *testing.T) {
	assert.Error(t, FilePathExists("files2.go"), "File path does not exist: files2.go")
}

func TestFilePathDirExists(t *testing.T) {
	assert.Equal(t, FilePathExists(".."), nil)
}

func TestFileExists(t *testing.T) {
	assert.Equal(t, FileExists("files.go"), nil)
}

func TestFileExistsButDir(t *testing.T) {
	assert.Error(t, FileExists("testsData"), "A file was expected: testsData")
}

func TestFileNotExists(t *testing.T) {
	assert.Error(t, FileExists("files2.go"), "File does not exist: files2.go")
}

func TestDirExists(t *testing.T) {
	assert.Equal(t, DirExists(".."), nil)
}

func TestDirNotExists(t *testing.T) {
	assert.Error(t, DirExists("dirNotExists"), "Directory path does not exist: dirNotExists")
}

func TestDirExistsNotADirectory(t *testing.T) {
	assert.Error(t, DirExists("files.go"), "A directory was expected: files.go")
}

func TestMd5FromFileButDir(t *testing.T) {
	md5, err := Md5SumFromFile("testsData")
	assert.Equal(t, md5, "")
	assert.Error(t, err, "read testsData: is a directory")
}

func TestMd5FromFileOk(t *testing.T) {
	md5, err := Md5SumFromFile("testsData/testMd5.txt")
	assert.Equal(t, md5, "772ac1a55fab1122f3b369ee9cd31549")
	assert.NilError(t, err, "error should have been nil")
}

func TestMd5FromFileNotExists(t *testing.T) {
	md5, err := Md5SumFromFile("testsData/notExists.txt")
	assert.Equal(t, md5, "")
	assert.Error(t, err, "open testsData/notExists.txt: no such file or directory")
}
