package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fchastanet/bash-compiler/internal/logger"
)

func CreateArchive(
	files []string,
	relativeDir string,
	buf io.Writer,
	updateFileInfoHeader func(info *tar.Header, fi fs.FileInfo) error,
) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToArchive(
			tw,
			file,
			relativeDir,
			updateFileInfoHeader,
		)
		if logger.FancyHandleError(err) {
			return err
		}
	}

	return nil
}

// https://stackoverflow.com/a/77679184/3045926
func ReproducibleTarOptions(info *tar.Header, _ fs.FileInfo) error {
	info.Uid = 0
	info.Gid = 0
	info.Uname = ""
	info.Gname = ""
	info.ModTime = time.Time{}
	info.AccessTime = time.Time{}
	info.ChangeTime = time.Time{}
	return nil
}

func addToArchive(
	tw *tar.Writer,
	filename string,
	relativeDir string,
	updateFileInfoHeader func(info *tar.Header, fi fs.FileInfo) error,
) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	// relative paths are used to preserve the directory paths in each file path
	relativePath, err := filepath.Rel(relativeDir, filename)
	if err != nil {
		return err
	}
	header.Name = relativePath

	err = updateFileInfoHeader(header, info)
	if err != nil {
		return err
	}

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
