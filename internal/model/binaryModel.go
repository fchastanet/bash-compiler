// Package model allowing to load different data models
package model

import (
	"os"

	"github.com/goccy/go-yaml"
)

type BinFileModel struct {
	TargetFile   string   `yaml:"targetFile"`
	TemplateFile string   `yaml:"templateFile"`
	TemplateName string   `yaml:"templateName"`
	TemplateDirs []string `yaml:"templateDirs"`
	SrcDirs      []string `yaml:"srcDirs"`
}
type BinaryModel struct {
	BinFile BinFileModel `yaml:"binFile"`
	BinData interface{}  `yaml:"binData"`
}

// LoadBinaryModel loads yaml file containing binary related data
func LoadBinaryModel(filePath string) (binaryModel *BinaryModel, err error) {
	yamlFileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	binaryModel = &BinaryModel{}
	err = yaml.Unmarshal(yamlFileContent, binaryModel)
	if err != nil {
		return nil, err
	}

	// basic structure checks (json schema)
	return binaryModel, nil
}
