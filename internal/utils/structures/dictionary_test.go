package structures

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetStringValue(t *testing.T) {
	dic := Dictionary{}

	t.Run("missing key", func(t *testing.T) {
		val, err := dic.GetStringValue("missingKey")
		assert.Equal(t, val, "")
		assert.Error(t, err, "missing key: missingKey")
	})

	t.Run("invalid value", func(t *testing.T) {
		dic["key"] = 1
		val, err := dic.GetStringValue("key")
		assert.Equal(t, val, "")
		assert.Error(t, err, "invalid type: 1")
	})

	t.Run("valid value", func(t *testing.T) {
		dic["key"] = "valid"
		val, err := dic.GetStringValue("key")
		assert.Equal(t, val, "valid")
		assert.Equal(t, err, nil)
	})

	t.Run("expand env", func(t *testing.T) {
		t.Setenv("ENV_KEY", "envValue")
		dic["key"] = "${ENV_KEY}"
		val, err := dic.GetStringValue("key")
		assert.Equal(t, val, "envValue")
		assert.Equal(t, err, nil)
	})
}

func TestStringList(t *testing.T) {
	dic := Dictionary{}

	t.Run("missing key", func(t *testing.T) {
		val, err := dic.GetStringList("missingKey")
		var expectedList []string
		assert.DeepEqual(t, val, expectedList)
		assert.Error(t, err, "missing key: missingKey")
	})

	t.Run("invalid value", func(t *testing.T) {
		dic["key"] = 1
		val, err := dic.GetStringList("key")
		var expectedList []string
		assert.DeepEqual(t, val, expectedList)
		assert.Error(t, err, "invalid type: 1")
	})

	t.Run("valid value", func(t *testing.T) {
		expectedList := []string{"valid"}
		dic["key"] = expectedList
		list, err := dic.GetStringList("key")
		assert.DeepEqual(t, list, expectedList)
		assert.Equal(t, err, nil)
	})

	t.Run("expand env", func(t *testing.T) {
		t.Setenv("ENV_KEY", "envValue")
		expectedList := []string{"envValue"}
		dic["key"] = []string{"${ENV_KEY}"}
		list, err := dic.GetStringList("key")
		assert.DeepEqual(t, list, expectedList)
		assert.Equal(t, err, nil)
	})
}
