package services

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fchastanet/bash-compiler/internal/compiler"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
)

type CodeGeneratorInterface interface {
	GenerateCode() (codeCompiled string, err error)
}

type CodeCompilerInterface interface {
	Init() error
	Compile(code string) (codeCompiled string, err error)
	GenerateCode(code string) (generatedCode string, err error)
}

type BinaryModelServiceContext struct {
	targetDir             string
	keepIntermediateFiles bool
	binaryModelFilePath   string
	binaryModelBaseName   string
	binaryModel           *model.BinaryModel
	binaryModelContext    *model.BinaryModelContext
	templateContext       *render.Context
	codeCompiler          CodeCompilerInterface
}

func NewBinaryModelService(
	targetDir string,
	keepIntermediateFiles bool,
	binaryModelFilePath string,
) (_ *BinaryModelServiceContext) {
	return &BinaryModelServiceContext{
		targetDir:             targetDir,
		keepIntermediateFiles: keepIntermediateFiles,
		binaryModelFilePath:   binaryModelFilePath,
	}
}

func (binaryModelServiceContext *BinaryModelServiceContext) Init() error {
	binaryModelBaseName := files.BaseNameWithoutExtension(binaryModelServiceContext.binaryModelFilePath)
	binaryModelServiceContext.binaryModelBaseName = binaryModelBaseName
	referenceDir := filepath.Dir(binaryModelServiceContext.binaryModelFilePath)

	// init binary Model
	binaryModelContext := model.NewBinaryModel(
		binaryModelServiceContext.targetDir,
		binaryModelServiceContext.binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		binaryModelServiceContext.keepIntermediateFiles,
	)
	binaryModelServiceContext.binaryModelContext = binaryModelContext
	binaryModel, err := binaryModelServiceContext.binaryModelContext.Load()
	if err != nil {
		return err
	}
	binaryModelServiceContext.binaryModel = binaryModel

	// init template context
	data := make(map[string]interface{})
	data["binData"] = binaryModel.BinData
	data["compilerConfig"] = binaryModel.CompilerConfig
	data["vars"] = binaryModel.Vars
	templateDirs := structures.ExpandStringList(binaryModel.CompilerConfig.TemplateDirs)

	templateContext := render.NewTemplateContext(
		templateDirs,
		binaryModel.CompilerConfig.TemplateFile,
		data,
	)
	err = templateContext.Init(render.FuncMap())
	if err != nil {
		return err
	}
	binaryModelServiceContext.templateContext = templateContext

	// init code compiler
	codeCompiler := compiler.NewCompiler(
		templateContext,
		binaryModel.CompilerConfig,
	)
	binaryModelServiceContext.codeCompiler = CodeCompilerInterface(codeCompiler)
	err = binaryModelServiceContext.codeCompiler.Init()
	if logger.FancyHandleError(err) {
		return err
	}

	return nil
}

func (binaryModelServiceContext *BinaryModelServiceContext) Compile() error {
	codeCompiled, err := binaryModelServiceContext.renderCode()
	if logger.FancyHandleError(err) {
		return err
	}

	// Save resulting file
	targetFile := structures.ExpandStringValue(
		binaryModelServiceContext.binaryModel.CompilerConfig.TargetFile,
	)

	err = os.WriteFile(targetFile, []byte(codeCompiled), files.UserReadWriteExecutePerm)
	if logger.FancyHandleError(err) {
		return err
	}
	slog.Info("Compiled", logger.LogFieldFilePath, targetFile)

	return nil
}

func (binaryModelServiceContext *BinaryModelServiceContext) renderBinaryCodeFromTemplate() (codeCompiled string, err error) {
	// Render code using template
	code, err := (*binaryModelServiceContext.templateContext).RenderFromTemplateName()
	if err != nil {
		return "", err
	}
	if binaryModelServiceContext.keepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			binaryModelServiceContext.binaryModelContext.TargetDir,
			binaryModelServiceContext.binaryModelBaseName,
			"-3-afterTemplateRendering.sh",
			code,
		)
		if err != nil {
			return "", err
		}
	}
	return code, err
}

func (binaryModelServiceContext *BinaryModelServiceContext) renderCode() (codeCompiled string, err error) {
	code, err := binaryModelServiceContext.renderBinaryCodeFromTemplate()
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Compile to get functions loaded once
	_, err = binaryModelServiceContext.codeCompiler.Compile(code)
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Generate code with all functions that has been loaded
	codeCompiled, err = binaryModelServiceContext.codeCompiler.GenerateCode(code)
	if err != nil {
		return "", err
	}

	return codeCompiled, nil
}
