// Package compiler
package compiler

import (
	"os"

	"github.com/fchastanet/bash-compiler/internal/model"
	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

// GenerateCode generates code from given model
func GenerateCode(binaryModel *model.BinaryModel) (code string, err error) {
	// load template system
	templateContext, err := myTemplate.NewTemplate(
		binaryModel.BinFile.TemplateDirs,
		binaryModel.BinFile.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		return "", err
	}

	// render
	templateContext.Data = &binaryModel.BinData
	templateContext.RootData = templateContext.Data
	code, err = templateContext.Render("commands")
	if err != nil {
		return "", err
	}

	return code, err
}
