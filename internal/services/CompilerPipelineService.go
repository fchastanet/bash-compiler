package services

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/dotenv"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

type CompilerPipelineService struct {
	rootDirectory         string
	compilerRootDir       string
	targetDir             string
	yamlFiles             []string
	binaryFilesExtension  string
	debug                 bool
	keepIntermediateFiles bool

	binaryModelService *BinaryModelServiceContext
}

func NewCompilerPipelineService(
	rootDirectory string,
	compilerRootDir string,
	targetDir string,
	yamlFiles []string,
	binaryFilesExtension string,
	debug bool,
	keepIntermediateFiles bool,
) (_ *CompilerPipelineService) {
	return &CompilerPipelineService{
		rootDirectory:         rootDirectory,
		compilerRootDir:       compilerRootDir,
		targetDir:             targetDir,
		yamlFiles:             yamlFiles,
		binaryFilesExtension:  binaryFilesExtension,
		debug:                 debug,
		keepIntermediateFiles: keepIntermediateFiles,
		binaryModelService:    nil,
	}
}

func (service *CompilerPipelineService) Init() error {
	// set useful env variables that can be interpolated during template rendering
	err := setEnvVariable("COMPILER_ROOT_DIR", service.compilerRootDir)
	if err != nil {
		return err
	}
	err = setEnvVariable("ROOT_DIR", service.rootDirectory)
	if err != nil {
		return err
	}

	// load config file
	err = service.loadConfFile()
	if err != nil {
		return err
	}
	if service.debug {
		envVars := os.Environ()
		for _, envVar := range envVars {
			slog.Debug("env", "var", envVar)
		}
	}

	service.initBinaryModelService()
	return nil
}

func (service *CompilerPipelineService) initBinaryModelService() {
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
	var templateContextInterface TemplateContextInterface = templateContext
	var compilerInterface CodeCompilerInterface = compilerService
	intermediateFileCallback := skipIntermediateFilesCallback
	intermediateFileContentCallback := skipIntermediateFilesCallback
	if service.keepIntermediateFiles {
		intermediateFileCallback = logger.DebugCopyIntermediateFile
		intermediateFileContentCallback = logger.DebugSaveIntermediateFile
	}
	service.binaryModelService = NewBinaryModelService(
		model.NewBinaryModelLoader(),
		templateContextInterface,
		compilerInterface,
		intermediateFileCallback,
		intermediateFileContentCallback,
	)
}

func (service *CompilerPipelineService) ProcessPipeline() error {
	defaultLogger := slog.Default()

	err := service.computeYamlFiles()
	if err != nil {
		return err
	}
	for _, binaryModelFilePath := range service.yamlFiles {
		slog.SetDefault(defaultLogger.With("binaryModelFilePath", binaryModelFilePath))
		binaryModelServiceContextData, err := service.binaryModelService.Init(
			service.targetDir,
			binaryModelFilePath,
		)
		if err != nil {
			return err
		}
		err = service.binaryModelService.Compile(binaryModelServiceContextData)
		if err != nil {
			return err
		}
	}
	return nil
}

// load .bash-compiler file in current directory if exists
func (service *CompilerPipelineService) loadConfFile() error {
	configFile := filepath.Join(service.rootDirectory, ".bash-compiler")
	err := files.FileExists(configFile)
	if err != nil {
		slog.Warn("Config file is not available or not readable", "configFile", configFile)
		return nil //nolint:nilerr // error ignored
	}
	slog.Info("Loading", logger.LogFieldFilePath, configFile)
	return dotenv.LoadEnvFile(configFile)
}

func (service *CompilerPipelineService) computeYamlFiles() (err error) {
	if len(service.yamlFiles) == 0 {
		filesList, err := files.MatchPatterns(
			service.rootDirectory,
			"**/*"+service.binaryFilesExtension,
		)
		if err != nil {
			return err
		}
		if len(filesList) == 0 {
			slog.Error(
				"cannot find any file with specified suffix and directory",
				"rootDirectory", service.rootDirectory,
				"extension", service.binaryFilesExtension,
			)
			return err
		}
		service.yamlFiles = filesList
	}
	return nil
}

func skipIntermediateFilesCallback(
	_ string, _ string, _ string, _ string,
) error {
	return nil
}

func setEnvVariable(name string, value string) error {
	slog.Debug(
		"main",
		logger.LogFieldVariableName, name,
		logger.LogFieldVariableValue, value,
	)
	return os.Setenv(name, value)
}