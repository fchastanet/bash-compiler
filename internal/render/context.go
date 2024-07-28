package render

import (
	"bytes"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/utils/bash"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

type TemplateContext struct {
}

type TemplateContextData struct {
	TemplateContext *TemplateContext
	TemplateDirs    []string
	TemplateFile    string
	TemplateName    string
	Template        *template.Template
	RootData        interface{}
	Data            interface{}
}

type TemplateContextInterface interface {
	Init(funcMap map[string]interface{}) error
	Render(templateName string) (string, error)
	RenderFromTemplateName() (code string, err error)
	RenderFromTemplateContent(templateContent string) (codeStr string, err error)
}

func NewTemplateContext() (templateContext *TemplateContext) {
	return &TemplateContext{}
}

func (templateContext *TemplateContext) Init(
	templateDirs []string,
	templateFile string,
	data interface{},
	funcMap map[string]interface{},
) (*TemplateContextData, error) {
	templateContextData := &TemplateContextData{
		TemplateContext: templateContext,
		TemplateDirs:    templateDirs,
		TemplateFile:    templateFile,
		TemplateName:    "",
		RootData:        data,
		Data:            data,
	}
	// load template system
	myTemplate, templateName, err := newTemplate(
		templateContextData.TemplateDirs,
		templateContextData.TemplateFile,
		funcMap,
	)
	if err != nil {
		return nil, err
	}
	templateContextData.Template = myTemplate
	templateContextData.TemplateName = templateName
	return templateContextData, nil
}

func (templateContext *TemplateContext) Render(
	templateContextData *TemplateContextData,
	templateName string,
) (string, error) {
	var tplWriter bytes.Buffer
	slog.Debug("Render template", logger.LogFieldTemplateName, templateName)
	err := templateContextData.Template.ExecuteTemplate(&tplWriter, templateName, templateContextData)
	if err != nil {
		return "", err
	}
	return tplWriter.String(), err
}

func (templateContext *TemplateContext) RenderFromTemplateName(
	templateContextData *TemplateContextData,
) (code string, err error) {
	code, err = templateContext.Render(templateContextData, templateContextData.TemplateName)
	if err != nil {
		return "", err
	}

	return code, err
}

func (templateContext *TemplateContext) RenderFromTemplateContent(
	templateContextData *TemplateContextData,
	templateContent string,
) (codeStr string, err error) {
	template, err := templateContextData.Template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	var tplWriter bytes.Buffer
	err = template.Execute(&tplWriter, templateContextData)
	if err != nil {
		return "", err
	}

	return bash.RemoveFirstShebangLineIfAny(tplWriter.String()), err
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
