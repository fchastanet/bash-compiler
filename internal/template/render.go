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
	RootData *any
	Data     *any
}

func NewTemplate(templateDirs []string, templateFile string,
	funcMap template.FuncMap) (templateContext *Context, err error) {
	var patterns = [3]string{
		"**/**/*.*",
		"**/*.*",
		"*.*",
	}
	templateDirPatterns := make([]string, len(templateDirs)*len(patterns))
	for _, templateDir := range templateDirs {
		for _, pattern := range patterns {
			templateDirPatterns = append(templateDirPatterns, filepath.Join(templateDir, pattern))
		}
	}
	files, err := files.MatchPatterns(templateDirPatterns...)
	if err != nil {
		return nil, err
	}
	name := path.Base(templateFile)
	slog.Info("Loaded template", "Name", name, "AvailableFile", files)

	myTemplate := template.New(name).Option("missingkey=zero").Funcs(funcMap)
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
	if err != nil {
		return "", err
	}
	return tplWriter.String(), err
}
