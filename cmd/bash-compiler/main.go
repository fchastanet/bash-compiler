// main package
package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/services"
	"github.com/fchastanet/bash-compiler/internal/utils/dotenv"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

func skipIntermediateFilesCallback(
	_ string, _ string, _ string, _ string,
) error {
	return nil
}

// load .bash-compiler file in current directory if exists
func loadConfFile(cli *cli) {
	configFile := filepath.Join(string(cli.RootDirectory), ".bash-compiler")
	err := files.FileExists(configFile)
	if err == nil {
		slog.Info("Loading", logger.LogFieldFilePath, configFile)
		err = dotenv.LoadEnvFile(configFile)
		logger.Check(err)
	} else {
		slog.Warn("Config file is not available or not readable", "configFile", configFile)
	}
}

func getYamlFiles(cli *cli) (yamlFiles YamlFiles, err error) {
	yamlFiles = cli.YamlFiles
	if len(cli.YamlFiles) == 0 {
		filesList, err := files.MatchPatterns(
			string(cli.RootDirectory),
			"**/*"+string(cli.BinaryFilesExtension),
		)
		logger.Check(err)
		if len(filesList) == 0 {
			slog.Error(
				"cannot find any file with specified suffix and directory",
				"rootDirectory", string(cli.RootDirectory),
				"extension", cli.BinaryFilesExtension,
			)
			return cli.YamlFiles, err
		}
		yamlFiles = filesList
	}
	return yamlFiles, nil
}

func setEnvVariable(name string, value string) {
	slog.Debug(
		"main",
		logger.LogFieldVariableName, name,
		logger.LogFieldVariableValue, value,
	)
	err := os.Setenv(name, value)
	logger.Check(err)
}

func main() {
	// This controls the maxprocs environment variable in container runtimes.
	// see https://tinyurl.com/3rfwknuv
	_, err := maxprocs.Set()
	logger.Check(err)

	// get current dir
	currentDir, err := os.Getwd()
	logger.Check(err)
	err = os.Setenv("PWD", currentDir)
	logger.Check(err)

	// parse arguments
	var cli cli
	err = parseArgs(&cli)
	logger.Check(err)

	// set useful env variables that can be interpolated during template rendering
	setEnvVariable("COMPILER_ROOT_DIR", string(cli.CompilerRootDir))
	setEnvVariable("ROOT_DIR", string(cli.RootDirectory))

	// load config file
	loadConfFile(&cli)
	logger.InitLogger(cli.LogLevel)
	if cli.Debug {
		envVars := os.Environ()
		for _, envVar := range envVars {
			slog.Debug("env", "var", envVar)
		}
	}

	// create BinaryModelService
	templateContext := render.NewTemplateContext()
	requireAnnotationProcessor := compiler.NewRequireAnnotationProcessor()
	embedAnnotationProcessor := compiler.NewEmbedAnnotationProcessor()
	compilerService := compiler.NewCompiler(
		templateContext,
		[]compiler.AnnotationProcessorInterface{
			requireAnnotationProcessor,
			embedAnnotationProcessor,
		},
	)
	var templateContextInterface services.TemplateContextInterface = templateContext
	var compilerInterface services.CodeCompilerInterface = compilerService
	intermediateFileCallback := skipIntermediateFilesCallback
	intermediateFileContentCallback := skipIntermediateFilesCallback
	if cli.KeepIntermediateFiles {
		intermediateFileCallback = logger.DebugCopyIntermediateFile
		intermediateFileContentCallback = logger.DebugSaveIntermediateFile
	}
	binaryModelService := services.NewBinaryModelService(
		model.NewBinaryModelLoader(),
		templateContextInterface,
		compilerInterface,
		intermediateFileCallback,
		intermediateFileContentCallback,
	)
	defaultLogger := slog.Default()

	yamlFiles, err := getYamlFiles(&cli)
	logger.Check(err)
	for _, binaryModelFilePath := range yamlFiles {
		slog.SetDefault(defaultLogger.With("binaryModelFilePath", binaryModelFilePath))
		binaryModelServiceContextData, err := binaryModelService.Init(
			string(cli.TargetDir),
			binaryModelFilePath,
		)
		logger.Check(err)
		err = binaryModelService.Compile(binaryModelServiceContextData)
		logger.Check(err)
	}
}
