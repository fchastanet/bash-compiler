package encoding

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func ChecksumFromBytes(buffer []byte) string {
	encoder := sha256.New()
	hash := encoder.Sum(buffer)
	return hex.EncodeToString(hash)
}

func ChecksumFromFile(file *os.File) (string, error) {
	encoder := sha256.New()
	if _, err := io.Copy(encoder, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(encoder.Sum(nil)), nil
}
