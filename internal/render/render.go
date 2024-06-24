package render

import (
	"bytes"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/files"
)

type Context struct {
	Template     *template.Template
	TemplateName string
	RootData     interface{}
	Data         interface{}
}

const (
	UserReadWritePerm        os.FileMode = 0600
	UserReadWriteExecutePerm os.FileMode = 0700
)

func NewTemplate(
	templateDirs []string,
	templateFile string,
	funcMap template.FuncMap,
) (templateInstance *template.Template, templateName string, err error) {
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
		return nil, "", err
	}

	templateBaseFile := path.Base(templateFile)
	templateName = strings.TrimSuffix(templateBaseFile, filepath.Ext(templateBaseFile))
	slog.Info("Loaded template", "Name", templateName, "AvailableFile", files)

	myTemplate := template.New(templateName).Option("missingkey=zero").Funcs(funcMap)
	_, err = myTemplate.ParseFiles(files...)
	if err != nil {
		return nil, "", err
	}

	return myTemplate, templateName, nil
}

func (templateContext Context) RenderFromTemplateName() (code string, err error) {
	code, err = templateContext.Render(templateContext.TemplateName)
	if err != nil {
		return "", err
	}

	return code, err
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
