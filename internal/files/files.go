// Package files
package files

import (
	"crypto/md5" //nolint:golint,gosec
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

var errFilePathMissing = errors.New("File path does not exist")

func ErrFilePathMissing(file string) error {
	return fmt.Errorf("%w: %s", errFilePathMissing, file)
}

func FilePathExists(filePath string) (err error) {
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return ErrFilePathMissing(filePath)
	}
	return nil
}

var errFileMissing = errors.New("File does not exist")

func ErrFileMissing(file string) error {
	return fmt.Errorf("%w: %s", errFileMissing, file)
}

var errFileWasExpected = errors.New("A file was expected")

func ErrFileWasExpected(file string) error {
	return fmt.Errorf("%w: %s", errFileWasExpected, file)
}

var errDirectoryWasExpected = errors.New("A directory was expected")

func ErrDirectoryWasExpected(file string) error {
	return fmt.Errorf("%w: %s", errDirectoryWasExpected, file)
}

func FileExists(filePath string) (err error) {
	stat, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return ErrFileMissing(filePath)
	}
	if stat.IsDir() {
		return ErrFileWasExpected(filePath)
	}
	return nil
}

func DirExists(filePath string) (err error) {
	stat, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return ErrFilePathMissing(filePath)
	}
	if !stat.IsDir() {
		return ErrDirectoryWasExpected(filePath)
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
	defer r.Close() // ignore error: file was opened read-only.

	w, err := os.Create(dstPath)
	if err != nil {
		return err
	}

	defer func() {
		// Report the error, if any, from Close, but do so
		// only if there isn't already an outgoing error.
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

func Md5SumFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New() //nolint:all
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
