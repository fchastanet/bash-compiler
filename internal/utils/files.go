// Package utils
package utils

import (
	"errors"
	"fmt"
	"os"
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
