// Package files
package files

import (
	"embed"
	"os"
	"path/filepath"
)

func CopyEmbeddedFiles(templateFs embed.FS, subDir string, tempDir string) error {
	return copyDir(templateFs, subDir, tempDir)
}

func copyDir(templateFs embed.FS, srcDir, dstDir string) error {
	entries, err := templateFs.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(templateFs, srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(templateFs, srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyDirectory(templateFs embed.FS, srcPath, dstPath string) error {
	if err := os.MkdirAll(dstPath, AllReadExecutePerm); err != nil {
		return err
	}
	return copyDir(templateFs, srcPath, dstPath)
}

func copyFile(templateFs embed.FS, srcPath, dstPath string) error {
	data, err := templateFs.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, AllReadPerm)
}
