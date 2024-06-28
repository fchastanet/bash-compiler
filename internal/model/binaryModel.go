// Package model allowing to load different data models
package model

import (
	"bytes"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/render"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render/functions"
	"github.com/goccy/go-yaml"
)

type BinFileModel struct {
	TargetFile                      string   `yaml:"targetFile"`
	RelativeRootDirBasedOnTargetDir string   `yaml:"relativeRootDirBasedOnTargetDir"`
	CommandDefinitionFiles          []string `yaml:"commandDefinitionFiles"`
	TemplateFile                    string   `yaml:"templateFile"`
	TemplateDirs                    []string `yaml:"templateDirs"`
	SrcDirs                         []string `yaml:"srcDirs"`
}
type CompilerConfig struct {
	FunctionsIgnoreRegexpList []string `yaml:"functionsIgnoreRegexpList"`
	SrcDirs                   []string `yaml:"srcDirs"`
	SrcDirsExpanded           []string `yaml:"-"`
}

type BinaryModel struct {
	BinFile        BinFileModel   `yaml:"binFile"`
	Vars           interface{}    `yaml:"vars"`
	BinData        interface{}    `yaml:"binData"`
	CompilerConfig CompilerConfig `yaml:"compilerConfig"`
}

type BinaryModelContext struct {
	BinaryModel           *BinaryModel
	TargetDir             string
	BinaryModelFilePath   string
	BinaryModelBaseName   string
	ReferenceDir          string
	KeepIntermediateFiles bool
}

type BinaryModelInterface interface {
	LoadBinaryModel() (err error)
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

func (binaryModelContext *BinaryModelContext) LoadBinaryModel() (err error) {
	modelMap := map[string]interface{}{}
	err = LoadModel(
		binaryModelContext.ReferenceDir,
		binaryModelContext.BinaryModelFilePath,
		&modelMap,
	)
	if err != nil {
		return err
	}

	// create temp file
	tempYamlFile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempYamlFile.Name())
	err = WriteYamlFile(modelMap, *tempYamlFile)
	if err != nil {
		return err
	}
	if binaryModelContext.KeepIntermediateFiles {
		err = logger.DebugSaveGeneratedFile(
			binaryModelContext.TargetDir,
			binaryModelContext.BinaryModelBaseName,
			"-merged.yaml",
			tempYamlFile.Name(),
		)
		if err != nil {
			return err
		}
	}

	var resultWriter bytes.Buffer
	err = transformModel(*tempYamlFile, &resultWriter)
	if err != nil {
		return err
	}
	if binaryModelContext.KeepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			binaryModelContext.TargetDir,
			binaryModelContext.BinaryModelBaseName,
			"-cue-transformed.yaml",
			resultWriter.String(),
		)
		if err != nil {
			return err
		}
	}

	// load command yaml data model
	slog.Info("Loading", "binaryModelFilePath", binaryModelContext.BinaryModelFilePath)
	binaryModel := BinaryModel{}
	err = yaml.Unmarshal(resultWriter.Bytes(), &binaryModel)
	if err != nil {
		return err
	}
	binaryModelContext.BinaryModel = &binaryModel

	binaryModelContext.expandVars()

	return err
}

func (binaryModelContext *BinaryModelContext) expandVars() {
	binaryModelContext.BinaryModel.CompilerConfig.SrcDirsExpanded = []string{}
	for _, srcDir := range binaryModelContext.BinaryModel.CompilerConfig.SrcDirs {
		binaryModelContext.BinaryModel.CompilerConfig.SrcDirsExpanded = append(
			binaryModelContext.BinaryModel.CompilerConfig.SrcDirsExpanded,
			os.ExpandEnv(srcDir),
		)
	}
}

func NewTemplateContext(binaryModelContext BinaryModelContext) (templateContext *render.Context, err error) {
	// load template system
	myTemplate, templateName, err := render.NewTemplate(
		binaryModelContext.BinaryModel.BinFile.TemplateDirs,
		binaryModelContext.BinaryModel.BinFile.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["binData"] = binaryModelContext.BinaryModel.BinData
	data["binFile"] = binaryModelContext.BinaryModel.BinFile
	data["vars"] = binaryModelContext.BinaryModel.Vars

	return &render.Context{
		Template:     myTemplate,
		TemplateName: templateName,
		RootData:     data,
		Data:         data,
	}, nil
}
