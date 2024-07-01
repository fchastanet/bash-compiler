// Package compiler
package compiler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fchastanet/bash-compiler/internal/code"
	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render/functions"
	"github.com/fchastanet/bash-compiler/internal/utils"
)

var (
	functionsDirectiveRegexp    = regexp.MustCompile(`^# FUNCTIONS$`)
	commentRegexp               = regexp.MustCompile(`^[[:blank:]]*(#.*)?$`)
	bashFrameworkFunctionRegexp = regexp.MustCompile(
		`(?P<funcName>([A-Z]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))`)
)

var errFunctionNotFound = errors.New("Function not found")

func ErrFunctionNotFound(functionName string, srcDirs []string) error {
	return fmt.Errorf("%w: %s in any srcDirs %v", errFunctionNotFound, functionName, srcDirs)
}

var errDuplicatedFunctionsDirective = errors.New("Duplicated FUNCTIONS directive")

func ErrDuplicatedFunctionsDirective() error {
	return fmt.Errorf("%w", errDuplicatedFunctionsDirective)
}

type InsertPosition int8

const (
	InsertPositionFirst  InsertPosition = 0
	InsertPositionMiddle InsertPosition = 1
	InsertPositionLast   InsertPosition = 2
)

type AnnotationProcessorInterface interface {
	Process(compileContext *compileContext, code string) error
}

type Annotation struct {
	kind       string
	properties map[string]interface{}
}

type functionInfoStruct struct {
	FunctionName         string
	SrcFile              string // "" if not computed yet
	Inserted             bool
	InsertPosition       InsertPosition
	SourceCode           string // the src file content
	SourceCodeLoaded     bool
	SourceCodeAsTemplate bool
	AnnotationMap        map[string]Annotation
}

type CodeCompilerInterface interface {
	Compile(code string) (codeCompiled string, err error)
}

type compileContext struct {
	templateContext       *render.Context
	functionsMap          map[string]functionInfoStruct
	ignoreFunctionsRegexp []*regexp.Regexp
	config                model.CompilerConfig
	annotationProcessors  []*AnnotationProcessorInterface
}

// Compile generates code from given model
func NewCompiler(
	templateContext *render.Context,
	config model.CompilerConfig,
) CodeCompilerInterface {
	requireProcessor := NewRequireAnnotationProcessor()
	return &compileContext{
		templateContext: templateContext,
		functionsMap:    make(map[string]functionInfoStruct),
		config:          config,
		annotationProcessors: []*AnnotationProcessorInterface{
			&requireProcessor,
		},
	}
}

func (context *compileContext) Compile(code string) (codeCompiled string, err error) {
	context.extractUniqueFrameworkFunctions(code)
	_, err = context.retrieveEachFunctionPath()
	if err != nil {
		return "", err
	}
	newFunctionAdded := true
	for newFunctionAdded {
		newFunctionAdded, err = context.retrieveAllFunctionsContent()
		if err != nil {
			return "", err
		}
	}

	err = context.renderEachFunctionAsTemplate()
	if err != nil {
		return "", err
	}

	for _, annotationProcessor := range context.annotationProcessors {
		err = (*annotationProcessor).Process(context, code)
		if err != nil {
			return "", err
		}
	}

	functionsCode, err := context.generateFunctionCode()
	if err != nil {
		return "", err
	}

	generatedCode, err := injectFunctionCode(code, functionsCode)
	if err != nil {
		return "", err
	}

	return generatedCode, nil
}

func (context *compileContext) renderEachFunctionAsTemplate() (err error) {
	var functionNames []string = utils.MapKeys(context.functionsMap)
	for _, functionName := range functionNames {
		functionInfo := context.functionsMap[functionName]
		if functionInfo.SourceCodeAsTemplate || !functionInfo.SourceCodeLoaded {
			continue
		}
		if functionInfo.SourceCode != "" {
			slog.Debug("renderEachFunctionAsTemplate", "functionName", functionName)
			newCode, err := myTemplateFunctions.RenderFromTemplateContent(
				context.templateContext, functionInfo.SourceCode)
			if err != nil {
				return err
			}
			slog.Debug("renderEachFunctionAsTemplate", "functionName", functionName, "code", newCode)
			functionInfo.SourceCode = newCode
		}
		functionInfo.SourceCodeAsTemplate = true
		context.functionsMap[functionName] = functionInfo
	}
	return nil
}

func (context *compileContext) isNonFrameworkFunction(functionName string) bool {
	context.nonFrameworkFunctionRegexpCompile()
	for _, re := range context.ignoreFunctionsRegexp {
		if re.MatchString(functionName) {
			return true
		}
	}

	return false
}

