// Package files
package files

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"

	myErrors "github.com/fchastanet/bash-compiler/internal/utils/errors"
)

const (
	UserReadWritePerm        os.FileMode = 0o600
	UserReadWriteExecutePerm os.FileMode = 0o700
	AllReadExecutePerm       os.FileMode = 0o755
	AllReadPerm              os.FileMode = 0o644
)

type filePathMissingError struct {
	error
	FilePath string
}

func (e *filePathMissingError) Error() string {
	return "file path does not exist: " + e.FilePath
}

func FilePathExists(filePath string) (err error) {
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return &filePathMissingError{nil, filePath}
	}
	return nil
}

type fileWasExpectedError struct {
	error
	File string
}

func (e *fileWasExpectedError) Error() string {
	return "a file was expected: " + e.File
}

type directoryWasExpectedError struct {
	error
	Directory string
}

func (e *directoryWasExpectedError) Error() string {
	return "a directory was expected: " + e.Directory
}

type directoryPathMissingError struct {
	error
	DirPath string
}

func (e *directoryPathMissingError) Error() string {
	return "directory path does not exist: " + e.DirPath
}

func FileExists(filePath string) (err error) {
	stat, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return &filePathMissingError{nil, filePath}
	}
	if stat.IsDir() {
		return &fileWasExpectedError{err, filePath}
	}
	return nil
}

func DirExists(filePath string) (err error) {
	stat, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return &directoryPathMissingError{err, filePath}
	}
	if !stat.IsDir() {
		return &directoryWasExpectedError{nil, filePath}
	}
	return nil
}

// Copy copies the contents of the file at srcPath to a regular file
// at dstPath. If the file named by dstPath already exists, it is
// truncated. The function does not copy the file mode, file
// permission bits, or file attributes.
func Copy(srcPath string, dstPath string) (err error) {
	r, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer myErrors.SafeCloseDeferCallback(r, &err)

	w, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer myErrors.SafeCloseDeferCallback(w, &err)

	_, err = io.Copy(w, r)
	return err
}

func ChecksumFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	// skipcq: GO-S2307 // no need Sync as readOnly open
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
