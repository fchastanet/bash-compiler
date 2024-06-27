// main package
package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/model"
	"go.uber.org/automaxprocs/maxprocs"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
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

func loadBinaryModel(cli cli, binaryModelFilePath string, binaryModelBaseName string) (
	binaryModel *model.BinaryModel, err error) {
	referenceDir := filepath.Dir(binaryModelFilePath)
	modelMap := map[string]interface{}{}
	err = model.LoadModel(referenceDir, binaryModelFilePath, &modelMap)
	if err != nil {
		return nil, err
	}

	// create temp file
	tempYamlFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempYamlFile.Name())
	err = model.WriteYamlFile(modelMap, *tempYamlFile)
	if err != nil {
		return nil, err
	}
	err = debugSaveGeneratedFile(cli, binaryModelBaseName, "-merged.yaml", tempYamlFile.Name())
	if err != nil {
		return nil, err
	}

	var resultWriter bytes.Buffer
	err = model.TransformModel(*tempYamlFile, &resultWriter)
	if err != nil {
		return nil, err
	}
	err = debugCopyGeneratedFile(cli, binaryModelBaseName, "-cue-transformed.yaml", resultWriter.String())
	if err != nil {
		return nil, err
	}

	// load command yaml data model
	slog.Info("Loading", "binaryModelFilePath", binaryModelFilePath)
	var binaryModelVar model.BinaryModel
	binaryModelVar, err = model.LoadBinaryModel(resultWriter.Bytes())
	if err != nil {
		return nil, err
	}
	return &binaryModelVar, nil
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

	// Load binary model
	binaryModelFilePath := string(cli.YamlFile)
	binaryModelBaseName := files.BaseNameWithoutExtension(binaryModelFilePath)
	binaryModel, err := loadBinaryModel(cli, binaryModelFilePath, binaryModelBaseName)
	logger.Check(err)

	// Render code using template
	templateContext, err := binaryModel.InitTemplateContext()
	logger.Check(err)
	code, err := templateContext.RenderFromTemplateName()
	logger.Check(err)
	err = debugCopyGeneratedFile(cli, binaryModelBaseName, "-afterTemplateRendering.sh", code)
	logger.Check(err)

	// Compile
	codeCompiled, err := compiler.Compile(code, templateContext, *binaryModel)
	logger.Check(err)

	// Save resulting file
	targetFile := os.ExpandEnv(binaryModel.BinFile.TargetFile)
	err = os.WriteFile(targetFile, []byte(codeCompiled), UserReadWriteExecutePerm)
	logger.Check(err)
	slog.Info("Compiled", "file", targetFile)
}

func debugSaveGeneratedFile(
	cli cli, basename string, suffix string, tempYamlFile string,
) (err error) {
	if cli.KeepIntermediateFiles {
		targetFile := filepath.Join(
			string(cli.TargetDir),
			fmt.Sprintf("%s%s", basename, suffix),
		)
		err := files.Copy(tempYamlFile, targetFile)
		if err != nil {
			return err
		}
		slog.Info("KeepIntermediateFiles", "merged config file", targetFile)
	}
	return nil
}

func debugCopyGeneratedFile(
	cli cli, basename string, suffix string, code string,
) (err error) {
	if cli.KeepIntermediateFiles {
		targetFile := filepath.Join(
			string(cli.TargetDir),
			fmt.Sprintf("%s%s", basename, suffix),
		)
		err := os.WriteFile(targetFile, []byte(code), UserReadWriteExecutePerm)
		slog.Info("KeepIntermediateFiles", "merged config file", targetFile)
		return err
	}
	return nil
}
