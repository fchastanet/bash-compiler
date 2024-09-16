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

	// load .bash-compiler file in current directory if exists
	bashCompilerConfFile := filepath.Join(currentDir, ".bash-compiler")
	err = files.FileExists(bashCompilerConfFile)
	if err == nil {
		slog.Info("Loading", logger.LogFieldFilePath, bashCompilerConfFile)
		err = dotenv.LoadEnvFile(bashCompilerConfFile)
		logger.Check(err)
	} else {
		slog.Warn(".bash-compiler file not available")
	}

	// parse arguments
	var cli cli
	err = parseArgs(&cli)
	logger.Check(err)
	logger.InitLogger(cli.LogLevel)

	// set useful env variables that can be interpolated during template rendering
	slog.Debug(
		"main",
		logger.LogFieldVariableName, "COMPILER_ROOT_DIR",
		logger.LogFieldVariableValue, string(cli.CompilerRootDir),
	)
	err = os.Setenv("COMPILER_ROOT_DIR", string(cli.CompilerRootDir))
	logger.Check(err)

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
	binaryModelService := services.NewBinaryModelService(
		model.NewBinaryModelLoader(),
		templateContextInterface,
		compilerInterface,
	)
	defaultLogger := slog.Default()
	for _, binaryModelFilePath := range cli.YamlFiles {
		slog.SetDefault(defaultLogger.With(
			"binaryModelFilePath", binaryModelFilePath,
		))
		binaryModelServiceContextData, err := binaryModelService.Init(
			string(cli.TargetDir),
			cli.KeepIntermediateFiles,
			binaryModelFilePath,
		)
		logger.Check(err)
		err = binaryModelService.Compile(binaryModelServiceContextData)
		logger.Check(err)
	}
}
