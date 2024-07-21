// Package model allowing to load different data models
package model

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/render"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render/functions"
	"github.com/goccy/go-yaml"
)

var errMissingKey = errors.New("Invalid key")

func ErrMissingKey(key string) error {
	return fmt.Errorf("%w: %s", errMissingKey, key)
}

var errInvalidType = errors.New("Invalid type")

func ErrInvalidType(myVar interface{}) error {
	return fmt.Errorf("%w: %T", errInvalidType, myVar)
}

type Dictionary map[string]interface{}

func (dic Dictionary) GetStringValue(key string) (value string, err error) {
	val, ok := dic[key]
	if !ok {
		return "", ErrMissingKey(key)
	}
	if value, ok := val.(string); ok {
		return ExpandStringValue(value), nil
	}

	return "", ErrInvalidType(value)
}

func ExpandStringValue(value string) string {
	return os.ExpandEnv(value)
}

func (dic Dictionary) GetStringList(key string) (values []string, err error) {
	val, ok := dic[key]
	if !ok {
		return nil, ErrMissingKey(key)
	}
	if values, ok := val.([]string); ok {
		return ExpandStringList(values), nil
	}

	return nil, ErrInvalidType(values)
}

func ExpandStringList(values []string) []string {
	slice := make([]string, len(values))
	for i := len(values) - 1; i >= 0; i-- {
		slice[i] = os.ExpandEnv(values[i])
	}
	return slice
}

type CompilerConfig struct {
	AnnotationsConfig               Dictionary `yaml:"annotationsConfig"`
	TargetFile                      string     `yaml:"targetFile"`
	RelativeRootDirBasedOnTargetDir string     `yaml:"relativeRootDirBasedOnTargetDir"`
	CommandDefinitionFiles          []string   `yaml:"commandDefinitionFiles"`
	TemplateFile                    string     `yaml:"templateFile"`
	TemplateDirs                    []string   `yaml:"templateDirs"`
	FunctionsIgnoreRegexpList       []string   `yaml:"functionsIgnoreRegexpList"`
	SrcDirs                         []string   `yaml:"srcDirs"`
	SrcDirsExpanded                 []string   `yaml:"-"`
}

type BinaryModel struct {
	CompilerConfig CompilerConfig `yaml:"compilerConfig"`
	Vars           Dictionary     `yaml:"vars"`
	BinData        interface{}    `yaml:"binData"`
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
			"-1-merged.yaml",
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
			"-2-cue-transformed.yaml",
			resultWriter.String(),
		)
		if err != nil {
			return err
		}
	}

	// load command yaml data model
	slog.Info("Loading binaryModel", logger.LogFieldFilePath, binaryModelContext.BinaryModelFilePath)
	binaryModel := BinaryModel{}
	err = yaml.Unmarshal(resultWriter.Bytes(), &binaryModel)
	if err != nil {
		return err
	}
	binaryModelContext.BinaryModel = &binaryModel

	binaryModelContext.setEnvVars()
	binaryModelContext.expandVars()

	return err
}

func (binaryModelContext *BinaryModelContext) setEnvVars() {
	for key, value := range binaryModelContext.BinaryModel.Vars {
		if val, ok := value.(string); ok {
			os.Setenv(key, val)
		}
	}
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
	templateDirs := ExpandStringList(binaryModelContext.BinaryModel.CompilerConfig.TemplateDirs)
	// load template system
	myTemplate, templateName, err := render.NewTemplate(
		templateDirs,
		binaryModelContext.BinaryModel.CompilerConfig.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["binData"] = binaryModelContext.BinaryModel.BinData
	data["compilerConfig"] = binaryModelContext.BinaryModel.CompilerConfig
	data["vars"] = binaryModelContext.BinaryModel.Vars

	return &render.Context{
		Template:     myTemplate,
		TemplateName: templateName,
		RootData:     data,
		Data:         data,
	}, nil
}
