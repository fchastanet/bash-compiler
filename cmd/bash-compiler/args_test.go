package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDefaultExpectedCli(expectedCli *cli) error {
	var expectedYamlFiles []string
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	targetDir, err := filepath.Abs(currentDir + "/../..")
	if err != nil {
		return err
	}
	expectedCli.YamlFiles = YamlFiles(expectedYamlFiles)

	expectedCli.RootDirectory = Directory(filepath.Join(currentDir, "testsData"))
	expectedCli.TargetDir = Directory(targetDir)
	expectedCli.BinaryFilesExtension = BinaryFilesExtension("-binary.yaml")
	expectedCli.Version = VersionFlag("")
	expectedCli.KeepIntermediateFiles = false
	expectedCli.Debug = false
	expectedCli.LogLevel = int(slog.LevelInfo)
	expectedCli.CompilerRootDir = Directory(targetDir)
	return nil
}

func TestArgs(t *testing.T) {
	t.Run("no arg except mandatory", func(t *testing.T) {
		os.Args = []string{"cmd", "-r", "testsData"}
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.Equal(t, nil, err)
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.Nil(t, err)
		assert.Equal(t, expectedCli, cli)
	})

	t.Run("target dir", func(t *testing.T) {
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.Equal(t, nil, err)
		expectedTargetDir := string(expectedCli.RootDirectory)
		os.Args = []string{"cmd", "-r", "testsData", "-t", expectedTargetDir}
		expectedCli.TargetDir = Directory(expectedTargetDir)
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.Nil(t, err)
		assert.Equal(t, expectedCli, cli)
	})

	t.Run("debug", func(t *testing.T) {
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.Equal(t, nil, err)
		expectedCli.Debug = true
		expectedCli.LogLevel = int(slog.LevelDebug)
		os.Args = []string{"cmd", "-r", "testsData", "-d"}
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.Nil(t, err)
		assert.Equal(t, expectedCli, cli)
	})

	t.Run("yaml file", func(t *testing.T) {
		os.Args = []string{"cmd", "-r", "testsData", "testsData/file-binary.yaml"}
		expectedCli := &cli{} //nolint:exhaustruct //test
		err := getDefaultExpectedCli(expectedCli)
		assert.Equal(t, nil, err)
		expectedCli.YamlFiles = append(
			expectedCli.YamlFiles,
			filepath.Join(string(expectedCli.RootDirectory), "file-binary.yaml"),
		)
		cli := &cli{} //nolint:exhaustruct //test
		err = parseArgs(cli)
		assert.Nil(t, err)
		assert.Equal(t, expectedCli, cli)
	})
}
