package generator

import (
	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/render"
)

type CodeGeneratorInterface interface {
	GenerateCode() (codeCompiled string, err error)
}

type codeGeneratorContext struct {
	yamlFilePath          string
	targetDir             string
	binaryModelBaseName   string
	keepIntermediateFiles bool
	templateContext       *render.Context
}

func NewCodeGenerator(
	yamlFilePath string,
	targetDir string,
	binaryModelBaseName string,
	templateContext *render.Context,
	keepIntermediateFiles bool,
) CodeGeneratorInterface {
	return &codeGeneratorContext{
		yamlFilePath:          yamlFilePath,
		targetDir:             targetDir,
		binaryModelBaseName:   binaryModelBaseName,
		keepIntermediateFiles: keepIntermediateFiles,
		templateContext:       templateContext,
	}
}

func (codeGeneratorContext *codeGeneratorContext) GenerateCode() (codeCompiled string, err error) {
	// Render code using template
	code, err := codeGeneratorContext.templateContext.RenderFromTemplateName()
	if err != nil {
		return "", err
	}
	if codeGeneratorContext.keepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			codeGeneratorContext.targetDir,
			codeGeneratorContext.binaryModelBaseName,
			"-afterTemplateRendering.sh",
			code,
		)
		if err != nil {
			return "", err
		}
	}
	return code, err
}
