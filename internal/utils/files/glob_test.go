package files

import (
	"os"
	"path"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMatchFullDirectoryRelativeTestsData(t *testing.T) {
	files, err := MatchFullDirectoryRelative("testsData")
	expectedFiles := []string{"testMd5.txt"}
	assert.DeepEqual(t, files, expectedFiles)
	assert.NilError(t, err, "error should have been nil")
}

func TestMatchFullDirectoryRelativeEmptyDir(t *testing.T) {
	emptyDir, _ := os.MkdirTemp("/tmp", "emptyDir*")

	t.Run("empty directory", func(t *testing.T) {
		files, err := MatchFullDirectoryRelative(emptyDir)
		var expectedFiles []string
		assert.DeepEqual(t, files, expectedFiles)
		assert.NilError(t, err, "error should have been nil")
	})

	t.Run("directory with dot files", func(t *testing.T) {
		os.Create(path.Join(emptyDir, ".file"))
		files, err := MatchFullDirectoryRelative(emptyDir)
		assert.DeepEqual(t, files, []string{".file"})
		assert.NilError(t, err, "error should have been nil")
	})

	t.Run("directory with several directories", func(t *testing.T) {
		os.Create(path.Join(emptyDir, ".file1"))
		os.Mkdir(path.Join(emptyDir, "dir"), os.ModePerm)
		os.Create(path.Join(emptyDir, "dir", "file2"))
		files, err := MatchFullDirectoryRelative(emptyDir)
		assert.DeepEqual(t, files, []string{
			".file", ".file1", "dir/file2",
		})
		assert.NilError(t, err, "error should have been nil")
	})

	os.Remove(emptyDir)
}
