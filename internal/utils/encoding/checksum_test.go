package encoding

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumFromBytes(t *testing.T) {
	checksum := ChecksumFromBytes([]byte("Hello"))
	assert.Equal(
		t,
		"48656c6c6fe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		checksum,
	)
}

func TestChecksumFromFile(t *testing.T) {
	file, _ := os.Open("testsData/content.txt")
	defer file.Close()
	checksum, err := ChecksumFromFile(file)
	assert.Equal(t, nil, err)
	assert.Equal(
		t,
		"b2fc955e752d3ea39c4ee61a1725252c95a83467fb4779c10f7078b6836020e2",
		checksum,
	)
}

func TestChecksumFromUnknownFile(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "test******")
	tempFile.Close()
	err := os.Remove(tempFile.Name())
	assert.Equal(t, nil, err)
	checksum, err := ChecksumFromFile(tempFile)
	assert.IsType(t, &fs.PathError{Op: "", Path: "", Err: nil}, err)
	assert.Equal(t, "", checksum)
}
