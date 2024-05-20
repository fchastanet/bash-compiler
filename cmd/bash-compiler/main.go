// main package
package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/log"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/utils"
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

func (yamlFile *YamlFile) Validate() error {
	yamlFilePath := string(*yamlFile)
	return utils.FileExists(yamlFilePath)
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

	// load command yaml data model
	binaryModelFilePath := string(cli.YamlFile)
	slog.Info("Loading", "binaryModelFilePath", binaryModelFilePath)
	binaryModel, err := model.LoadBinaryModel(binaryModelFilePath)
	if err != nil {
		panic(err)
	}

	code, err := compiler.GenerateCode(binaryModel)
	if err != nil {
		panic(err)
	}

	// Save resulting file
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.beforeCompile.sh", []byte(code), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.beforeCompile.sh")

	// Compile
	codeCompiled, err := compiler.Compile(code, binaryModel)
	if err != nil {
		panic(err)
	}

	// Save resulting file
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.sh", []byte(codeCompiled), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.sh")
}
