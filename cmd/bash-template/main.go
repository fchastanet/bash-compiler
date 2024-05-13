// main package
package main

import (
	"fmt"
	"log/slog"
	"os"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	"github.com/fchastanet/bash-compiler/internal/template/functions"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
	"github.com/goccy/go-yaml"
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

	// load an transform yaml data
	filePath := "templates-examples/shellcheckLint.yaml"
	yamlData := functions.FromYAMLFile(filePath)
	templateContext.Data = &yamlData
	templateContext.RootData = templateContext.Data
	yamlDataTransformed, err := templateContext.Render("dataModel")
	slog.Info("Data Model transformed", "filePath", filePath)
	if err != nil {
		panic(err)
	}
	var yamlData1 interface{}
	err = yaml.Unmarshal([]byte(yamlDataTransformed), &yamlData1)
	if err != nil {
		panic(err)
	}
	// we have to transform twice to take the value included into account
	templateContext.Data = &yamlData1
	templateContext.RootData = templateContext.Data
	yamlDataTransformed, err = templateContext.Render("dataModel")
	slog.Info("Data Model transformed pass 2", "filePath", filePath)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.dataModel.yaml", []byte(yamlDataTransformed), UserReadWritePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.dataModel.yaml")

	var yamlData2 interface{}
	err = yaml.Unmarshal([]byte(yamlDataTransformed), &yamlData2)
	if err != nil {
		panic(err)
	}
	out, err := yaml.MarshalWithOptions(yamlData2, yaml.Indent(2), yaml.IndentSequence(true))
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("templates-examples/testsData/shellcheckLint.dataModelMarshalled.yaml", out, UserReadWritePerm); err != nil {
		panic(err)
	}
	slog.Info("Check templates-examples/testsData/shellcheckLint.dataModelMarshalled.yaml")

	// render
	var str string
	templateContext.Data = &yamlData2
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
