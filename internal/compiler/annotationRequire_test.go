package compiler

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestRequireInitCompileContextDataNil(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(nil)
	assert.Error(t, err, "validation failed invalid value : "+
		"context annotationEmbed field compileContextData value <nil>")
}

func TestRequireInitEmptyCompileContextData(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(&CompileContextData{}) //nolint:exhaustruct // test
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.compileContext value <nil>")
}

func TestRequireInitInvalidCompileContextDataMissingTemplateContextData(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(
		&CompileContextData{
			&CompileContext{},
			nil, nil, nil, nil,
		},
	) //exhaustruct:ignore
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.templateContextData value <nil>")
}

func TestRequireInitInvalidCompileContextDataMissingConfig(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(
		&CompileContextData{
			&CompileContext{},             //nolint:exhaustruct // test
			&render.TemplateContextData{}, //nolint:exhaustruct // test
			nil,
			nil,
			nil,
		},
	)
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.config value <nil>")
}

func TestRequireInitInvalidCompileContextDataMissingFunctionMap(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{},       //nolint:exhaustruct // test
		nil,
		nil,
	})
	assert.Error(t, err, "validation failed invalid value : "+
		"context compiler field CompileContextData.functionsMap value map[]")
}

func TestRequireInitInvalidCompileContextDataMissingRequirementsTemplateName(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	// jscpd:ignore-start
	err := requireProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{},       //nolint:exhaustruct // test
		make(map[string]functionInfoStruct),
		[]*regexp.Regexp{},
	})
	// jscpd:ignore-end
	assert.Error(t, err, "validation failed invalid value : "+
		"context compileContextData.config.AnnotationsConfig field checkRequirementsTemplateName value <nil> inner error missing key: checkRequirementsTemplateName")
}

func TestRequireInitInvalidCompileContextDataMissingRequireTemplate(t *testing.T) {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(&CompileContextData{
		compileContext:      &CompileContext{},             //nolint:exhaustruct // test
		templateContextData: &render.TemplateContextData{}, //nolint:exhaustruct // test
		config: &model.CompilerConfig{ //nolint:exhaustruct // test
			AnnotationsConfig: structures.Dictionary{
				"checkRequirementsTemplateName": "checkRequirementsTemplate",
			},
		},
		functionsMap:          make(map[string]functionInfoStruct),
		ignoreFunctionsRegexp: []*regexp.Regexp{},
	})
	assert.Error(t, err, "validation failed invalid value : "+
		"context compileContextData.config.AnnotationsConfig field requireTemplateName value <nil> inner error missing key: requireTemplateName")
}

func getValidRequireProcessor(t *testing.T) AnnotationProcessorInterface {
	requireProcessor := NewRequireAnnotationProcessor()
	err := requireProcessor.Init(&CompileContextData{
		&CompileContext{},             //nolint:exhaustruct // test
		&render.TemplateContextData{}, //nolint:exhaustruct // test
		&model.CompilerConfig{ //nolint:exhaustruct // test
			AnnotationsConfig: structures.Dictionary{
				"checkRequirementsTemplateName": "templateName",
				"requireTemplateName":           "template",
			},
		},
		make(map[string]functionInfoStruct),
		[]*regexp.Regexp{},
	})
	assert.Equal(t, nil, err)
	return requireProcessor
}

func TestRequireParseFunctionNoRequiredFunction(t *testing.T) {
	requireProcessor := getValidRequireProcessor(t)
	functionStruct := &functionInfoStruct{ //nolint:exhaustruct // test
		AnnotationMap: make(map[string]any),
		SourceCode:    "",
	}
	err := requireProcessor.ParseFunction(
		&CompileContextData{}, //nolint:exhaustruct // test
		functionStruct,
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", functionStruct.SourceCode)
}

func TestRequireParseFunctionWithARequiredFunction(t *testing.T) {
	requireProcessor := getValidRequireProcessor(t)
	functionStruct := &functionInfoStruct{ //nolint:exhaustruct // test
		AnnotationMap: make(map[string]any),
		FunctionName:  "Function::myFunction",
		SourceCode:    "# @require MyRequired::function\nfunction Function::myFunction() {\n:;\n}",
	}
	err := requireProcessor.ParseFunction(
		&CompileContextData{ //nolint:exhaustruct // test
			templateContextData: &render.TemplateContextData{ //nolint:exhaustruct // test
				TemplateContext: &templateContextMock{
					templateContextRenderFunc: func(
						templateContextData *render.TemplateContextData, _ string,
					) (string, error) {
						if data, ok := templateContextData.Data.(map[string]any); ok {
							return fmt.Sprintf("%v %v %v", data["code"], data["functionName"], data["requires"]), nil
						}
						return "", validationError("invalid Data", templateContextData.Data)
					},
					templateContextRenderFromTemplateContent: func(
						_ *render.TemplateContextData, _ string,
					) (string, error) {
						return "", nil
					},
				},
			},
		},
		functionStruct,
	)
	assert.Equal(t, nil, err)
	golden.Assert(t, functionStruct.SourceCode, "expectedTestRequireParseFunctionWithARequiredFunction.txt")
}

func TestRequireProcess(t *testing.T) {
	requireProcessor := getValidRequireProcessor(t)
	err := requireProcessor.Process(&CompileContextData{}) //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
}

func TestRequirePostProcessEmptyString(t *testing.T) {
	requireProcessor := getValidRequireProcessor(t)
	code, err := requireProcessor.PostProcess(&CompileContextData{}, "") //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
	assert.Equal(t, "", code)
}

func TestRequirePostProcessNoMatch(t *testing.T) {
	requireProcessor := getValidRequireProcessor(t)
	code, err := requireProcessor.PostProcess(&CompileContextData{}, "noMatchCode") //nolint:exhaustruct // test
	assert.Equal(t, nil, err)
	assert.Equal(t, "noMatchCode", code)
}

type annotationRequireGenerateMock struct {
	GenerateCodeFunc
}

func (annotation *annotationRequireGenerateMock) RenderResource(_ string, _ string, _ int) (string, error) {
	return annotation.GenerateCodeFunc()
}

func getRequireProcessorMocked() *requireAnnotationProcessor {
	requireProcessor := &requireAnnotationProcessor{ //nolint:exhaustruct // test
		annotationProcessor:           annotationProcessor{},
		checkRequirementsTemplateName: "templateName",
		requireTemplateName:           "requireTemplateName",
	}

	return requireProcessor
}

func TestRequirePostProcess(t *testing.T) {
	mock := getRequireProcessorMocked()
	var requireProcessor AnnotationProcessorInterface = mock
	compileContextData := getCompileContextData()
	code, err := requireProcessor.PostProcess(compileContextData, "# @require Env::requireLoad")
	assert.Equal(t, nil, err)
	assert.Equal(t, "# @require Env::requireLoad", code)
}
