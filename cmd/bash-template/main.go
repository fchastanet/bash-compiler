// main package
package main

import (
	"fmt"
	"log/slog"
	"os"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	"github.com/fchastanet/bash-compiler/internal/template/functions"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stderr, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	initLogger()

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	argsWithoutProg := os.Args[1:]
	templateDir := fmt.Sprintf("%s/%s", currentDir, argsWithoutProg[0])
	var templateFile = fmt.Sprintf("%s/%s", templateDir, argsWithoutProg[1])
	templateContext, err := myTemplate.NewTemplate(templateDir, templateFile, myTemplateFunctions.FuncMap())
	if err != nil {
		panic(err)
	}

	// load command yaml data model
	filePath := "templates-examples/testsData/shellcheckLint.yaml"
	yamlData := functions.FromYAMLFile(filePath)
	templateContext.Data = &yamlData
	templateContext.RootData = templateContext.Data

	// render
	var str string
	templateContext.Data = &yamlData
	templateContext.RootData = templateContext.Data
	str, err = templateContext.Render("commands")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.sh", []byte(str), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.sh")
}
