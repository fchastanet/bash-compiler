package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"gotest.tools/v3/assert"
)

func getDefaultExpectedCli(expectedCli *cli) error {
	var expectedYamlFiles []string
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	expectedCli.YamlFiles = YamlFiles(expectedYamlFiles)

	expectedCli.RootDirectory = RootDirectory(currentDir)
	expectedCli.BinaryFilesExtension = BinaryFilesExtension("-binary.yaml")
	expectedCli.Version = VersionFlag("")
	expectedCli.IntermediateFilesDir = IntermediateFilesDir("")
	expectedCli.Debug = false
	expectedCli.LogLevel = int(slog.LevelInfo)
	return nil
}

func defaultCase(t *testing.T, args []string) {
	os.Args = args
	expectedCli := &cli{} //nolint:exhaustruct //test
	err := getDefaultExpectedCli(expectedCli)
	assert.NilError(t, err)
	cli := &cli{} //nolint:exhaustruct //test
	err = parseArgs(cli)
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedCli, cli)
}

func TestArgs(t *testing.T) {
	currentDir, err := os.Getwd()
	assert.NilError(t, err)
	testsDataDir := filepath.Join(currentDir, "testsData")
	err = os.Chdir(testsDataDir)
	assert.NilError(t, err)

	t.Run("no arg", func(t *testing.T) {
		defaultCase(t, []string{"cmd"})
	})

	t.Run("rootDir provided (go run mode simulated)", func(t *testing.T) {
		defaultCase(t, []string{"__debug_bin_cmd", "-r", testsDataDir})
	})

	t.Run("root does not exist", func(t *testing.T) {
		os.Args = []string{"cmd", "-r", "inexistent"}
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.ErrorContains(t, err, "rootDir option should not be provided")
	})

	t.Run("no arg (go run mode simulated)", func(t *testing.T) {
		os.Args = []string{"__debug_bin_cmd"}
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.ErrorContains(t, err, "please provide rootDir option")
	})

	t.Run("target dir", func(t *testing.T) {
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.NilError(t, err)
		expectedTargetDir := string(expectedCli.RootDirectory)
		os.Args = []string{"cmd", "-t", expectedTargetDir}
		expectedCli.IntermediateFilesDir = IntermediateFilesDir(expectedTargetDir)
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.NilError(t, err)
		assert.DeepEqual(t, expectedCli, cli)
	})

	t.Run("debug", func(t *testing.T) {
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.NilError(t, err)
		expectedCli.Debug = true
		expectedCli.LogLevel = int(slog.LevelDebug)
		os.Args = []string{"cmd", "-d"}
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.NilError(t, err)
		assert.DeepEqual(t, expectedCli, cli)
	})

	t.Run("yaml file", func(t *testing.T) {
		os.Args = []string{"cmd", "file-binary.yaml"}
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.NilError(t, err)
		expectedCli.YamlFiles = append(
			expectedCli.YamlFiles,
			filepath.Join(string(expectedCli.RootDirectory), "file-binary.yaml"),
		)
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.NilError(t, err)
		assert.DeepEqual(t, expectedCli, cli)
	})

	err = os.Chdir(currentDir)
	assert.NilError(t, err)
}
