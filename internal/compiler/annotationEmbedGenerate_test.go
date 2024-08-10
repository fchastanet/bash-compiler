package compiler

import (
	"os"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/render"
	"gotest.tools/v3/assert"
)

const shouldNotBeCalledCodeStr = "should not be called"

var shouldNotBeCalledError = unsupportedEmbeddedResourceError{nil, "asName", "resource", 12}

type TemplateContextRenderFunc func(
	templateContextData *render.TemplateContextData,
	templateName string,
) (string, error)

type TemplateContextRenderFromTemplateContent func(
	templateContextData *render.TemplateContextData,
	templateContent string,
) (codeStr string, err error)

type templateContextMock struct {
	templateContextRenderFunc                TemplateContextRenderFunc
	templateContextRenderFromTemplateContent TemplateContextRenderFromTemplateContent
}

func (templateContext *templateContextMock) Render(
	templateContextData *render.TemplateContextData,
	code string,
) (string, error) {
	return templateContext.templateContextRenderFunc(templateContextData, code)
}

func (templateContext *templateContextMock) RenderFromTemplateContent(
	templateContextData *render.TemplateContextData,
	code string,
) (codeStr string, err error) {
	return templateContext.templateContextRenderFromTemplateContent(
		templateContextData,
		code,
	)
}

func newAnnotationEmbedGenerate(
	templateContextRenderFunc TemplateContextRenderFunc,
	templateContextRenderFromTemplateContent TemplateContextRenderFromTemplateContent,
) *annotationEmbedGenerate {
	return &annotationEmbedGenerate{
		embedDirTemplateName:  "embedDirTemplateName",
		embedFileTemplateName: "embedFileTemplateName",
		templateContextData: &render.TemplateContextData{
			TemplateContext: &templateContextMock{
				templateContextRenderFunc,
				templateContextRenderFromTemplateContent,
			},
			TemplateName: "fakeTemplate",
			Template:     nil,
			RootData:     nil,
			Data:         nil,
		},
	}
}

func TestRenderResourceNotFound(t *testing.T) {
	embedGenerate := newAnnotationEmbedGenerate(
		func(_ *render.TemplateContextData, _ string) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
		func(_ *render.TemplateContextData, _ string) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	code, err := embedGenerate.RenderResource("asName", "resource", 1)
	assert.Error(t, err, "Embedded resource 'resource' - name 'asName' on line 1 cannot be embedded - inner error: stat resource: no such file or directory")
	assert.Equal(t, "", code)
}

func TestRenderResourceFile(t *testing.T) {
	embedGenerate := newAnnotationEmbedGenerate(
		func(
			templateContextData *render.TemplateContextData,
			templateName string,
		) (string, error) {
			rootDataMap, ok := templateContextData.RootData.(map[string]string)
			assert.Equal(t, true, ok)
			assert.Equal(t, "asName", rootDataMap["asName"])
			assert.Equal(t, "embedFileTemplateName", templateName)
			return "transformed code", nil
		},
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	code, err := embedGenerate.RenderResource("asName", "annotationEmbed.go", 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, "transformed code", code)
}

func TestRenderResourceDir(t *testing.T) {
	embedGenerate := newAnnotationEmbedGenerate(
		func(
			templateContextData *render.TemplateContextData,
			templateName string,
		) (string, error) {
			rootDataMap, ok := templateContextData.RootData.(map[string]string)
			assert.Equal(t, true, ok)
			assert.Equal(t, "asName", rootDataMap["asName"])
			assert.Equal(t, "embedDirTemplateName", templateName)
			return "transformed code", nil
		},
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	pwd, err := os.Getwd()
	assert.Equal(t, err, nil)
	code, err := embedGenerate.RenderResource("asName", pwd, 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, "transformed code", code)
}
