package template

import (
	"bytes"
	"log/slog"
	"path"
	"path/filepath"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/files"
)

type Context struct {
	Template *template.Template
	Name     string
	RootData any
	Data     any
}

func NewTemplate(templateDir string, templateFile string,
	funcMap template.FuncMap) (templateContext *Context, err error) {
	files, err := files.MatchPatterns(
		filepath.Join(templateDir, "**/*.gtpl"),
		filepath.Join(templateDir, "**.gtpl"),
	)
	if err != nil {
		return nil, err
	}
	name := path.Base(templateFile)
	slog.Info("Loaded template", name, files)

	myTemplate := template.New(name).Funcs(funcMap)
	_, err = myTemplate.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	templateContext = &Context{myTemplate, name, nil, nil}

	return templateContext, nil
}

func (templateContext Context) Render(template string) (string, error) {
	var tplWriter bytes.Buffer
	slog.Debug("Render template", slog.String("template", template))
	err := templateContext.Template.ExecuteTemplate(&tplWriter, template, templateContext)
	return tplWriter.String(), err
}
