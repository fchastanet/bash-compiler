// Package compiler
package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/bash"
	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
)

const (
	LogFieldSourceCodeLen    = "sourceCodeLen"
	LogFieldInsertPosition   = "insertPosition"
	LogFieldCode             = "code"
	LogFieldLength           = "length"
	LogFieldSourceCodeLoaded = "sourceCodeLoaded"
)

var removeTrailingSpacesRegexp = regexp.MustCompile(
	`[ \t]+$`,
)

type functionNotFoundError struct {
	error
	FunctionName string
	SrcDirs      []string
}

func (e *functionNotFoundError) Error() string {
	return fmt.Sprintf(
		"function not found: %s in any srcDirs %v",
		e.FunctionName,
		e.SrcDirs,
	)
}

type annotationCastError struct {
	error
	FunctionName string
}

func (e *annotationCastError) Error() string {
	return "cannot cast annotation on function: " + e.FunctionName
}

type duplicatedFunctionsDirectiveError struct {
	error
	LineNumber int
}

func (e *duplicatedFunctionsDirectiveError) Error() string {
	return fmt.Sprintf("duplicated FUNCTIONS directive on line %d", e.LineNumber)
}

type InsertPosition int8

const (
	InsertPositionFirst  InsertPosition = 0
	InsertPositionMiddle InsertPosition = 1
	InsertPositionLast   InsertPosition = 2
)

type AnnotationProcessorInterface interface {
	GetTitle() string
	Init(compileContextData *CompileContextData) error
	Reset()
	ParseFunction(compileContextData *CompileContextData, functionStruct *functionInfoStruct) error
	Process(compileContextData *CompileContextData) error
	PostProcess(compileContextData *CompileContextData, code string) (newCode string, err error)
}

type annotationProcessor struct{}

type annotation struct{}

type functionInfoStruct struct {
	FunctionName         string
	SrcFile              string // "" if not computed yet
	SourceCode           string // the src file content
	AnnotationMap        map[string]interface{}
	Inserted             bool
	InsertPosition       InsertPosition
	SourceCodeLoaded     bool
	SourceCodeAsTemplate bool
}

type CompileContext struct {
	templateContext      render.TemplateContextInterface
	annotationProcessors []AnnotationProcessorInterface
}

type CompileContextData struct {
	compileContext        *CompileContext
	templateContextData   *render.TemplateContextData
	config                *model.CompilerConfig
	functionsMap          map[string]functionInfoStruct
	ignoreFunctionsRegexp []*regexp.Regexp
}

func compilerValidationError(fieldName string, fieldValue any) error {
	return &errors.ValidationError{
		InnerError: nil,
		Context:    "compiler",
		FieldName:  fieldName,
		FieldValue: fieldValue,
	}
}

func (context CompileContextData) Validate() error {
	if context.compileContext == nil {
		return compilerValidationError("CompileContextData.compileContext", context.compileContext)
	}
	if context.templateContextData == nil {
		return compilerValidationError("CompileContextData.templateContextData", context.templateContextData)
	}
	if context.config == nil {
		return compilerValidationError("CompileContextData.config", context.config)
	}
	if context.functionsMap == nil {
		return compilerValidationError("CompileContextData.functionsMap", context.functionsMap)
	}

	return nil
}

// Compile generates code from given model
func NewCompiler(
	templateContext render.TemplateContextInterface,
	annotationProcessors []AnnotationProcessorInterface,
) CompileContext {
	return CompileContext{
		templateContext:      templateContext,
		annotationProcessors: annotationProcessors,
	}
}

func (context CompileContext) Init(
	templateContextData *render.TemplateContextData,
	config *model.CompilerConfig,
) (*CompileContextData, error) {
	compileContextData := &CompileContextData{
		compileContext:        &context,
		templateContextData:   templateContextData,
		config:                config,
		functionsMap:          make(map[string]functionInfoStruct),
		ignoreFunctionsRegexp: nil,
	}
	for _, annotationProcessor := range context.annotationProcessors {
		err := annotationProcessor.Init(compileContextData)
		if logger.FancyHandleError(err) {
			return nil, err
		}
	}

	return compileContextData, nil
}

