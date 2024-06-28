package binary

import (
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/generator"
	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
)

type CompilerInterface interface {
	Compile(targetDir string,
		binaryModelFilePath string,
		binaryModelBaseName string,
		referenceDir string,
		keepIntermediateFiles bool,
	) (binaryModelContext *model.BinaryModelContext, codeCompiled string, err error)
}

type codeCompiler struct {
	//nolint:unused
	binaryModel model.BinaryModelContext
	//nolint:unused
	templateContext render.Context
	//nolint:unused
	codeGenerator generator.CodeGeneratorInterface
	//nolint:unused
	codeCompiler compiler.CodeCompilerInterface
}

func NewCompiler() CompilerInterface {
	return codeCompiler{}
}

func (codeCompiler) Compile(
	targetDir string,
	binaryModelFilePath string,
	binaryModelBaseName string,
	referenceDir string,
	keepIntermediateFiles bool,

) (binaryModelContext *model.BinaryModelContext, codeCompiled string, err error) {
	binaryModelContext = model.NewBinaryModel(
		targetDir,
		binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		keepIntermediateFiles,
	)
	err = binaryModelContext.LoadBinaryModel()
	logger.Check(err)

	templateContext, err := model.NewTemplateContext(*binaryModelContext)
	logger.Check(err)

	codeGenerator := generator.NewCodeGenerator(
		binaryModelFilePath,
		targetDir,
		binaryModelBaseName,
		templateContext,
		keepIntermediateFiles,
	)

	codeCompiler := compiler.NewCompiler(
		templateContext,
		binaryModelContext.BinaryModel.CompilerConfig,
	)

	code, err := codeGenerator.GenerateCode()
	if err != nil {
		return nil, "", err
	}

	// Compile
	codeCompiled, err = codeCompiler.Compile(code)
	if err != nil {
		return nil, "", err
	}
	return binaryModelContext, codeCompiled, nil
}
