package encoding

import (
	"crypto/md5" //nolint:all
	"encoding/hex"
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

	return hex.EncodeToString(encoder.Sum(nil)), nil
}
