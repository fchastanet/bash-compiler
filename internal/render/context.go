package render

import (
	"bytes"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/code"
	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/logger"
)

type Context struct {
	TemplateDirs []string
	Template     *template.Template
	TemplateFile string
	TemplateName string
	RootData     interface{}
	Data         interface{}
}

func (templateContext *Context) Render(templateName string) (string, error) {
	var tplWriter bytes.Buffer
	slog.Debug("Render template", logger.LogFieldTemplateName, templateName)
	err := templateContext.Template.ExecuteTemplate(&tplWriter, templateName, templateContext)
	if err != nil {
		return "", err
	}
	return tplWriter.String(), err
}

func (templateContext *Context) RenderFromTemplateName() (code string, err error) {
	code, err = templateContext.Render(templateContext.TemplateName)
	if err != nil {
		return "", err
	}

	return code, err
}

func (templateContext *Context) RenderFromTemplateContent(templateContent string) (codeStr string, err error) {
	template, err := templateContext.Template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	var tplWriter bytes.Buffer
	err = template.Execute(&tplWriter, templateContext)
	if err != nil {
		return "", err
	}

	return code.RemoveFirstShebangLineIfAny(tplWriter.String()), err
}

func NewTemplateContext(
	templateDirs []string,
	templateFile string,
	data interface{},
) (templateContext *Context) {
	return &Context{
		TemplateDirs: templateDirs,
		TemplateFile: templateFile,
		Template:     nil,
		TemplateName: "",
		RootData:     data,
		Data:         data,
	}
}

func (templateContext *Context) Init(funcMap map[string]interface{}) error {
	// load template system
	myTemplate, templateName, err := newTemplate(
		templateContext.TemplateDirs,
		templateContext.TemplateFile,
		funcMap,
	)
	if err != nil {
		return err
	}
	templateContext.Template = myTemplate
	templateContext.TemplateName = templateName
	return nil
}

func newTemplate(
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
	slog.Info(
		"Loaded template",
		logger.LogFieldTemplateName, templateName,
		logger.LogFieldAvailableTemplateFiles, files,
	)

	myTemplate := template.New(templateName).Option("missingkey=zero").Funcs(funcMap)
	_, err = myTemplate.ParseFiles(files...)
	if err != nil {
		return nil, "", err
	}

	return myTemplate, templateName, nil
}