func (context CompileContext) Compile(
	compileContextData *CompileContextData, code string,
) (codeCompiled string, err error) {
	_, err = context.computeFunctions(compileContextData, code)
	if err != nil {
		return "", err
	}
	compileContextData.config.DebugCopyGeneratedFile(code, "-compiler::Compile1")

	context.markAllFunctionsAsNotInserted(compileContextData)
	_, generatedCode, err := context.generateCode(compileContextData, code)
	if err != nil {
		return "", err
	}
	compileContextData.config.DebugCopyGeneratedFile(generatedCode, "-compiler::Compile2")

	for _, annotationProcessor := range context.annotationProcessors {
		annotationProcessor.Reset()
		generatedCode, err := annotationProcessor.PostProcess(compileContextData, generatedCode)
		if err != nil {
			return "", err
		}
		compileContextData.config.DebugCopyGeneratedFile(generatedCode, "-after-"+annotationProcessor.GetTitle())
	}

	return context.formatCode(generatedCode), nil
}

func (context CompileContext) computeFunctions(
	compileContextData *CompileContextData, code string,
) (codeCompiled string, err error) {
	err = context.functionsAnalysis(compileContextData, code)
	if err != nil {
		return "", err
	}

	compileContextData.config.DebugCopyGeneratedFile(code, "-compiler::computeFunctions1")

	needAnotherCompilerPass, generatedCode, err := context.generateCode(compileContextData, code)
	if err != nil {
		return "", err
	}
	compileContextData.config.DebugCopyGeneratedFile(generatedCode, "-compiler::computeFunctions2")

	if needAnotherCompilerPass {
		generatedCode, err = context.computeFunctions(compileContextData, generatedCode)
		if err != nil {
			return "", err
		}
	}

	return generatedCode, nil
}

func (context CompileContext) markAllFunctionsAsNotInserted(
	compileContextData *CompileContextData,
) {
	functionNames := getSortedFunctionNamesFromMap(compileContextData.functionsMap)
	for _, functionName := range functionNames {
		functionInfoStruct := compileContextData.functionsMap[functionName]
		functionInfoStruct.Inserted = false
		compileContextData.functionsMap[functionName] = functionInfoStruct
	}
}

func (context CompileContext) formatCode(code string) string {
	var newCodeBuffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := scanner.Text()
		line = removeTrailingSpacesRegexp.ReplaceAllLiteralString(line, "")
		newCodeBuffer.Write([]byte(line))
		newCodeBuffer.WriteByte('\n')
	}
	return newCodeBuffer.String()
}

func (context CompileContext) generateCode(compileContextData *CompileContextData, code string) (
	needAnotherCompilerPass bool,
	generatedCode string,
	err error,
) {
	functionsCode, err := context.generateFunctionCode(compileContextData)
	if err != nil {
		return false, "", err
	}
	compileContextData.config.DebugCopyGeneratedFile(functionsCode, "-compiler::generateCode1")

	generatedCode, err = injectFunctionCode(code, functionsCode)
	if err != nil {
		return false, "", err
	}
	compileContextData.config.DebugCopyGeneratedFile(generatedCode, "-compiler::generateCode2")

	newCode := generatedCode
	for _, annotationProcessor := range context.annotationProcessors {
		annotationProcessor.Reset()
		newCode, err = annotationProcessor.PostProcess(compileContextData, newCode)
		if err != nil {
			return false, "", err
		}
		compileContextData.config.DebugCopyGeneratedFile(newCode, "-after-"+annotationProcessor.GetTitle())
	}

	return newCode != generatedCode, newCode, nil
}

