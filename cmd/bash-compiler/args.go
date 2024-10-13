package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
)

const constMaxScreenSize = 80

type getCurrentFilenameError struct {
	error
}

func (*getCurrentFilenameError) Error() string {
	return "unable to get the current filename"
}

type rootDirError struct {
	error
}

func (*rootDirError) Error() string {
	return "please provide rootDir option"
}

type missingBashCompilerFileError struct {
	error
}

func (*missingBashCompilerFileError) Error() string {
	return "current directory should contain file .bash-compiler"
}

type rootDirOptionShouldNotBeProvidedError struct {
	error
}

func (*rootDirOptionShouldNotBeProvidedError) Error() string {
	return "rootDir option should not be provided"
}

type cli struct {
	YamlFiles            YamlFiles            `arg:""    optional:"" type:"path"                help:"Yaml files"`                                                            //nolint:tagalign //avoid reformat annotations
	RootDirectory        RootDirectory        `short:"r" optional:"" type:"path" name:"rootDir" help:"Root directory containing binary files"`                                //nolint:tagalign //avoid reformat annotations
	IntermediateFilesDir IntermediateFilesDir `short:"t" optional:""                            help:"Directory that will contain generated files (no save if not provided)"` //nolint:tagalign //avoid reformat annotations
	BinaryFilesExtension BinaryFilesExtension `          optional:"" default:"-binary.yaml"     help:"Provide the extension for automatic search of binary files"`            //nolint:tagalign //avoid reformat annotations
	Version              VersionFlag          `short:"v" name:"version"                         help:"Print version information and quit"`                                    //nolint:tagalign //avoid reformat annotations
	Debug                bool                 `short:"d"                                        help:"Set log in debug level"`                                                //nolint:tagalign //avoid reformat annotations
	LogLevel             int                  `hidden:""`
}

type (
	VersionFlag          string
	IntermediateFilesDir string
	RootDirectory        string
	ConfigFile           string
	BinaryFilesExtension string
	YamlFiles            []string
)

func isUsingGoRun() bool {
	executable := filepath.Base(os.Args[0])
	return strings.HasPrefix(executable, "__debug_bin") ||
		strings.HasPrefix(os.Args[0], "/tmp")
}

func (
	o cli, //nolint:gocritic // hugeparam: no need to optimize, called one time
) BeforeReset(ctx *kong.Context) error { //nolint:unparam // need to conform to interface
	if isUsingGoRun() {
		return nil
	}
	for _, flag := range ctx.Kong.Model.Flags {
		if flag.Name == "rootDir" {
			// hide rootDir as needed only for debug or when using go run
			flag.Hidden = true
		}
	}
	return nil
}

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

func (intermediateFilesDir *IntermediateFilesDir) Validate(
	_ *kong.Kong, _ kong.Vars,
) error {
	if (*intermediateFilesDir) == "" {
		return nil
	}
	return files.IsWritableDirectory(string(*intermediateFilesDir))
}

func parseArgs(cli *cli) (err error) {
	// just need the yaml file, from which all the dependencies will be deduced
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
			"version": "3.0.0",
		},
	)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if isUsingGoRun() {
		if cli.RootDirectory == "" {
			return &rootDirError{nil}
		}
	} else {
		if cli.RootDirectory != "" {
			return &rootDirOptionShouldNotBeProvidedError{nil}
		}
		cli.RootDirectory = RootDirectory(currentDir)
	}
	bashCompilerFile := filepath.Join(string(cli.RootDirectory), ".bash-compiler")
	if _, err = os.Stat(bashCompilerFile); err != nil {
		slog.Error("current directory should contain file .bash-compiler", "expectedFile", bashCompilerFile)
		return &missingBashCompilerFileError{err}
	}

	if cli.Debug {
		cli.LogLevel = int(slog.LevelDebug)
	}
	return nil
}
