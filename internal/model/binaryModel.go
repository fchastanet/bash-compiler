// Package model allowing to load different data models
package model

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/utils"
	"github.com/goccy/go-yaml"
)

type CompilerConfig struct {
	AnnotationsConfig               utils.Dictionary `yaml:"annotationsConfig"`
	TargetFile                      string           `yaml:"targetFile"`
	RelativeRootDirBasedOnTargetDir string           `yaml:"relativeRootDirBasedOnTargetDir"`
	CommandDefinitionFiles          []string         `yaml:"commandDefinitionFiles"`
	TemplateFile                    string           `yaml:"templateFile"`
	TemplateDirs                    []string         `yaml:"templateDirs"`
	FunctionsIgnoreRegexpList       []string         `yaml:"functionsIgnoreRegexpList"`
	SrcDirs                         []string         `yaml:"srcDirs"`
	SrcDirsExpanded                 []string         `yaml:"-"`
}

type BinaryModel struct {
	CompilerConfig CompilerConfig   `yaml:"compilerConfig"`
	Vars           utils.Dictionary `yaml:"vars"`
	BinData        interface{}      `yaml:"binData"`
}

type BinaryModelContext struct {
	TargetDir             string
	BinaryModelFilePath   string
	BinaryModelBaseName   string
	ReferenceDir          string
	KeepIntermediateFiles bool
}

type BinaryModelInterface interface {
	Load() (binaryModel *BinaryModel, err error)
}

func NewBinaryModel(
	targetDir string,
	binaryModelFilePath string,
	binaryModelBaseName string,
	referenceDir string,
	keepIntermediateFiles bool,
) *BinaryModelContext {
	return &BinaryModelContext{
		TargetDir:             targetDir,
		BinaryModelFilePath:   binaryModelFilePath,
		BinaryModelBaseName:   binaryModelBaseName,
		ReferenceDir:          referenceDir,
		KeepIntermediateFiles: keepIntermediateFiles,
	}
}

func (binaryModelContext *BinaryModelContext) Load() (_ *BinaryModel, err error) {
	modelMap := map[string]interface{}{}
	err = loadModel(
		binaryModelContext.ReferenceDir,
		binaryModelContext.BinaryModelFilePath,
		&modelMap,
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
	if binaryModelContext.KeepIntermediateFiles {
		err = logger.DebugSaveGeneratedFile(
			binaryModelContext.TargetDir,
			binaryModelContext.BinaryModelBaseName,
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
	if binaryModelContext.KeepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			binaryModelContext.TargetDir,
			binaryModelContext.BinaryModelBaseName,
			"-2-cue-transformed.yaml",
			resultWriter.String(),
		)
		if err != nil {
			return nil, err
		}
	}

	// load command yaml data model
	slog.Info("Loading binaryModel", logger.LogFieldFilePath, binaryModelContext.BinaryModelFilePath)
	binaryModel := BinaryModel{}
	err = yaml.Unmarshal(resultWriter.Bytes(), &binaryModel)
	if err != nil {
		return nil, err
	}

	binaryModelContext.setEnvVars(&binaryModel)
	binaryModelContext.expandVars(&binaryModel)

	return &binaryModel, err
}

func (binaryModelContext *BinaryModelContext) setEnvVars(binaryModel *BinaryModel) {
	for key, value := range binaryModel.Vars {
		if val, ok := value.(string); ok {
			os.Setenv(key, val)
		}
	}
}

func (binaryModelContext *BinaryModelContext) expandVars(binaryModel *BinaryModel) {
	binaryModel.CompilerConfig.SrcDirsExpanded = []string{}
	for _, srcDir := range binaryModel.CompilerConfig.SrcDirs {
		binaryModel.CompilerConfig.SrcDirsExpanded = append(
			binaryModel.CompilerConfig.SrcDirsExpanded,
			os.ExpandEnv(srcDir),
		)
	}
}
