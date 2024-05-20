// Package utils
package utils

import (
	"errors"
	"fmt"
	"os"
)

var errFileMissing = errors.New("File does not exist")

func ErrFileMissing(file string) error {
	return fmt.Errorf("%w: %s", errFileMissing, file)
}

func FileExists(filePath string) (err error) {
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return ErrFileMissing(filePath)
	}
	return nil
}
