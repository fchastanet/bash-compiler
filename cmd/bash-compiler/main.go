// main package
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/log"
	"go.uber.org/automaxprocs/maxprocs"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

type cli struct {
	YamlFile YamlFile `arg:"" help:"Yaml file" type:"path"`
}
type YamlFile string

var errFileMissing = errors.New("File does not exist")

func ErrFileMissing(file string) error {
	return fmt.Errorf("%w: %s", errFileMissing, file)
}

func (yamlFile *YamlFile) Validate() error {
	yamlFilePath := string(*yamlFile)
	if _, err := os.Stat(yamlFilePath); errors.Is(err, os.ErrNotExist) {
		return ErrFileMissing(yamlFilePath)
	}
	return nil
}

func main() {
	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_, err := maxprocs.Set()
	if err != nil {
		panic(err)
	}

	log.InitLogger()

	// parse arguments
	var cli cli

	// just need the yaml file, from which all the dependencies will deduced
	kong.Parse(&cli,
		kong.Name("bash-compiler"),
		kong.Description("From a yaml file describing the bash application, interprets the templates and import the necessary bash functions"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	compiler.GenerateCode(string(cli.YamlFile))
}
