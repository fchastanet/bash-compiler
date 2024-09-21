package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

func CreateArchive(
	files []string,
	relativeDir string,
	buf io.Writer,
	updateFileInfoHeader func(info *tar.Header, fi fs.FileInfo) error,
) (err error) {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer errors.SafeCloseDeferCallback(gw, &err)
	tw := tar.NewWriter(gw)
	defer errors.SafeCloseDeferCallback(tw, &err)

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

func getFileHeader(
	fileInfo fs.FileInfo,
	filename string,
	relativeDir string,
) (header *tar.Header, err error) {
	// Create a tar Header from the FileInfo data
	header, err = tar.FileInfoHeader(fileInfo, filename)
	if err != nil {
		return nil, err
	}
	// relative paths are used to preserve the directory paths in each file path
	if filepath.IsAbs(filename) {
		relativePath, err := filepath.Rel(relativeDir, filename)
		if err != nil {
			return nil, err
		}
		header.Name = relativePath
	} else {
		header.Name = filename
	}
	return header, nil
}

func addToArchive(
	tw *tar.Writer,
	filename string,
	relativeDir string,
	updateFileInfoHeader func(info *tar.Header, fi fs.FileInfo) error,
) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		// skip directories
		return nil
	}
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer errors.SafeCloseDeferCallback(file, &err)

	header, err := getFileHeader(fileInfo, filename, relativeDir)
	if err != nil {
		return err
	}
	err = updateFileInfoHeader(header, fileInfo)
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
