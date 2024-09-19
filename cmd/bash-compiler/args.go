package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

const constMaxScreenSize = 80

type getCurrentFilenameError struct {
	error
}

func (*getCurrentFilenameError) Error() string {
	return "unable to get the current filename"
}

type cli struct {
	YamlFiles             YamlFiles            `arg:""                                             help:"Yaml files"                                                 optional:"" type:"path"`
	RootDirectory         Directory            `help:"Root directory containing binary files"      short:"r"                                                         type:"path"`
	TargetDir             Directory            `help:"Directory that will contain generated files" optional:""                                                       short:"t"`
	BinaryFilesExtension  BinaryFilesExtension `default:"-binary.yaml"                             help:"Provide the extension for automatic search of binary files" optional:""`
	Version               VersionFlag          `help:"Print version information and quit"          name:"version"                                                    short:"v"`
	KeepIntermediateFiles bool                 `help:"Keep intermediate files in target directory" short:"k"`
	Debug                 bool                 `help:"Set log in debug level"                      short:"d"`
	LogLevel              int                  `hidden:""`
	CompilerRootDir       Directory            `hidden:""`
}

type (
	VersionFlag          string
	Directory            string
	ConfigFile           string
	BinaryFilesExtension string
	YamlFiles            []string
)

func (yamlFiles *YamlFiles) Validate() error {
	for _, yamlFile := range *yamlFiles {
		err := files.FileExists(yamlFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (VersionFlag) IsBool() bool                       { return true }
func (VersionFlag) BeforeApply(
	app *kong.Kong, vars kong.Vars,
) error { //nolint:unparam // need to conform to interface
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
		kong.Description("From a yaml file describing the bash application, "+
			"interprets the templates and import the necessary bash functions"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary:        false,
			Summary:             true,
			Compact:             true,
			Tree:                false,
			FlagsLast:           true,
			Indenter:            kong.LineIndenter,
			NoExpandSubcommands: true,
			WrapUpperBound:      constMaxScreenSize,
		}),
		kong.Vars{
			"version": "0.1.0",
		},
	)

	// current dir
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return &getCurrentFilenameError{nil}
	}
	compilerRootDir, err := filepath.Abs(filepath.Dir(filename) + "/../..")
	if err != nil {
		return err
	}
	slog.Debug(
		"parseArgs",
		logger.LogFieldVariableName, "compilerRootDir",
		logger.LogFieldVariableValue, compilerRootDir,
	)

	cli.CompilerRootDir = Directory(compilerRootDir)
	if cli.TargetDir == "" {
		cli.TargetDir = Directory(compilerRootDir)
	}
	if cli.RootDirectory == "" {
		currentDir, err := os.Getwd()
		logger.Check(err)
		cli.RootDirectory = Directory(currentDir)
	}
	if cli.Debug {
		cli.LogLevel = int(slog.LevelDebug)
	}
	return nil
}
