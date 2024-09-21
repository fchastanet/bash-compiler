package dotenv

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

const (
	msgErrorShouldBeNil = "error should be nil"
)

func setupTest(_ testing.TB) func(tb testing.TB) {
	os.Clearenv()
	return func(_ testing.TB) {
		os.Clearenv()
	}
}

func TestLoadSimpleFileMyVarNotExisting(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	err := LoadEnvFile("./testsData/simpleFile.txt")
	assert.NilError(t, err, msgErrorShouldBeNil)
	envValue, exists := os.LookupEnv("FRAMEWORK_ROOT_DIR")
	assert.Equal(t, exists, true)
	assert.Equal(t, envValue, "", envValue)
}

func TestLoadSimpleFileMyVarExists(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	t.Setenv("MY_VAR", "myValue")
	err := LoadEnvFile("./testsData/simpleFile.txt")
	assert.NilError(t, err, msgErrorShouldBeNil)
	envValue, exists := os.LookupEnv("FRAMEWORK_ROOT_DIR")
	assert.Equal(t, exists, true)
	assert.Equal(t, envValue, "myValue", envValue)
}

func TestLoadFileVarsDefaultValues(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	err := LoadEnvFile("./testsData/vars.txt")
	assert.NilError(t, err, msgErrorShouldBeNil)
	envValue, exists := os.LookupEnv("FRAMEWORK_ROOT_DIR")
	assert.Equal(t, exists, true)
	assert.Equal(t, envValue, "dummy", envValue)
}

func TestLoadFileDependentVarsDefaultValues(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	err := LoadEnvFile("./testsData/dependentVars.txt")
	assert.NilError(t, err, "error should be nil")
	envValue, exists := os.LookupEnv("FRAMEWORK_ROOT_DIR")
	assert.Equal(t, exists, true)
	assert.Equal(t, envValue, "dummy", envValue)
	envValue, exists = os.LookupEnv("SRC_FILE")
	assert.Equal(t, exists, true)
	assert.Equal(t, envValue, "dummy/srcFile", envValue)
}
