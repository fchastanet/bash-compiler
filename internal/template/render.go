package template

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/files"
)

var TemplateDir string

func Render(
	templateDir string,
	templateFile string,
	templateData any,
	funcMap template.FuncMap) (string, error) {
	TemplateDir = templateDir // TODO find another to store this

	files, err := files.MatchPatterns(
		filepath.Join(templateDir, "**/*.tpl"),
		filepath.Join(templateDir, "**.tpl"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)

	name := path.Base(templateFile)
	myTemplate := template.New(name).
		Funcs(funcMap)
	_, err = myTemplate.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	// TODO render should only do this part, myTemplate should be stored somewhere
	var tplWriter bytes.Buffer
	err = myTemplate.ExecuteTemplate(&tplWriter, name, templateData)
	return tplWriter.String(), err
}
