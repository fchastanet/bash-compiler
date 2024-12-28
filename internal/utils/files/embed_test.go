// Package files
package files

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testsData/*
var templateFs embed.FS

func TestCopyEmbeddedFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "bash-compiler-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = CopyEmbeddedFiles(templateFs, "testsData", tempDir)
	assert.NoError(t, err)

	// Check if files are copied
	entries, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)

	for _, entry := range entries {
		srcPath := filepath.Join("testsData", entry.Name())
		dstPath := filepath.Join(tempDir, entry.Name())

		srcInfo, err := os.Stat(srcPath)
		assert.NoError(t, err)

		dstInfo, err := os.Stat(dstPath)
		assert.NoError(t, err)

		assert.Equal(t, srcInfo.Size(), dstInfo.Size())
	}
}
