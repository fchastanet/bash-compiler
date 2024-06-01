// Package compiler
package compiler

import (
	"os"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/template"
	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

func InitTemplateContext(binaryModel model.BinaryModel) (templateContext *template.Context, err error) {
	// load template system
	templateContext, err = myTemplate.NewTemplate(
		binaryModel.BinFile.TemplateDirs,
		binaryModel.BinFile.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		return nil, err
	}

	// render
	data := make(map[string]interface{})
	data["binData"] = binaryModel.BinData
	data["binFile"] = binaryModel.BinFile
	data["vars"] = binaryModel.Vars

	templateContext.Data = data
	templateContext.RootData = templateContext.Data
	return templateContext, nil
}

func RenderFromTemplateName(templateContext *template.Context, templateName string) (code string, err error) {
	code, err = templateContext.Render(templateName)
	if err != nil {
		return "", err
	}

	return code, err
}
