package structures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeys(t *testing.T) {
	var emptyMap map[string]string
	oneElemMap := map[string]string{"id1": "elem1"}
	intMap := map[int]string{1: "elem1", 2: "elem2"}

	t.Run("empty map", func(t *testing.T) {
		list := MapKeys(emptyMap)
		// skipcq: GO-W1027 // here not nil expected but real slice
		expectedList := []string{}
		assert.Equal(t, list, expectedList)
	})

	t.Run("one element", func(t *testing.T) {
		list := MapKeys(oneElemMap)
		expectedList := []string{"id1"}
		assert.Equal(t, list, expectedList)
	})

	t.Run("map with int keys", func(t *testing.T) {
		list := MapKeys(intMap)
		assert.IsType(t, []int{}, list)
		assert.Contains(t, list, 1)
		assert.Contains(t, list, 2)
	})
}
