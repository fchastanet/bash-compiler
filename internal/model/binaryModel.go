// Package model allowing to load different data models
package model

import (
	"os"

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
type BinaryModel struct {
	BinFile BinFileModel `yaml:"binFile"`
	Vars    interface{}  `yaml:"vars"`
	BinData interface{}  `yaml:"binData"`
}

// LoadBinaryModel loads yaml file containing binary related data
func LoadBinaryModel(filePath string) (binaryModel BinaryModel, err error) {
	yamlFileContent, err := os.ReadFile(filePath)
	binaryModel = BinaryModel{}
	if err != nil {
		return binaryModel, err
	}
	err = yaml.Unmarshal(yamlFileContent, &binaryModel)
	if err != nil {
		return binaryModel, err
	}

	// basic structure checks (json schema)
	return binaryModel, nil
}

func (binaryModel BinaryModel) InitTemplateContext() (templateContext *render.Context, err error) {
	// load template system
	myTemplate, templateName, err := render.NewTemplate(
		binaryModel.BinFile.TemplateDirs,
		binaryModel.BinFile.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["binData"] = binaryModel.BinData
	data["binFile"] = binaryModel.BinFile
	data["vars"] = binaryModel.Vars

	return &render.Context{
		Template:     myTemplate,
		TemplateName: templateName,
		RootData:     data,
		Data:         data,
	}, nil
}
