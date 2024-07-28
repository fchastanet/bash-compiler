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
	Init(
		templateContextData *render.TemplateContextData,
		config model.CompilerConfig,
	) (*compiler.CompileContextData, error)
	Compile(compileContextData *compiler.CompileContextData, code string) (codeCompiled string, err error)
	GenerateCode(compileContextData *compiler.CompileContextData, code string) (generatedCode string, err error)
}

type BinaryModelInterface interface {
	Load(
		targetDir string,
		binaryModelFilePath string,
		binaryModelBaseName string,
		referenceDir string,
		keepIntermediateFiles bool,
	) (binaryModel *model.BinaryModel, err error)
}

type BinaryModelServiceContext struct {
	binaryModel     BinaryModelInterface
	templateContext *render.TemplateContext
	codeCompiler    CodeCompilerInterface
}

type BinaryModelServiceContextData struct {
	binaryModel           *model.BinaryModel
	compileContextData    *compiler.CompileContextData
	templateContextData   *render.TemplateContextData
	targetDir             string
	keepIntermediateFiles bool
	binaryModelFilePath   string
	binaryModelBaseName   string
}

func NewBinaryModelService(
	binaryModel BinaryModelInterface,
	templateContext *render.TemplateContext,
	codeCompiler CodeCompilerInterface,
) (_ *BinaryModelServiceContext) {
	return &BinaryModelServiceContext{
		binaryModel:     binaryModel,
		templateContext: templateContext,
		codeCompiler:    codeCompiler,
	}
}

func (binaryModelServiceContext *BinaryModelServiceContext) Init(
	targetDir string,
	keepIntermediateFiles bool,
	binaryModelFilePath string,
) (*BinaryModelServiceContextData, error) {
	binaryModelServiceContextData := &BinaryModelServiceContextData{
		targetDir:             targetDir,
		keepIntermediateFiles: keepIntermediateFiles,
		binaryModelFilePath:   binaryModelFilePath,
	}
	binaryModelBaseName := files.BaseNameWithoutExtension(binaryModelServiceContextData.binaryModelFilePath)
	binaryModelServiceContextData.binaryModelBaseName = binaryModelBaseName
	referenceDir := filepath.Dir(binaryModelServiceContextData.binaryModelFilePath)

	// init binary Model
	binaryModel, err := binaryModelServiceContext.binaryModel.Load(
		binaryModelServiceContextData.targetDir,
		binaryModelServiceContextData.binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		binaryModelServiceContextData.keepIntermediateFiles,
	)
	if err != nil {
		return nil, err
	}
	binaryModelServiceContextData.binaryModel = binaryModel

	// init template context
	data := make(map[string]interface{})
	data["binData"] = binaryModel.BinData
	data["compilerConfig"] = binaryModel.CompilerConfig
	data["vars"] = binaryModel.Vars
	templateDirs := structures.ExpandStringList(binaryModel.CompilerConfig.TemplateDirs)

	templateContextData, err := binaryModelServiceContext.templateContext.Init(
		templateDirs,
		binaryModel.CompilerConfig.TemplateFile,
		data,
		render.FuncMap(),
	)
	if err != nil {
		return nil, err
	}
	binaryModelServiceContextData.templateContextData = templateContextData

	// init code compiler
	compileContextData, err := binaryModelServiceContext.codeCompiler.Init(
		templateContextData,
		binaryModel.CompilerConfig,
	)
	if logger.FancyHandleError(err) {
		return nil, err
	}
	binaryModelServiceContextData.compileContextData = compileContextData

	return binaryModelServiceContextData, nil
}

func (binaryModelServiceContext *BinaryModelServiceContext) Compile(
	binaryModelServiceContextData *BinaryModelServiceContextData,
) error {
	codeCompiled, err := binaryModelServiceContext.renderCode(binaryModelServiceContextData)
	if logger.FancyHandleError(err) {
		return err
	}

	// Save resulting file
	targetFile := structures.ExpandStringValue(
		binaryModelServiceContextData.binaryModel.CompilerConfig.TargetFile,
	)

	err = os.WriteFile(targetFile, []byte(codeCompiled), files.UserReadWriteExecutePerm)
	if logger.FancyHandleError(err) {
		return err
	}
	slog.Info("Compiled", logger.LogFieldFilePath, targetFile)

	return nil
}

func (binaryModelServiceContext *BinaryModelServiceContext) renderBinaryCodeFromTemplate(
	binaryModelServiceContextData *BinaryModelServiceContextData,
) (codeCompiled string, err error) {
	// Render code using template
	code, err := (*binaryModelServiceContext.templateContext).RenderFromTemplateName(
		binaryModelServiceContextData.templateContextData,
	)
	if err != nil {
		return "", err
	}
	if binaryModelServiceContextData.keepIntermediateFiles {
		err = logger.DebugCopyGeneratedFile(
			binaryModelServiceContextData.targetDir,
			binaryModelServiceContextData.binaryModelBaseName,
			"-3-afterTemplateRendering.sh",
			code,
		)
		if err != nil {
			return "", err
		}
	}
	return code, err
}

func (binaryModelServiceContext *BinaryModelServiceContext) renderCode(
	binaryModelServiceContextData *BinaryModelServiceContextData,
) (codeCompiled string, err error) {
	code, err := binaryModelServiceContext.renderBinaryCodeFromTemplate(binaryModelServiceContextData)
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Compile to get functions loaded once
	_, err = binaryModelServiceContext.codeCompiler.Compile(
		binaryModelServiceContextData.compileContextData,
		code,
	)
	if logger.FancyHandleError(err) {
		return "", err
	}

	// Generate code with all functions that has been loaded
	codeCompiled, err = binaryModelServiceContext.codeCompiler.GenerateCode(
		binaryModelServiceContextData.compileContextData,
		code,
	)
	if err != nil {
		return "", err
	}

	return codeCompiled, nil
}
