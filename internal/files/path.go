package files

import (
	"path/filepath"
	"strings"
)

func BaseNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}
