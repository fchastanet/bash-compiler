// main package
package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fchastanet/bash-compiler/internal/services"
	"github.com/fchastanet/bash-compiler/internal/utils/dotenv"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_, err := maxprocs.Set()
	logger.Check(err)

	// get current dir
	currentDir, err := os.Getwd()
	logger.Check(err)
	os.Setenv("PWD", currentDir)

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
	slog.Info(
		"main",
		logger.LogFieldVariableName, "COMPILER_ROOT_DIR",
		logger.LogFieldVariableValue, string(cli.CompilerRootDir),
	)
	os.Setenv("COMPILER_ROOT_DIR", string(cli.CompilerRootDir))

	// create BinaryModelService
	for _, binaryModelFilePath := range cli.YamlFiles {
		binaryModelService := services.NewBinaryModelService(
			string(cli.TargetDir),
			cli.KeepIntermediateFiles,
			binaryModelFilePath,
		)
		err = binaryModelService.Init()
		logger.Check(err)
		err = binaryModelService.Compile()
		logger.Check(err)
	}
}
