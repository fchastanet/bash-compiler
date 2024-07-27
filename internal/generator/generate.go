package generator

import (
	"github.com/fchastanet/bash-compiler/internal/logger"
)

type CodeGeneratorInterface interface {
	GenerateCode() (codeCompiled string, err error)
}

type TemplateRenderingInterface interface {
	Init(funcMap map[string]interface{}) error
	RenderFromTemplateName() (code string, err error)
}

type CodeGeneratorContext struct {
	targetDir             string
	binaryModelBaseName   string
	keepIntermediateFiles bool
	templateContext       *TemplateRenderingInterface
}

func NewCodeGenerator(
	targetDir string,
	binaryModelBaseName string,
	templateContext *TemplateRenderingInterface,
	keepIntermediateFiles bool,
) CodeGeneratorInterface {
	return &CodeGeneratorContext{
		targetDir:             targetDir,
		binaryModelBaseName:   binaryModelBaseName,
		keepIntermediateFiles: keepIntermediateFiles,
		templateContext:       templateContext,
	}
}

func (codeGeneratorContext *CodeGeneratorContext) GenerateCode() (codeCompiled string, err error) {
	// Render code using template
	code, err := (*codeGeneratorContext.templateContext).RenderFromTemplateName()
	if err != nil {
		return "", err
	}
	if codeGeneratorContext.keepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			codeGeneratorContext.targetDir,
			codeGeneratorContext.binaryModelBaseName,
			"-3-afterTemplateRendering.sh",
			code,
		)
		if err != nil {
			return "", err
		}
	}
	return code, err
}
