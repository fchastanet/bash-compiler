// Package files
package files

import (
	"path/filepath"
)

func MatchFullDirectory(directory string) (files []string, err error) {
	directoryGlob := filepath.Join(directory, "**")
	return MatchPatterns(directoryGlob)
}

// NewAtLevel Initializes logger with provided level
func MatchPatterns(patterns ...string) (files []string, err error) {
	for _, pattern := range patterns {
		patternFiles, err := matchPattern(pattern)
		if err != nil {
			return nil, err
		}
		files = append(files, patternFiles...)
	}
	return files, err
}

func matchPattern(pattern string) (files []string, err error) {
	files, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return files, err
}
