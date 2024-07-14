package files

import (
	"path/filepath"
	"sort"
	"strings"
)

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
