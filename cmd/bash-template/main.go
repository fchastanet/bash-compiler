// main package
package main

import (
	"fmt"
	"log/slog"
	"os"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
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
	var str string
	templateContext, err := myTemplate.NewTemplate(templateDir, templateFile, myTemplateFunctions.FuncMap())
	if err != nil {
		panic(err)
	}

	templateData := TemplateData{
		Name: "Example",
	}

	str, err = templateContext.Render(templateContext.Name, templateData)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
