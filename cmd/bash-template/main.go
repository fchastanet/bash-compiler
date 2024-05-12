// main package
package main

import (
	"fmt"
	"log/slog"
	"os"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	"github.com/fchastanet/bash-compiler/internal/template/functions"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
	"gopkg.in/yaml.v3"
)

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

type TemplateData struct {
	Name string
}

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
	yamlData := functions.FromYAMLFile("templates-examples/shellcheckLint.yaml")
	templateContext.Data = &yamlData
	templateContext.RootData = templateContext.Data
	yamlDataTransformed, err := templateContext.Render("command.dataModel")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("logs/main.transformed.yaml", []byte(yamlDataTransformed), UserReadWritePerm); err != nil {
		panic(err)
	}
	slog.Info("Check logs/main.transformed.yaml")

	var yamlData2 interface{}
	err = yaml.Unmarshal([]byte(yamlDataTransformed), &yamlData2)
	if err != nil {
		panic(err)
	}
	out, err := yaml.Marshal(yamlData2)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("logs/main.marshalled.yaml", out, UserReadWritePerm); err != nil {
		panic(err)
	}
	slog.Info("Check logs/main.marshalled.yaml")

	// render
	var str string
	templateContext.Data = &yamlData2
	templateContext.RootData = templateContext.Data
	str, err = templateContext.Render(templateContext.Name)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("logs/main.sh", []byte(str), UserReadWriteExecutePerm); err != nil {
		panic(err)
	}
	slog.Info("Check logs/main.marshalled.yaml")
}
