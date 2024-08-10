package render

import (
	"bytes"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fchastanet/bash-compiler/internal/utils/bash"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

type TemplateContext struct{}

type templateInterface interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
	Parse(text string) (*template.Template, error)
}

type TemplateContextInterface interface {
	Render(
		templateContextData *TemplateContextData,
		templateName string,
	) (string, error)
	RenderFromTemplateContent(
		templateContextData *TemplateContextData,
		templateContent string,
	) (codeStr string, err error)
}

type TemplateContextData struct {
	TemplateContext TemplateContextInterface
	TemplateName    string
	Template        templateInterface
	RootData        interface{}
	Data            interface{}
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
	// load template system
	myTemplate, templateName, err := newTemplate(
		templateDirs,
		templateFile,
		funcMap,
	)
	if err != nil {
		return nil, err
	}

	templateContextData := &TemplateContextData{
		TemplateContext: templateContext,
		TemplateName:    templateName,
		Template:        myTemplate,
		RootData:        data,
		Data:            data,
	}

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
	myTemplate, err := templateContextData.Template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	var tplWriter bytes.Buffer
	err = myTemplate.Execute(&tplWriter, templateContextData)
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
	patterns := [3]string{
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
	myFiles, err := files.MatchPatterns(templateDirPatterns...)
	if err != nil {
		return nil, "", err
	}

	templateBaseFile := path.Base(templateFile)
	templateName = strings.TrimSuffix(templateBaseFile, filepath.Ext(templateBaseFile))
	slog.Info(
		"Loaded template",
		logger.LogFieldTemplateName, templateName,
		logger.LogFieldAvailableTemplateFiles, myFiles,
	)

	myTemplate := template.New(templateName).Option("missingkey=zero").Funcs(funcMap)
	_, err = myTemplate.ParseFiles(myFiles...)
	if err != nil {
		return nil, "", err
	}

	return myTemplate, templateName, nil
}
