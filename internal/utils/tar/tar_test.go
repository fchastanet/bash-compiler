package tar

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/utils/files"
)

func TestCreateArchive(assert *testing.T) {
	var filesList []string
	directory, _ := os.MkdirTemp("", "directory")
	file := filepath.Join(directory, "file1")
	os.Create(file)
	filesList = append(filesList, file)
	subDirectory1 := filepath.Join(directory, "directory")
	os.Mkdir(subDirectory1, files.UserGroupReadExecutePerm)
	file2 := filepath.Join(subDirectory1, "file2")
	os.Create(file2)
	filesList = append(filesList, file2)

	defer os.RemoveAll(directory)

	directoryArchive, _ := os.CreateTemp("", "directoryArchive*.tgz")
	defer os.Remove(directoryArchive.Name())
	err := CreateArchive(
		filesList,
		directory,
		directoryArchive,
		ReproducibleTarOptions,
	)
	if err != nil {
		assert.Errorf("No error should be raised : %q", err)
	}
}
