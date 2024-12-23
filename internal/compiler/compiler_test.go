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

func (*mockAnnotation) Init(
	_ *CompileContextData,
) error {
	return nil
}

func (*mockAnnotation) ParseFunction(
	_ *CompileContextData, _ *functionInfoStruct,
) error {
	return nil
}

func (*mockAnnotation) Process(
	_ *CompileContextData,
) error {
	return nil
}

func (*mockAnnotation) PostProcess(
	_ *CompileContextData, code string,
) (newCode string, err error) {
	return code, nil
}

func (*mockAnnotation) GetTitle() string {
	return ""
}

func (*mockAnnotation) Reset() {
}

// *******************************************
// newCompiler
func newMockedCompiler(
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
			IntermediateFilesDir: "",
		},
	)
	return compileContextData
}

func simulateGoodRenderingCallback(_ *render.TemplateContextData, code string) (string, error) {
	return code, nil
}

func simulateFailingRenderingCallback(_ *render.TemplateContextData, _ string) (string, error) {
	return shouldNotBeCalledCodeStr, &shouldNotBeCalledError
}

func compile(
	inputCode string,
	functionsIgnoreRegexpList []string,
	renderingCallback TemplateContextRenderFromTemplateContent,
	srcDirs []string,
) (codeCompiled string, err error) {
	compilerContextData := newMockedCompiler(renderingCallback)
	compilerContextData.config.FunctionsIgnoreRegexpList = functionsIgnoreRegexpList
	compilerContextData.config.SrcDirs = srcDirs
	return compilerContextData.compileContext.Compile(
		compilerContextData,
		inputCode,
	)
}

func TestCompileEmptyCode(t *testing.T) {
	resultCode, err := compile(
		"",
		[]string{},
		simulateFailingRenderingCallback,
		[]string{},
	)
	assert.Equal(t, nil, err)
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionNotFound(t *testing.T) {
	resultCode, err := compile(
		"MyPackage::function",
		[]string{},
		simulateFailingRenderingCallback,
		[]string{},
	)
	assert.Error(t, err, "function not found: MyPackage::function in any srcDirs []")
	assert.Equal(t, "", resultCode)
}

func TestCompileDuplicatedFunctionDirective(t *testing.T) {
	resultCode, err := compile(
		"# FUNCTIONS\n# FUNCTIONS\nMyPackage::function",
		[]string{}, simulateGoodRenderingCallback, []string{"./testdata"},
	)
	assert.Error(t, err, "duplicated FUNCTIONS directive on line 2")
	assert.Equal(t, "", resultCode)
}

func TestCompileFunctionIgnoredFunction(t *testing.T) {
	resultCode, err := compile("Ignore::ignoredFunction", []string{
		"Ignore::ignoredFunction",
	}, simulateFailingRenderingCallback, []string{"./testdata"})
	assert.Equal(t, nil, err)
	assert.Equal(t, "Ignore::ignoredFunction\n", resultCode)
}

func TestCompileOneFunctionFound(t *testing.T) {
	resultCode, err := compile(
		"# FUNCTIONS\nMyPackage::function",
		[]string{},
		simulateGoodRenderingCallback,
		[]string{"./testdata"},
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileOneFunctionFound.txt")
}

func TestCompileOneFunctionFoundWith_AndZZZ(t *testing.T) {
	resultCode, err := compile(
		"# FUNCTIONS\nMyCompletePackage::function",
		[]string{},
		simulateGoodRenderingCallback,
		[]string{"./testdata"},
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileOneFunctionFoundWith_AndZZZ.txt")
}

func TestCompileDependentFunction(t *testing.T) {
	resultCode, err := compile(
		"# FUNCTIONS\nMyPackage::useDependentFunction",
		[]string{},
		simulateGoodRenderingCallback,
		[]string{"./testdata"},
	)
	assert.Equal(t, err, nil)
	golden.Assert(t, resultCode, "expectedTestCompileDependentFunction.txt")
}
