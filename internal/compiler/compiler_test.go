package compiler

import (
	"testing"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

// *******************************************
// mockAnnotation
type mockAnnotation struct {
	annotation
}

func (annotationProcessor *mockAnnotation) Init(
	_ *CompileContextData,
) error {
	return nil
}

func (annotationProcessor *mockAnnotation) ParseFunction(
	_ *CompileContextData, _ *functionInfoStruct,
) error {
	return nil
}

func (annotationProcessor *mockAnnotation) Process(
	_ *CompileContextData,
) error {
	return nil
}

func (annotationProcessor *mockAnnotation) PostProcess(
	_ *CompileContextData, code string,
) (newCode string, err error) {
	return code, nil
}

func (annotationProcessor *mockAnnotation) GetTitle() string {
	return ""
}

func (annotationProcessor *mockAnnotation) Reset() {
}

// *******************************************
// newCompiler
func newCompiler(
	templateContextRenderFromTemplateContent TemplateContextRenderFromTemplateContent,
) *CompileContextData {
	templateContext := templateContextMock{
		nil,
		templateContextRenderFromTemplateContent,
	}
	annotation := mockAnnotation{
		annotation: annotation{},
	}
	compileContext := NewCompiler(
		&templateContext,
		[]AnnotationProcessorInterface{&annotation},
	)
	templateName := "fakeTemplate"
	compileContextData, _ := compileContext.Init(
		&render.TemplateContextData{
			TemplateContext: &templateContext,
			TemplateName:    &templateName,
			Template:        nil,
			RootData:        nil,
			Data:            nil,
		},
		&model.CompilerConfig{ //nolint:exhaustruct // test
			KeepIntermediateFiles: false,
		},
	)
	return compileContextData
}

func TestCompileEmptyCode(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"",
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionNotFound(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"MyPackage::function",
	)
	assert.Error(t, err, "function not found: MyPackage::function in any srcDirs []")
	assert.Equal(t, "", resultCode)
}

func TestCompileDuplicatedFunctionDirective(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			code string,
		) (string, error) {
			return code, nil
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{}
	compilerContextData.config.SrcDirs = []string{
		"./testdata",
	}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"# FUNCTIONS\n# FUNCTIONS\nMyPackage::function",
	)
	assert.Error(t, err, "duplicated FUNCTIONS directive on line 2")
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionIgnoredFunction(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			_ string,
		) (string, error) {
			return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{
		"Ignore::ignoredFunction",
	}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"Ignore::ignoredFunction",
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "Ignore::ignoredFunction\n", resultCode)
}

func TestCompileOneFunctionFound(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			code string,
		) (string, error) {
			return code, nil
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{}
	compilerContextData.config.SrcDirs = []string{
		"./testdata",
	}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"# FUNCTIONS\nMyPackage::function",
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileOneFunctionFound.txt")
}

func TestCompileOneFunctionFoundWith_AndZZZ(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			code string,
		) (string, error) {
			return code, nil
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{}
	compilerContextData.config.SrcDirs = []string{
		"./testdata",
	}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"# FUNCTIONS\nMyCompletePackage::function",
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileOneFunctionFoundWith_AndZZZ.txt")
}

func TestCompileDependentFunction(t *testing.T) {
	compilerContextData := newCompiler(
		func(
			_ *render.TemplateContextData,
			code string,
		) (string, error) {
			return code, nil
		},
	)
	compilerContextData.config.FunctionsIgnoreRegexpList = []string{}
	compilerContextData.config.SrcDirs = []string{
		"./testdata",
	}
	resultCode, err := compilerContextData.compileContext.Compile(
		compilerContextData,
		"# FUNCTIONS\nMyPackage::useDependentFunction",
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileDependentFunction.txt")
}
