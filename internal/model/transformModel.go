package model

import (
	"bytes"
	_ "embed"
	"log"
	"os"

	sidekick "cuelang.org/go/cmd/cue/cmd"
)

//go:embed binFile.cue
var binFileCueSchema string

func TransformModel(tempYamlFile os.File, resultWriter *bytes.Buffer) (err error) {
	// write cue file to temp file
	tempCueFile, err := os.CreateTemp("", "binFile*.cue")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempCueFile.Name())
	_, err = tempCueFile.Write([]byte(binFileCueSchema))
	if err != nil {
		return err
	}

	log.Printf("Temp file containing cue file : %s\n", tempCueFile.Name())

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
