// Package files
package files

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func MatchFullDirectoryRelative(directory string) (files []string, err error) {
	return MatchPatternRelative(directory, "**/*")
}

func MatchPatternRelative(directory string, pattern string) (files []string, err error) {
	return doublestar.Glob(
		os.DirFS(directory),
		pattern,
		doublestar.WithFailOnIOErrors(),
		doublestar.WithFilesOnly(),
		doublestar.WithNoFollow(),
	)
}

func MatchPatterns(directory string, patterns ...string) (files []string, err error) {
	filesList := []string{}
	for _, pattern := range patterns {
		list, err := matchPattern(directory, pattern)
		if err != nil {
			return []string{}, err
		}
		filesList = append(filesList, list...)
	}
	return filesList, nil
}

func matchPattern(directory string, pattern string) (files []string, err error) {
	return doublestar.FilepathGlob(
		filepath.Join(directory, pattern),
		doublestar.WithFailOnIOErrors(),
		doublestar.WithFilesOnly(),
		doublestar.WithNoFollow(),
	)
}
