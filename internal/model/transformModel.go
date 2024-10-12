package model

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os"
	"path"

	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/goccy/go-yaml"
	"kcl-lang.io/kcl-go"
)

//go:embed kcl/binFile.k
var kclBinFileSchema string

//go:embed kcl/libs.k
var kclLibs string

func transformModel(tempYamlFile os.File, resultWriter *bytes.Buffer) (err error) {
	tempKclTempDir, err := os.MkdirTemp("", "kcl")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempKclTempDir)

	// write k files to temp files
	tempKclFilePath := path.Join(tempKclTempDir, "binFile.k")
	err = os.WriteFile(tempKclFilePath, []byte(kclBinFileSchema), files.UserReadWriteExecutePerm)
	if err != nil {
		return err
	}
	slog.Debug("Temp file containing binFile.k file", logger.LogFieldFilePath, tempKclFilePath)

	tempKclLibsFilePath := path.Join(tempKclTempDir, "libs.k")
	err = os.WriteFile(tempKclLibsFilePath, []byte(kclLibs), files.UserReadWriteExecutePerm)
	if err != nil {
		return err
	}
	slog.Debug("Temp file containing libs.k file", logger.LogFieldFilePath, tempKclLibsFilePath)

	// Run the KCL script
	result, err := kcl.Run(
		tempKclFilePath,
		kcl.WithOptions(
			"-D", "configFile="+tempYamlFile.Name(),
			"-S", "configYaml",
		),
		kcl.WithSortKeys(true),
	)
	if err != nil {
		return err
	}

	// Check the result
	first := result.First()
	configYaml, err := first.ToMap()
	if err != nil {
		return err
	}
	yamlResult, err := yaml.Marshal(configYaml["configYaml"])
	if err != nil {
		return err
	}
	_, err = resultWriter.Write(yamlResult)

	return err
}
