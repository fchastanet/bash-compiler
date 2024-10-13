package files

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type dirError struct {
	error
	Path   string
	Reason string
}

func (err *dirError) Error() string {
	return fmt.Sprintf("Directory %s %s", err.Path, err.Reason)
}

func BaseNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func SortFilesByPath(files []string) {
	sort.Slice(files, func(i, j int) bool {
		dirI, fileI := filepath.Split(files[i])
		dirJ, fileJ := filepath.Split(files[j])

		// First sort by directory path
		if dirI != dirJ {
			return dirI < dirJ
		}
		// If directories are the same, sort by file name
		return fileI < fileJ
	})
}

func IsWritableDirectory(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return &dirError{nil, path, "is not a directory"}
	}
	if stat.Mode()&0o600 != 0 {
		return &dirError{nil, path, "is not writable"}
	}
	return nil
}
