// Package model allowing to load different data models
package model

import (
	"bytes"
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
}

type BinaryModel struct {
	CompilerConfig CompilerConfig        `yaml:"compilerConfig"`
	Vars           structures.Dictionary `yaml:"vars"`
	BinData        interface{}           `yaml:"binData"`
}

type BinaryModelContext struct{}

func NewBinaryModel() *BinaryModelContext {
	return &BinaryModelContext{}
}

func (binaryModelContext *BinaryModelContext) Load(
	targetDir string,
	binaryModelFilePath string,
	binaryModelBaseName string,
	referenceDir string,
	keepIntermediateFiles bool,
) (_ *BinaryModel, err error) {
	modelMap := map[string]interface{}{}
	err = loadModel(
		referenceDir,
		binaryModelFilePath,
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
