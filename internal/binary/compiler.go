package binary

import (
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/generator"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
)

type CompilerInterface interface {
	Compile(targetDir string,
		binaryModelFilePath string,
		binaryModelBaseName string,
		referenceDir string,
		keepIntermediateFiles bool,
	) (codeCompiled string, err error)
}

type CodeCompilerContext struct {
	BinaryModelContext *model.BinaryModelContext
	TemplateContext    *render.Context
	CodeGenerator      *generator.CodeGeneratorInterface
	CodeCompiler       *compiler.CodeCompilerInterface
}

func NewCompiler() CodeCompilerContext {
	return CodeCompilerContext{}
}

func (codeCompilerContext *CodeCompilerContext) Compile(
	targetDir string,
	binaryModelFilePath string,
	binaryModelBaseName string,
	referenceDir string,
	keepIntermediateFiles bool,
) (codeCompiled string, err error) {
	codeCompilerContext.BinaryModelContext = model.NewBinaryModel(
		targetDir,
		binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		keepIntermediateFiles,
	)
	err = codeCompilerContext.BinaryModelContext.LoadBinaryModel()
	if err != nil {
		return "", err
	}

	codeCompilerContext.TemplateContext, err = model.NewTemplateContext(*codeCompilerContext.BinaryModelContext)
	if err != nil {
		return "", err
	}

	codeGenerator := generator.NewCodeGenerator(
		binaryModelFilePath,
		targetDir,
		binaryModelBaseName,
		codeCompilerContext.TemplateContext,
		keepIntermediateFiles,
	)
	codeCompilerContext.CodeGenerator = &codeGenerator

	codeCompiler := compiler.NewCompiler(
		codeCompilerContext.TemplateContext,
		codeCompilerContext.BinaryModelContext.BinaryModel.CompilerConfig,
	)
	codeCompilerContext.CodeCompiler = &codeCompiler

	code, err := codeGenerator.GenerateCode()
	if err != nil {
		return "", err
	}

	// Compile
	codeCompiled, err = codeCompiler.Compile(code)
	if err != nil {
		return "", err
	}

	return codeCompiled, nil
}
