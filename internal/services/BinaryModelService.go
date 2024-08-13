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

type CodeCompilerInterface interface {
	Init(
		templateContextData *render.TemplateContextData,
		config *model.CompilerConfig,
	) (*compiler.CompileContextData, error)
	Compile(compileContextData *compiler.CompileContextData, code string) (codeCompiled string, err error)
}

type BinaryModelLoaderInterface interface {
	Load(
		targetDir string,
		binaryModelFilePath string,
		binaryModelBaseName string,
		referenceDir string,
		keepIntermediateFiles bool,
	) (binaryModel *model.BinaryModel, err error)
}

type TemplateContextInterface interface {
	Init(
		templateDirs []string,
		templateFile string,
		data interface{},
		funcMap map[string]interface{},
	) (*render.TemplateContextData, error)
	Render(
		templateContextData *render.TemplateContextData,
		templateName string,
	) (string, error)
	RenderFromTemplateName(
		templateContextData *render.TemplateContextData,
	) (code string, err error)
	RenderFromTemplateContent(
		templateContextData *render.TemplateContextData,
		templateContent string,
	) (codeStr string, err error)
}

type BinaryModelServiceContext struct {
	binaryModelLoader BinaryModelLoaderInterface
	templateContext   TemplateContextInterface
	codeCompiler      CodeCompilerInterface
}

type BinaryModelServiceContextData struct {
	binaryModelData       *model.BinaryModel
	compileContextData    *compiler.CompileContextData
	templateContextData   *render.TemplateContextData
	targetDir             string
	keepIntermediateFiles bool
	binaryModelFilePath   string
	binaryModelBaseName   string
}

func NewBinaryModelService(
	binaryModelLoader BinaryModelLoaderInterface,
	templateContext TemplateContextInterface,
	codeCompiler CodeCompilerInterface,
) (_ *BinaryModelServiceContext) {
	return &BinaryModelServiceContext{
		binaryModelLoader: binaryModelLoader,
		templateContext:   templateContext,
		codeCompiler:      codeCompiler,
	}
}

func (binaryModelServiceContext *BinaryModelServiceContext) Init(
	targetDir string,
	keepIntermediateFiles bool,
	binaryModelFilePath string,
) (*BinaryModelServiceContextData, error) {
	binaryModelBaseName := files.BaseNameWithoutExtension(binaryModelFilePath)
	referenceDir := filepath.Dir(binaryModelFilePath)
	// init binary Model
	binaryModelData, err := binaryModelServiceContext.binaryModelLoader.Load(
		targetDir,
		binaryModelFilePath,
		binaryModelBaseName,
		referenceDir,
		keepIntermediateFiles,
	)
	if err != nil {
		return nil, err
	}
	binaryModelServiceContextData := &BinaryModelServiceContextData{
		binaryModelData:       binaryModelData,
		templateContextData:   nil, // computed later
		compileContextData:    nil, // computed later
		targetDir:             targetDir,
		keepIntermediateFiles: keepIntermediateFiles,
		binaryModelFilePath:   binaryModelFilePath,
		binaryModelBaseName:   binaryModelBaseName,
	}

	// init template context
	data := make(map[string]interface{})
	data["binData"] = binaryModelData.BinData
	data["compilerConfig"] = binaryModelData.CompilerConfig
	data["vars"] = binaryModelData.Vars
	templateDirs := structures.ExpandStringList(binaryModelData.CompilerConfig.TemplateDirs)

	templateContextData, err := binaryModelServiceContext.templateContext.Init(
		templateDirs,
		binaryModelData.CompilerConfig.TemplateFile,
		data,
		render.FuncMap(),
	)
	if err != nil {
		return nil, err
	}
	binaryModelServiceContextData.templateContextData = templateContextData

	// init code compiler
	compilerConfig := &binaryModelData.CompilerConfig
	compileContextData, err := binaryModelServiceContext.codeCompiler.Init(
		templateContextData,
		compilerConfig,
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
		binaryModelServiceContextData.binaryModelData.CompilerConfig.TargetFile,
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
	code, err := binaryModelServiceContext.templateContext.Render(
		binaryModelServiceContextData.templateContextData,
		*binaryModelServiceContextData.templateContextData.TemplateName,
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
	return binaryModelServiceContext.codeCompiler.Compile(
		binaryModelServiceContextData.compileContextData,
		code,
	)
}
