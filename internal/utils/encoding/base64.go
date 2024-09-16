package encoding

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
)

func Base64FromBytes(buffer []byte) string {
	return base64.StdEncoding.EncodeToString(buffer)
}

func Base64FromFile(file *os.File) (string, error) {
	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := io.Copy(encoder, file); err != nil {
		return "", err
	}
	err := encoder.Close()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
