package main

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

var errToGetCurrentFilename = errors.New("Unable to get the current filename")

type cli struct {
	YamlFiles             YamlFiles   `arg:"" help:"Yaml files" type:"path"`
	TargetDir             Directory   `short:"t" optional:"" help:"Directory that will contain generated files"`
	Version               VersionFlag `short:"v" name:"version" help:"Print version information and quit"`
	KeepIntermediateFiles bool        `short:"k" help:"Keep intermediate files in target directory"`
	Debug                 bool        `short:"d" help:"Set log in debug level"`
	LogLevel              int         `hidden:""`
	CompilerRootDir       Directory   `hidden:""`
}

type VersionFlag string
type Directory string
type YamlFiles []string

func (yamlFiles *YamlFiles) Validate() error {
	for _, yamlFile := range *yamlFiles {
		err := files.FileExists(yamlFile)
		if err != nil {
			return err
		}
	}
	return nil
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

	// current dir
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errToGetCurrentFilename
	}
	compilerRootDir, err := filepath.Abs(filepath.Dir(filename) + "/../..")
	if err != nil {
		return err
	}
	slog.Info(
		"parseArgs",
		logger.LogFieldVariableName, "compilerRootDir",
		logger.LogFieldVariableValue, compilerRootDir,
	)

	cli.CompilerRootDir = Directory(compilerRootDir)
	if cli.TargetDir == "" {
		cli.TargetDir = Directory(compilerRootDir)
	}
	if cli.Debug {
		cli.LogLevel = int(slog.LevelDebug)
	}
	return nil
}
