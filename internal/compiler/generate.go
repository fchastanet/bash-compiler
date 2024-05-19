// Package compiler
package compiler

import (
	"log/slog"
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
func GenerateCode(binaryModelFilePath string) {
	// load command yaml data model
	slog.Info("Loading", "binaryModelFilePath", binaryModelFilePath)
	binaryModel, err := model.LoadBinaryModel(binaryModelFilePath)
	if err != nil {
		panic(err)
	}

	// load template system
	templateContext, err := myTemplate.NewTemplate(
		binaryModel.BinFile.TemplateDir,
		binaryModel.BinFile.TemplateFile,
		myTemplateFunctions.FuncMap(),
	)
	if err != nil {
		panic(err)
	}

	// render
	var str string
	templateContext.Data = &binaryModel.BinData
	templateContext.RootData = templateContext.Data
	str, err = templateContext.Render("commands")
	if err != nil {
		panic(err)
	}

	// Save resulting file
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.sh", []byte(str), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.sh")
}