func (context *compileContext) nonFrameworkFunctionRegexpCompile() {
	if context.ignoreFunctionsRegexp != nil {
		return
	}
	context.ignoreFunctionsRegexp = []*regexp.Regexp{}
	for _, reg := range context.config.FunctionsIgnoreRegexpList {
		regStr := fmt.Sprint(reg)
		re, err := regexp.Compile(fmt.Sprint(regStr))
		if err != nil {
			slog.Warn("ignored invalid regexp", "regexp", regStr, "error", err)
		} else {
			context.ignoreFunctionsRegexp = append(context.ignoreFunctionsRegexp, re)
		}
	}
}

func (context *compileContext) generateFunctionCode() (code string, err error) {
	var functionNames []string = utils.MapKeys(context.functionsMap)
	sort.Strings(functionNames) // ensure to generate functions always in the same order
	var bufferFirst bytes.Buffer
	var bufferMiddle bytes.Buffer
	var bufferLast bytes.Buffer
	for _, functionName := range functionNames {
		functionInfo := context.functionsMap[functionName]
		if !functionInfo.SourceCodeLoaded {
			slog.Warn("Function source code not loaded", "functionName", functionName)
			continue
		}
		slog.Debug("Append ", "SourceCodeLen", len(functionInfo.SourceCode), "InsertPosition", functionInfo.InsertPosition)
		switch functionInfo.InsertPosition {
		case InsertPositionFirst:
			_, err = bufferFirst.Write([]byte(functionInfo.SourceCode))
		case InsertPositionMiddle:
			_, err = bufferMiddle.Write([]byte(functionInfo.SourceCode))
		case InsertPositionLast:
			_, err = bufferLast.Write([]byte(functionInfo.SourceCode))
		}
		if err != nil {
			return "", err
		}
	}
	var finalBuffer bytes.Buffer
	slog.Debug("Append ", "bufferFirstLen", bufferFirst.Len())
	_, err = finalBuffer.Write(bufferFirst.Bytes())
	if err != nil {
		return "", err
	}
	slog.Debug("Append ", "bufferMiddleLen", bufferMiddle.Len())
	_, err = finalBuffer.Write(bufferMiddle.Bytes())
	if err != nil {
		return "", err
	}
	slog.Debug("Append ", "bufferLastLen", bufferLast.Len())
	_, err = finalBuffer.Write(bufferLast.Bytes())
	if err != nil {
		return "", err
	}
	slog.Debug("Final Buffer ", "finalBufferLen", finalBuffer.Len())
	return finalBuffer.String(), nil
}

func (context *compileContext) retrieveAllFunctionsContent() (
	newFunctionAdded bool, err error,
) {
	var functionNames []string = utils.MapKeys(context.functionsMap)
	for _, functionName := range functionNames {
		if context.isNonFrameworkFunction(functionName) {
			continue
		}
		functionInfo := context.functionsMap[functionName]
		slog.Debug("retrieveAllFunctionsContent", "functionName", functionName, "SourceCodeLoaded", functionInfo.SourceCodeLoaded)
		if functionInfo.SourceCodeLoaded {
			slog.Debug("Function source code loaded", "functionName", functionName)
			continue
		}
		slog.Debug("Loading Function source code", "functionName", functionName, "SrcFile", functionInfo.SrcFile)
		fileContent, err := os.ReadFile(functionInfo.SrcFile)
		if err != nil {
			return false, err
		}
		functionInfo.SourceCode = code.RemoveFirstShebangLineIfAny(string(fileContent))
		functionInfo.SourceCodeLoaded = true
		context.functionsMap[functionName] = functionInfo
		newFunctionExtracted := context.extractUniqueFrameworkFunctions(functionInfo.SourceCode)
		addedFiles, err := context.retrieveEachFunctionPath()
		newFunctionAdded = newFunctionAdded || addedFiles || newFunctionExtracted
		if err != nil {
			return newFunctionAdded, err
		}
	}
	return newFunctionAdded, nil
}

