package binary

import (
	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/generator"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

func Render(
	codeGenerator generator.CodeGeneratorInterface,
	codeCompiler compiler.CodeCompilerInterface,
) (codeCompiled string, err error) {
	code, err := codeGenerator.GenerateCode()
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Compile to get functions loaded once
	_, err = codeCompiler.Compile(code)
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Generate code with all functions that has been loaded
	codeCompiled, err = codeCompiler.GenerateCode(code)
	if err != nil {
		return "", err
	}

	return codeCompiled, nil
}
