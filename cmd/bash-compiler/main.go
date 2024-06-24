// main package
package main

import (
	"fmt"
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
	YamlFile             YamlFile    `arg:"" help:"Yaml file" type:"path"`
	IntermediateFilesDir Directory   `short:"i" help:"save intermediate files to directory"`
	Version              VersionFlag `name:"version" help:"Print version information and quit"`
}

type VersionFlag string
type Directory string
type YamlFile string

func (yamlFile *YamlFile) Validate() error {
	yamlFilePath := string(*yamlFile)
	return utils.FileExists(yamlFilePath)
}

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                       { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error { //nolint: unparam
	fmt.Printf("Bash compiler version %s\n", vars["version"])
	app.Exit(0)
	return nil
}

func (directory *Directory) Validate() error {
	directoryPath := string(*directory)
	return utils.DirExists(directoryPath)
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
		}),
		kong.Vars{
			"version": "0.1.0",
		},
	)

	// load command yaml data model
	binaryModelFilePath := string(cli.YamlFile)
	slog.Info("Loading", "binaryModelFilePath", binaryModelFilePath)
	binaryModel, err := model.LoadBinaryModel(binaryModelFilePath)
	if err != nil {
		panic(err)
	}

	templateContext, err := binaryModel.InitTemplateContext()
	if err != nil {
		panic(err)
	}
	code, err := templateContext.RenderFromTemplateName()
	if err != nil {
		panic(err)
	}

	// Save resulting file
	if err := os.WriteFile("examples/binaries/shellcheckLint.beforeCompile.sh", []byte(code), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check examples/binaries/shellcheckLint.beforeCompile.sh")

	// Compile
	codeCompiled, err := compiler.Compile(code, templateContext, binaryModel)
	if err != nil {
		panic(err)
	}

	// Save resulting file
	targetFile := os.ExpandEnv(binaryModel.BinFile.TargetFile)

	if err := os.WriteFile(targetFile, []byte(codeCompiled), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check", "file", targetFile)
}
