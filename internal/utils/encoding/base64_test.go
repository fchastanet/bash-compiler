package encoding

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64FromBytes(t *testing.T) {
	base64 := Base64FromBytes([]byte("Hello"))
	assert.Equal(t, "SGVsbG8=", base64)
}

func TestBase64FromFile(t *testing.T) {
	file, _ := os.Open("testsData/content.txt")
	// skipcq: GO-S2307 // no need Sync as readOnly open
	defer file.Close()
	base64, err := Base64FromFile(file)
	assert.Equal(t, nil, err)
	expectedBase64, err := os.ReadFile("testsData/expectedBase64.txt")
	assert.Equal(t, nil, err)
	assert.Equal(t, strings.Trim(string(expectedBase64), "\n"), base64)
}

func TestBase64FromUnknownFile(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "test******")
	tempFile.Sync()
	tempFile.Close()
	err := os.Remove(tempFile.Name())
	assert.Equal(t, nil, err)
	base64, err := Base64FromFile(tempFile)
	assert.IsType(t, &fs.PathError{Op: "", Path: "", Err: nil}, err)
	assert.Equal(t, "", base64)
}
