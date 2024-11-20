package compiler

import (
	"regexp"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
	"gotest.tools/v3/assert"
)

func TestEmbedInitCompileContextDataNil(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(nil)
	assert.Error(t, err, "validation failed invalid value : "+
		"context annotationEmbed field compileContextData value <nil>")
}

func TestEmbedInitEmptyCompileContextData(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(&CompileContextData{}) //nolint:exhaustruct // test
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.compileContext value <nil>")
}

func TestEmbedInitInvalidCompileContextDataMissingTemplateContextData(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(&CompileContextData{&CompileContext{}, nil, nil, nil, nil}) //nolint:exhaustruct // test
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.templateContextData value <nil>")
}

func TestEmbedInitInvalidCompileContextDataMissingConfig(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		nil,
		nil,
		nil,
	})
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.config value <nil>")
}

func TestEmbedInitInvalidCompileContextDataMissingFunctionMap(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{},       //nolint:exhaustruct // test
		nil,
		nil,
	})
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.functionsMap value map[]")
}

func TestEmbedInitInvalidCompileContextDataMissingEmbedFileTemplateName(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	// jscpd:ignore-start
	err := embedProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{},       //nolint:exhaustruct // test
		make(map[string]functionInfoStruct),
		[]*regexp.Regexp{},
	})
	// jscpd:ignore-end
	assert.Error(t, err, "validation failed invalid value : "+
		"context compileContextData.config.AnnotationsConfig field embedFileTemplateName value <nil> inner error missing key: embedFileTemplateName")
}

func TestEmbedInitInvalidCompileContextDataMissingEmbedFileTemplateDir(t *testing.T) {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{ //nolint:exhaustruct // test
			AnnotationsConfig: structures.Dictionary{
				"embedFileTemplateName": "templateName",
			},
		},
		make(map[string]functionInfoStruct),
		[]*regexp.Regexp{},
	})
	assert.Error(t, err, "validation failed invalid value : "+
		"context compileContextData.config.AnnotationsConfig field embedDirTemplateName value <nil> inner error missing key: embedDirTemplateName")
}

func getValidEmbedProcessor(t *testing.T) AnnotationProcessorInterface {
	embedProcessor := NewEmbedAnnotationProcessor()
	err := embedProcessor.Init(getCompileContextData())
	assert.Equal(t, nil, err)
	return embedProcessor
}

func TestEmbedParseFunction(t *testing.T) {
	embedProcessor := getValidEmbedProcessor(t)
	err := embedProcessor.ParseFunction(&CompileContextData{}, &functionInfoStruct{}) //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
}

func TestEmbedProcess(t *testing.T) {
	embedProcessor := getValidEmbedProcessor(t)
	err := embedProcessor.Process(&CompileContextData{}) //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
}

func TestEmbedPostProcessEmptyString(t *testing.T) {
	embedProcessor := getValidEmbedProcessor(t)
	code, err := embedProcessor.PostProcess(&CompileContextData{}, "") //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
	assert.Equal(t, "", code)
}

func TestEmbedPostProcessNoMatch(t *testing.T) {
	embedProcessor := getValidEmbedProcessor(t)
	code, err := embedProcessor.PostProcess(&CompileContextData{}, "noMatchCode") //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
	assert.Equal(t, "noMatchCode\n", code)
}

type GenerateCodeFunc func() (string, error)

type annotationEmbedGenerateMock struct {
	GenerateCodeFunc
}

func (annotation *annotationEmbedGenerateMock) RenderResource(
	_ string, _ string, _ int,
) (string, error) {
	return annotation.GenerateCodeFunc()
}

func (*annotationEmbedGenerateMock) GetTitle() string {
	return ""
}

func (*annotationEmbedGenerateMock) Reset() {
}

func getEmbedProcessorMocked(generateCodeFunc GenerateCodeFunc) *embedAnnotationProcessor {
	embedProcessor := &embedAnnotationProcessor{
		annotationProcessor:     annotationProcessor{},
		annotationEmbedGenerate: &annotationEmbedGenerateMock{generateCodeFunc},
		embedMap:                make(map[string]string),
	}

	return embedProcessor
}

func getCompileContextData() *CompileContextData {
	return &CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{ //nolint:exhaustruct // test
			AnnotationsConfig: structures.Dictionary{
				"embedFileTemplateName": "templateName",
				"embedDirTemplateName":  "templateDir",
			},
		},
		make(map[string]functionInfoStruct),
		[]*regexp.Regexp{},
	}
}

func TestEmbedPostProcessOneMatch(t *testing.T) {
	mock := getEmbedProcessorMocked(
		func() (string, error) {
			return "mock", nil
		},
	)
	var embedProcessor AnnotationProcessorInterface = mock
	compileContextData := getCompileContextData()
	code, err := embedProcessor.PostProcess(compileContextData, "# @embed srcFile AS targetFile")
	assert.Equal(t, nil, err)
	assert.Equal(t, "mock\n", code)
}

func TestEmbedPostProcessTwoMatches(t *testing.T) {
	mock := getEmbedProcessorMocked(
		func() (string, error) {
			return "mock", nil
		},
	)
	var embedProcessor AnnotationProcessorInterface = mock
	compileContextData := getCompileContextData()
	code, err := embedProcessor.PostProcess(
		compileContextData,
		"# @embed srcFile AS targetFile\n# @embed srcFile AS targetFile",
	)
	assert.Error(t, err, "Embedded resource 'srcFile' - name 'targetFile' is duplicated on line 2")
	assert.Equal(t, "", code)
}

func TestEmbedPostProcessGenerateError(t *testing.T) {
	mock := getEmbedProcessorMocked(func() (string, error) {
		return "", &unsupportedEmbeddedResourceError{nil, "asName", "resource", 12}
	})
	var embedProcessor AnnotationProcessorInterface = mock
	compileContextData := getCompileContextData()
	code, err := embedProcessor.PostProcess(
		compileContextData,
		"# @embed srcFile AS targetFile\n# @embed srcFile AS targetFile",
	)
	assert.Error(t, err, "Embedded resource 'resource' - name 'asName' on line 12 cannot be embedded")
	assert.Equal(t, "", code)
}