func (context *compileContext) retrieveEachFunctionPath() (
	addedFiles bool, err error) {
	addedFiles = false
	var functionNames []string = utils.MapKeys(context.functionsMap)
	for _, functionName := range functionNames {
		if context.isNonFrameworkFunction(functionName) {
			continue
		}
		functionInfo := context.functionsMap[functionName]
		if functionInfo.SrcFile != "" {
			continue
		}
		functionRelativePath := convertFunctionNameToPath(functionName)
		filePath, _, found := context.findFileInSrcDirs(functionRelativePath)
		if !found {
			return addedFiles, ErrFunctionNotFound(functionName, context.config.SrcDirsExpanded)
		}
		functionInfo.SrcFile = filePath
		context.functionsMap[functionName] = functionInfo

		// compute relative filepath
		relativeFilePathDir := filepath.Dir(functionRelativePath)

		// check if _.sh in directory of the function is needed to be loaded
		underscoreShFile := filepath.Join(relativeFilePathDir, "_.sh")
		filePath, _, found = context.findFileInSrcDirs(underscoreShFile)
		if found {
			if _, ok := context.functionsMap[filePath]; !ok {
				slog.Debug("Adding file", "file", filePath)
				addedFiles = true
				context.functionsMap[filePath] = functionInfoStruct{
					FunctionName:         filePath,
					SrcFile:              filePath,
					Inserted:             false,
					InsertPosition:       InsertPositionFirst,
					SourceCode:           "",
					SourceCodeLoaded:     false,
					SourceCodeAsTemplate: false,
				}
			}
		}

		// check if ZZZ.sh in directory of the function is needed to be loaded
		zzzShFile := filepath.Join(relativeFilePathDir, "ZZZ.sh")
		filePath, _, found = context.findFileInSrcDirs(zzzShFile)
		if found {
			if _, ok := context.functionsMap[filePath]; !ok {
				addedFiles = true
				slog.Debug("Adding file", "file", filePath)
				context.functionsMap[filePath] = functionInfoStruct{
					FunctionName:         filePath,
					SrcFile:              filePath,
					Inserted:             false,
					InsertPosition:       InsertPositionLast,
					SourceCode:           "",
					SourceCodeLoaded:     false,
					SourceCodeAsTemplate: false,
				}
			}
		}
	}

	// TODO https://go.dev/play/p/0yJNk065ftB to format functionMap as json
	slog.Info("Found these", "bashFrameworkFunctions", utils.MapKeys(context.functionsMap))
	return addedFiles, nil
}

func (context *compileContext) extractUniqueFrameworkFunctions(code string) (newFunctionAdded bool) {
	var rewrittenCode bytes.Buffer
	newFunctionAdded = false
	funcNameGroupIndex := bashFrameworkFunctionRegexp.SubexpIndex("funcName")
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := scanner.Bytes()
		rewrittenCode.Write(line)
		if commentRegexp.Match(line) {
			continue
		}
		matches := bashFrameworkFunctionRegexp.FindSubmatch(line)
		if matches != nil {
			funcName := string(matches[funcNameGroupIndex])
			if _, keyExists := context.functionsMap[funcName]; !keyExists {
				slog.Debug("Found new", "bashFrameworkFunction", funcName)
				if context.isNonFrameworkFunction(funcName) {
					continue
				}

				context.functionsMap[funcName] = functionInfoStruct{
					FunctionName:         funcName,
					SrcFile:              "",
					Inserted:             false,
					InsertPosition:       InsertPositionMiddle,
					SourceCode:           "",
					SourceCodeLoaded:     false,
					SourceCodeAsTemplate: false,
				}
				newFunctionAdded = true
			}
		}
	}

	return newFunctionAdded
}

//nolint:unparam
func (context *compileContext) findFileInSrcDirs(relativeFilePath string) (
	filePath string, srcDir string, found bool,
) {
	for _, srcDir := range context.config.SrcDirs {
		srcFile := filepath.Join(srcDir, relativeFilePath)
		srcFileExpanded := os.ExpandEnv(srcFile)
		slog.Debug("Check if file exists", "srcDir", srcDir, "file", srcFile, "fileExpanded", srcFileExpanded)
		err := files.FileExists(srcFileExpanded)
		if err == nil {
			return srcFileExpanded, srcDir, true
		}
	}
	return "", "", false
}

func convertFunctionNameToPath(functionName string) string {
	return strings.ReplaceAll(functionName, "::", "/") + ".sh"
}

func injectFunctionCode(code string, functionsCode string) (newCode string, err error) {
	var rewrittenCode bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(code))
	slog.Debug("debugCode", "code", code)
	functionDirectiveFound := false
	for scanner.Scan() {
		line := scanner.Bytes()
		if functionsDirectiveRegexp.Match(line) {
			if functionDirectiveFound {
				return "", ErrDuplicatedFunctionsDirective()
			}
			rewrittenCode.Write([]byte(functionsCode))
			functionDirectiveFound = true
		}

		rewrittenCode.Write(line)
		rewrittenCode.WriteByte(byte('\n'))
	}
	return rewrittenCode.String(), nil
}
