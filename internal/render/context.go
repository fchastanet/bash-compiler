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
	TemplateName    *string
	Template        templateInterface
	RootData        any
	Data            any
}

func NewTemplateContext() (templateContext *TemplateContext) {
	return &TemplateContext{}
}

func (templateContext *TemplateContext) Init(
	templateDirs []string,
	templateFile string,
	data any,
	funcMap map[string]any,
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
		TemplateName:    &templateName,
		Template:        myTemplate,
		RootData:        data,
		Data:            data,
	}

	return templateContextData, nil
}

func (*TemplateContext) Render(
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
	code, err = templateContext.Render(templateContextData, *templateContextData.TemplateName)
	if err != nil {
		return "", err
	}

	return code, err
}

func (*TemplateContext) RenderFromTemplateContent(
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
	var filesList []string
	for _, templateDir := range templateDirs {
		myFiles, err := files.MatchPatterns(templateDir, "*.gtpl", "**/*.gtpl")
		if err != nil {
			return nil, "", err
		}
		filesList = append(filesList, myFiles...)
	}

	templateBaseFile := path.Base(templateFile)
	templateName = strings.TrimSuffix(templateBaseFile, filepath.Ext(templateBaseFile))
	slog.Debug(
		"Loaded template",
		logger.LogFieldTemplateName, templateName,
		logger.LogFieldAvailableTemplateFiles, filesList,
	)

	myTemplate := template.New(templateName).Option("missingkey=zero").Funcs(funcMap)
	_, err = myTemplate.ParseFiles(filesList...)
	if err != nil {
		return nil, "", err
	}

	return myTemplate, templateName, nil
}
