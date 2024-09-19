package bash

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestRemoveFirstShebangLineIfAny(t *testing.T) {
	t.Run("one line", func(t *testing.T) {
		oldCode := "#!/bin/bash"
		newCode := RemoveFirstShebangLineIfAny(oldCode)
		assert.Equal(t, "\n", newCode)
	})

	t.Run("only comment", func(t *testing.T) {
		oldCode := "#/bin/bash"
		newCode := RemoveFirstShebangLineIfAny(oldCode)
		assert.Equal(t, "#/bin/bash\n", newCode)
	})

	t.Run("different kind of shebangs", func(t *testing.T) {
		oldCode, err := os.ReadFile("testsData/multipleShebangs.sh")
		assert.NilError(t, err)
		newCode := RemoveFirstShebangLineIfAny(string(oldCode))

		expectedCode, err := os.ReadFile("testsData/multipleShebangs.expected.txt")
		assert.NilError(t, err)
		assert.Equal(t, string(expectedCode), newCode)
	})
}
