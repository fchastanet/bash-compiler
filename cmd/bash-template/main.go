// main package
package main

import (
	"fmt"
	"os"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/template/functions"
)

type TemplateData struct {
	Name string
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	argsWithoutProg := os.Args[1:]
	templateDir := fmt.Sprintf("%s/%s", currentDir, argsWithoutProg[0])
	var templateFile = fmt.Sprintf("%s/%s", templateDir, argsWithoutProg[1])
	templateData := TemplateData{
		Name: "Example",
	}
	var str string
	str, err = myTemplate.Render(templateDir, templateFile, templateData, myTemplateFunctions.FuncMap())
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
