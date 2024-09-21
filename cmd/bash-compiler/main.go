// main package
package main

import (
	"os"

	"github.com/fchastanet/bash-compiler/internal/services"
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

	// parse arguments
	var cli cli
	err = parseArgs(&cli)
	logger.Check(err)
	logger.InitLogger(cli.LogLevel)

	compilerPipelineService := services.NewCompilerPipelineService(
		string(cli.RootDirectory),
		string(cli.CompilerRootDir),
		string(cli.TargetDir),
		[]string(cli.YamlFiles),
		string(cli.BinaryFilesExtension),
		cli.Debug,
		cli.KeepIntermediateFiles,
	)
	compilerPipelineService.Init()
	err = compilerPipelineService.ProcessPipeline()
	logger.Check(err)
}
