package model

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os"

	sidekick "cuelang.org/go/cmd/cue/cmd"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

//go:embed binFile.cue
var binFileCueSchema string

func transformModel(tempYamlFile os.File, resultWriter *bytes.Buffer) (err error) {
	// write cue file to temp file
	tempCueFile, err := os.CreateTemp("", "binFile*.cue")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempCueFile.Name())
	_, err = tempCueFile.WriteString(binFileCueSchema)
	if err != nil {
		return err
	}

	slog.Debug("Temp file containing cue file", logger.LogFieldFilePath, tempCueFile.Name())
	// transform using cue
	cmd, err := sidekick.New([]string{
		"export",
		"-l", "input:", tempYamlFile.Name(),
		tempCueFile.Name(),
		"--out", "yaml", "-e", "output",
	})
	if err != nil {
		return err
	}

	// outputs result
	cmd.SetOutput(resultWriter)
	err = cmd.Run(cmd.Context())

	return err
}
