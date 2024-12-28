// main package
package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/services"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"go.uber.org/automaxprocs/maxprocs"
)

//go:embed defaultTemplates/*
var templateFs embed.FS

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

	// create temporary directory
	templateTempDir, err := os.MkdirTemp("", "bash-compiler")
	logger.Check(err)
	err = files.CopyEmbeddedFiles(templateFs, "defaultTemplates", templateTempDir)
	logger.Check(err)
	err = os.Setenv("DEFAULT_TEMPLATE_FOLDER", templateTempDir)
	logger.Check(err)
	slog.Info("Default template folder", "folder", templateTempDir)

	compilerPipelineService := services.NewCompilerPipelineService(
		string(cli.RootDirectory),
		[]string(cli.YamlFiles),
		string(cli.BinaryFilesExtension),
		cli.Debug,
		string(cli.IntermediateFilesDir),
	)
	err = compilerPipelineService.Init()
	logger.Check(err)
	err = compilerPipelineService.ProcessPipeline()
	logger.Check(err)
}
