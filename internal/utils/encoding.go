package utils

import (
	"bytes"
	"crypto/md5" //nolint:all
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func Md5SumFromBytes(buffer []byte) string {
	hash := md5.Sum(buffer) //nolint:all
	return hex.EncodeToString(hash[:])
}

func Md5SumFromFile(file *os.File) (string, error) {
	encoder := md5.New() //nolint:all
	if _, err := io.Copy(encoder, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", encoder.Sum(nil)), nil
}

func Base64FromBytes(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

func Base64FromFile(file *os.File) (string, error) {
	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := io.Copy(encoder, file); err != nil {
		return "", err
	}
	encoder.Close()

	return buf.String(), nil
}
