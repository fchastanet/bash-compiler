// main package
package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/binary"
	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

type cli struct {
	YamlFile              YamlFile    `arg:"" help:"Yaml file" type:"path"`
	TargetDir             Directory   `short:"t" optional:"" help:"Directory that will contain generated files"`
	Version               VersionFlag `short:"v" name:"version" help:"Print version information and quit"`
	KeepIntermediateFiles bool        `short:"k" help:"Keep intermediate files in target directory"`
	Debug                 bool        `short:"d" help:"Set log in debug level"`
	LogLevel              int         `hidden:""`
}

type VersionFlag string
type Directory string
type YamlFile string

func (yamlFile *YamlFile) Validate() error {
	yamlFilePath := string(*yamlFile)
	return files.FileExists(yamlFilePath)
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
	return files.DirExists(directoryPath)
}

func parseArgs(cli *cli) (err error) {
	// just need the yaml file, from which all the dependencies will deduced
	kong.Parse(cli,
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

	if cli.TargetDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		cli.TargetDir = Directory(currentDir)
	}
	if cli.Debug {
		cli.LogLevel = int(slog.LevelDebug)
	}
	return nil
}

func main() {
	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_, err := maxprocs.Set()
	logger.Check(err)

	// parse arguments
	var cli cli
	err = parseArgs(&cli)
	logger.Check(err)

	logger.InitLogger(cli.LogLevel)

	binaryModelFilePath := string(cli.YamlFile)
	binaryModelBaseName := files.BaseNameWithoutExtension(binaryModelFilePath)
	referenceDir := filepath.Dir(binaryModelFilePath)

	binaryCompiler := binary.NewCompiler()
	binaryModelContext, codeCompiled, err := binaryCompiler.Compile(
		string(cli.TargetDir),
		binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		cli.KeepIntermediateFiles,
	)
	logger.Check(err)

	// Save resulting file
	targetFile := os.ExpandEnv(binaryModelContext.BinaryModel.BinFile.TargetFile)
	err = os.WriteFile(targetFile, []byte(codeCompiled), files.UserReadWriteExecutePerm)
	logger.Check(err)
	slog.Info("Compiled", "file", targetFile)
}
