package compiler

import (
	"testing"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"gotest.tools/v3/assert"
)

func newCompiler(
	templateContextRenderFromTemplateContent TemplateContextRenderFromTemplateContent,
) CompileContext {
	templateContext := templateContextMock{
		nil,
		templateContextRenderFromTemplateContent,
	}
	compileContext := NewCompiler(
		&templateContext,
		[]AnnotationProcessorInterface{},
	)
	templateName := "fakeTemplate"
	compileContext.Init(
		&render.TemplateContextData{
			TemplateContext: &templateContext,
			TemplateName:    &templateName,
			Template:        nil,
			RootData:        nil,
			Data:            nil,
		},
		&model.CompilerConfig{}, //nolint:exhaustruct // test
	)
	return compileContext
}

func TestCompileEmptyCode(t *testing.T) {
	compilerContext := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	resultCode, err := compilerContext.Compile(
		&CompileContextData{}, //nolint:exhaustruct // test
		"",
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionNotFound(t *testing.T) {
	compilerContext := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	resultCode, err := compilerContext.Compile(
		&CompileContextData{ //nolint:exhaustruct // test
			config: &model.CompilerConfig{ //nolint:exhaustruct // test
				FunctionsIgnoreRegexpList: []string{},
			},
			functionsMap: make(map[string]functionInfoStruct),
		},
		"MyPackage::function",
	)
	assert.Error(t, err, "function not found: MyPackage::function in any srcDirs []")
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionIgnoredFunction(t *testing.T) {
	compilerContext := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	resultCode, err := compilerContext.Compile(
		&CompileContextData{ //nolint:exhaustruct // test
			config: &model.CompilerConfig{ //nolint:exhaustruct // test
				FunctionsIgnoreRegexpList: []string{
					"Ignore::ignoredFunction",
				},
			},
			functionsMap: make(map[string]functionInfoStruct),
		},
		"Ignore::ignoredFunction",
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Ignore::ignoredFunction\n", resultCode)
}
