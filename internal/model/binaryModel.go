// Package model allowing to load different data models
package model

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
	"github.com/goccy/go-yaml"
)

type CompilerConfig struct {
	AnnotationsConfig               structures.Dictionary `yaml:"annotationsConfig"`
	TargetFile                      string                `yaml:"targetFile"`
	RelativeRootDirBasedOnTargetDir string                `yaml:"relativeRootDirBasedOnTargetDir"`
	CommandDefinitionFiles          []string              `yaml:"commandDefinitionFiles"`
	TemplateFile                    string                `yaml:"templateFile"`
	TemplateDirs                    []string              `yaml:"templateDirs"`
	FunctionsIgnoreRegexpList       []string              `yaml:"functionsIgnoreRegexpList"`
	SrcDirs                         []string              `yaml:"srcDirs"`
	SrcDirsExpanded                 []string              `yaml:"-"`
	TargetDir                       string                `yaml:"-"`
	KeepIntermediateFiles           bool                  `yaml:"-"`
	BinaryModelFilePath             string                `yaml:"-"`
	BinaryModelBaseName             string                `yaml:"-"`
	IntermediateFilesCount          int                   `yaml:"-"`
}

func (compilerConfig *CompilerConfig) DebugCopyGeneratedFile(
	code string,
	suffix string,
) {
	if compilerConfig.KeepIntermediateFiles {
		compilerConfig.IntermediateFilesCount++
		logger.DebugCopyGeneratedFile(
			compilerConfig.TargetDir,
			compilerConfig.BinaryModelBaseName,
			fmt.Sprintf("-4-compiler-%d%s.sh", compilerConfig.IntermediateFilesCount, suffix),
			code,
		)
	}
}

type BinaryModel struct {
	CompilerConfig CompilerConfig        `yaml:"compilerConfig"`
	Vars           structures.Dictionary `yaml:"vars"`
	BinData        any                   `yaml:"binData"`
}

type BinaryModelLoader struct{}

func NewBinaryModelLoader() *BinaryModelLoader {
	return &BinaryModelLoader{}
}

func (binaryModelContext *BinaryModelLoader) Load(
	targetDir string,
	binaryModelFilePath string,
	binaryModelBaseName string,
	referenceDir string,
	keepIntermediateFiles bool,
) (_ *BinaryModel, err error) {
	modelMap := map[string]any{}
	loadedFiles := map[string]string{}
	err = loadModel(
		referenceDir,
		binaryModelFilePath,
		&modelMap,
		&loadedFiles,
		"",
	)
	if err != nil {
		return nil, err
	}

	// create temp file
	tempYamlFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempYamlFile.Name())
	err = writeYamlFile(modelMap, *tempYamlFile)
	if err != nil {
		return nil, err
	}
	if keepIntermediateFiles {
		err = logger.DebugSaveGeneratedFile(
			targetDir,
			binaryModelBaseName,
			"-1-merged.yaml",
			tempYamlFile.Name(),
		)
		if err != nil {
			return nil, err
		}
	}

	var resultWriter bytes.Buffer
	err = transformModel(*tempYamlFile, &resultWriter)
	if err != nil {
		return nil, err
	}
	if keepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			targetDir,
			binaryModelBaseName,
			"-2-cue-transformed.yaml",
			resultWriter.String(),
		)
		if err != nil {
			return nil, err
		}
	}

	// load command yaml data model
	slog.Info("Loading binaryModel", logger.LogFieldFilePath, binaryModelFilePath)
	binaryModel := BinaryModel{} //nolint:exhaustruct // load from yaml
	err = yaml.Unmarshal(resultWriter.Bytes(), &binaryModel)
	if err != nil {
		return nil, err
	}

	binaryModelContext.setEnvVars(&binaryModel)
	binaryModelContext.expandVars(&binaryModel)

	return &binaryModel, err
}

func (*BinaryModelLoader) setEnvVars(binaryModel *BinaryModel) {
	for key, value := range binaryModel.Vars {
		if val, ok := value.(string); ok {
			os.Setenv(key, val)
		}
	}
}

func (*BinaryModelLoader) expandVars(binaryModel *BinaryModel) {
	binaryModel.CompilerConfig.SrcDirsExpanded = []string{}
	for _, srcDir := range binaryModel.CompilerConfig.SrcDirs {
		binaryModel.CompilerConfig.SrcDirsExpanded = append(
			binaryModel.CompilerConfig.SrcDirsExpanded,
			os.ExpandEnv(srcDir),
		)
	}
}