func (context CompileContext) functionsAnalysis(
	compileContextData *CompileContextData,
	code string,
) (err error) {
	context.extractUniqueFrameworkFunctions(compileContextData, code)
	_, err = context.retrieveEachFunctionPath(compileContextData)
	if err != nil {
		return err
	}
	newFunctionAdded := true
	for newFunctionAdded {
		newFunctionAdded, err = context.retrieveAllFunctionsContent(compileContextData)
		if err != nil {
			return err
		}
	}

	err = context.renderEachFunctionAsTemplate(compileContextData)
	if err != nil {
		return err
	}

	for _, annotationProcessor := range context.annotationProcessors {
		err = annotationProcessor.Process(compileContextData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context CompileContext) renderEachFunctionAsTemplate(
	compileContextData *CompileContextData,
) (err error) {
	functionNames := getSortedFunctionNamesFromMap(compileContextData.functionsMap)
	for _, functionName := range functionNames {
		functionInfo := compileContextData.functionsMap[functionName]
		if functionInfo.SourceCodeAsTemplate || !functionInfo.SourceCodeLoaded {
			continue
		}
		if functionInfo.SourceCode != "" {
			slog.Debug("renderEachFunctionAsTemplate", logger.LogFieldFunc, functionName)
			newCode, err := context.templateContext.RenderFromTemplateContent(
				compileContextData.templateContextData,
				functionInfo.SourceCode,
			)
			if err != nil {
				return err
			}
			slog.Debug("renderEachFunctionAsTemplate",
				logger.LogFieldFunc, functionName,
				LogFieldCode, newCode,
			)
			functionInfo.SourceCode = newCode
			for _, annotationProcessor := range context.annotationProcessors {
				err = annotationProcessor.ParseFunction(compileContextData, &functionInfo)
				if err != nil {
					return err
				}
			}
		}
		functionInfo.SourceCodeAsTemplate = true
		compileContextData.functionsMap[functionName] = functionInfo
	}
	return nil
}

func (context CompileContext) isNonFrameworkFunction(
	compileContextData *CompileContextData,
	functionName string,
) bool {
	context.nonFrameworkFunctionRegexpCompile(compileContextData)
	for _, re := range compileContextData.ignoreFunctionsRegexp {
		if re.MatchString(functionName) {
			return true
		}
	}

	return false
}

func (context CompileContext) nonFrameworkFunctionRegexpCompile(
	compileContextData *CompileContextData,
) {
	if compileContextData.ignoreFunctionsRegexp != nil {
		return
	}
	compileContextData.ignoreFunctionsRegexp = []*regexp.Regexp{}
	for _, reg := range compileContextData.config.FunctionsIgnoreRegexpList {
		re, err := regexp.Compile(reg)
		if err != nil {
			slog.Warn("ignored invalid regexp",
				logger.LogFieldVariableValue, reg,
				logger.LogFieldErr, err,
			)
		} else {
			compileContextData.ignoreFunctionsRegexp = append(compileContextData.ignoreFunctionsRegexp, re)
		}
	}
}

func (context CompileContext) generateFunctionCode(
	compileContextData *CompileContextData,
) (
	code string,
	err error,
) {
	functionNames := getSortedFunctionNamesFromMap(compileContextData.functionsMap)
	sort.Strings(functionNames) // ensure to generate functions always in the same order

	var finalBuffer bytes.Buffer
	err = context.insertFunctionsCode(compileContextData, functionNames, &finalBuffer, InsertPositionFirst)
	if err != nil {
		return "", err
	}
	err = context.insertFunctionsCode(compileContextData, functionNames, &finalBuffer, InsertPositionMiddle)
	if err != nil {
		return "", err
	}
	err = context.insertFunctionsCode(compileContextData, functionNames, &finalBuffer, InsertPositionLast)
	if err != nil {
		return "", err
	}
	slog.Debug("Final Buffer length", LogFieldLength, finalBuffer.Len())
	return finalBuffer.String(), nil
}

func (context CompileContext) insertFunctionsCode(
	compileContextData *CompileContextData,
	functionNames []string,
	buffer *bytes.Buffer,
	insertPosition InsertPosition,
) error {
	for _, functionName := range functionNames {
		functionInfo := compileContextData.functionsMap[functionName]
		if functionInfo.Inserted || functionInfo.InsertPosition != insertPosition {
			continue
		}
		if !functionInfo.SourceCodeLoaded {
			slog.Warn("Function source code not loaded", logger.LogFieldFunc, functionName)
			continue
		}
		slog.Debug("Append",
			LogFieldSourceCodeLen, len(functionInfo.SourceCode),
			LogFieldInsertPosition, functionInfo.InsertPosition,
		)
		_, err := buffer.WriteString(functionInfo.SourceCode)
		if err != nil {
			return err
		}
		functionInfo.Inserted = true
		compileContextData.functionsMap[functionName] = functionInfo
	}

	return nil
}

func (context CompileContext) retrieveAllFunctionsContent(
	compileContextData *CompileContextData,
) (
	newFunctionAdded bool, err error,
) {
	functionNames := getSortedFunctionNamesFromMap(compileContextData.functionsMap)
	for _, functionName := range functionNames {
		if context.isNonFrameworkFunction(compileContextData, functionName) {
			continue
		}
		functionInfo := compileContextData.functionsMap[functionName]
		slog.Debug(
			"retrieveAllFunctionsContent",
			logger.LogFieldFunc, functionName,
			LogFieldSourceCodeLoaded, functionInfo.SourceCodeLoaded,
		)
		if functionInfo.SourceCodeLoaded {
			slog.Debug("Function source code loaded", logger.LogFieldFunc, functionName)
			continue
		}
		slog.Debug("Loading Function source code from file",
			logger.LogFieldFunc, functionName,
			logger.LogFieldFilePath, functionInfo.SrcFile,
		)
		fileContent, err := os.ReadFile(functionInfo.SrcFile)
		if err != nil {
			return false, err
		}
		functionInfo.SourceCode = bash.RemoveFirstShebangLineIfAny(string(fileContent))
		functionInfo.SourceCodeLoaded = true
		compileContextData.functionsMap[functionName] = functionInfo
		newFunctionExtracted := context.extractUniqueFrameworkFunctions(
			compileContextData,
			functionInfo.SourceCode,
		)
		addedFiles, err := context.retrieveEachFunctionPath(compileContextData)
		newFunctionAdded = newFunctionAdded || addedFiles || newFunctionExtracted
		if err != nil {
			return newFunctionAdded, err
		}
	}
	return newFunctionAdded, nil
}

func createFunctionInfoStruct(
	funcName string, srcFile string, insertPosition InsertPosition,
) functionInfoStruct {
	return functionInfoStruct{
		FunctionName:         funcName,
		SrcFile:              srcFile,
		Inserted:             false,
		InsertPosition:       insertPosition,
		SourceCode:           "",
		SourceCodeLoaded:     false,
		SourceCodeAsTemplate: false,
		AnnotationMap:        make(map[string]interface{}),
	}
}

func getSortedFunctionNamesFromMap(myMap map[string]functionInfoStruct) []string {
	functionNames := structures.MapKeys(myMap)
	sort.Strings(functionNames)
	return functionNames
}

func (context CompileContext) retrieveEachFunctionPath(
	compileContextData *CompileContextData,
) (
	addedFiles bool, err error,
) {
	addedFiles = false
	functionNames := getSortedFunctionNamesFromMap(compileContextData.functionsMap)
	for _, functionName := range functionNames {
		if context.isNonFrameworkFunction(compileContextData, functionName) {
			continue
		}
		functionInfo := compileContextData.functionsMap[functionName]
		if functionInfo.SrcFile != "" {
			continue
		}
		functionRelativePath := convertFunctionNameToPath(functionName)
		filePath, found := context.findFileInSrcDirs(compileContextData, functionRelativePath)
		if !found {
			return addedFiles, &functionNotFoundError{nil, functionName, compileContextData.config.SrcDirsExpanded}
		}
		functionInfo.SrcFile = filePath
		compileContextData.functionsMap[functionName] = functionInfo

		// compute relative filepath
		relativeFilePathDir := filepath.Dir(functionRelativePath)

		// check if _.sh in directory of the function is needed to be loaded
		underscoreShFile := filepath.Join(relativeFilePathDir, "_.sh")
		filePath, found = context.findFileInSrcDirs(compileContextData, underscoreShFile)
		if found {
			if _, ok := compileContextData.functionsMap[filePath]; !ok {
				slog.Debug("Adding file", logger.LogFieldFilePath, filePath)
				addedFiles = true
				compileContextData.functionsMap[filePath] = createFunctionInfoStruct(
					filePath, filePath, InsertPositionFirst,
				)
			}
		}

		// check if ZZZ.sh in directory of the function is needed to be loaded
		zzzShFile := filepath.Join(relativeFilePathDir, "ZZZ.sh")
		filePath, found = context.findFileInSrcDirs(compileContextData, zzzShFile)
		if found {
			if _, ok := compileContextData.functionsMap[filePath]; !ok {
				addedFiles = true
				slog.Debug("Adding file", logger.LogFieldFilePath, filePath)
				compileContextData.functionsMap[filePath] = createFunctionInfoStruct(
					filePath, filePath, InsertPositionLast,
				)
			}
		}
	}

	// TODO https://go.dev/play/p/0yJNk065ftB to format functionMap as json
	slog.Debug("Found these",
		logger.LogFieldVariableName, "bashFrameworkFunctions",
		logger.LogFieldVariableValue, getSortedFunctionNamesFromMap(compileContextData.functionsMap),
	)
	return addedFiles, nil
}

func (context CompileContext) extractUniqueFrameworkFunctions(
	compileContextData *CompileContextData,
	code string,
) (newFunctionAdded bool) {
	if code == "" {
		return false
	}
	var rewrittenCode bytes.Buffer
	newFunctionAdded = false
	funcNameGroupIndex := bashFrameworkFunctionRegexp.SubexpIndex("funcName")
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := scanner.Bytes()
		rewrittenCode.Write(line)
		if IsCommentLine(line) {
			continue
		}
		matches := bashFrameworkFunctionRegexp.FindSubmatch(line)
		if matches != nil {
			funcName := string(matches[funcNameGroupIndex])
			if _, keyExists := compileContextData.functionsMap[funcName]; !keyExists {
				slog.Debug("Found new",
					logger.LogFieldVariableName, "bashFrameworkFunction",
					logger.LogFieldVariableValue, funcName,
				)
				if context.isNonFrameworkFunction(compileContextData, funcName) {
					continue
				}

				compileContextData.functionsMap[funcName] = createFunctionInfoStruct(
					funcName, "", InsertPositionMiddle,
				)
				newFunctionAdded = true
			}
		}
	}

	return newFunctionAdded
}

func (context CompileContext) findFileInSrcDirs(
	compileContextData *CompileContextData,
	relativeFilePath string,
) (
	filePath string, found bool,
) {
	for _, srcDir := range compileContextData.config.SrcDirs {
		srcFile := filepath.Join(srcDir, relativeFilePath)
		srcFileExpanded := os.ExpandEnv(srcFile)
		slog.Debug(
			"Check if file exists",
			logger.LogFieldDirPath, srcDir,
			logger.LogFieldFilePath, srcFile,
			logger.LogFieldFilePathExpanded, srcFileExpanded,
		)
		err := files.FileExists(srcFileExpanded)
		if err == nil {
			return srcFileExpanded, true
		}
	}
	return "", false
}

func convertFunctionNameToPath(functionName string) string {
	return strings.ReplaceAll(functionName, "::", "/") + ".sh"
}

func injectFunctionCode(code string, functionsCode string) (newCode string, err error) {
	var rewrittenCode bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(code))
	slog.Debug("debugCode", LogFieldCode, code)
	functionDirectiveFound := false
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		lineNumber++
		if IsFunctionDirective(line) {
			if functionDirectiveFound {
				return "", &duplicatedFunctionsDirectiveError{nil, lineNumber}
			}
			rewrittenCode.Write([]byte(functionsCode))
			functionDirectiveFound = true
		}

		rewrittenCode.Write(line)
		rewrittenCode.WriteByte(byte('\n'))
	}
	return rewrittenCode.String(), nil
}
