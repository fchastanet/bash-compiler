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

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/render"
	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render/functions"
	"github.com/fchastanet/bash-compiler/internal/utils"
)

var (
	// shebangRegexp            = regexp.MustCompile(`^[ \t]*(#!.*)?$`)
	functionsDirectiveRegexp    = regexp.MustCompile(`^# FUNCTIONS$`)
	commentRegexp               = regexp.MustCompile(`^[[:blank:]]*(#.*)?$`)
	bashFrameworkFunctionRegexp = regexp.MustCompile(
		`(?P<funcName>([A-Z]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))`)
)

type InsertPosition int8

const (
	InsertPositionFirst  InsertPosition = 0
	InsertPositionMiddle InsertPosition = 1
	InsertPositionLast   InsertPosition = 2
)

type functionInfoStruct struct {
	FunctionName         string
	SrcFile              string // "" if not computed yet
	Inserted             bool
	InsertPosition       InsertPosition
	SourceCode           string // the src file content
	SourceCodeLoaded     bool
	SourceCodeAsTemplate bool
}

var errFunctionNotFound = errors.New("Function not found")

func ErrFunctionNotFound(functionName string, srcDirs []string) error {
	return fmt.Errorf("%w: %s in any srcDirs %v", errFunctionNotFound, functionName, srcDirs)
}

var errDuplicatedFunctionsDirective = errors.New("Duplicated FUNCTIONS directive")

func ErrDuplicatedFunctionsDirective() error {
	return fmt.Errorf("%w", errDuplicatedFunctionsDirective)
}

// Compile generates code from given model
func Compile(code string, templateContext *render.Context, binaryModel model.BinaryModel) (codeCompiled string, err error) {
	functionsMap := make(map[string]functionInfoStruct)
	extractUniqueFrameworkFunctions(functionsMap, code)
	_, err = retrieveEachFunctionPath(functionsMap, binaryModel.BinFile.SrcDirs)
	if err != nil {
		return "", err
	}
	newFunctionAdded := true
	for newFunctionAdded {
		newFunctionAdded, err = retrieveAllFunctionsContent(functionsMap, binaryModel)
		if err != nil {
			return "", err
		}
	}

	err = renderEachFunctionAsTemplate(functionsMap, templateContext)
	if err != nil {
		return "", err
	}

	functionsCode, err := generateFunctionCode(functionsMap)
	if err != nil {
		return "", err
	}

	generatedCode, err := injectFunctionCode(code, functionsCode)
	if err != nil {
		return "", err
	}

	return generatedCode, nil
}

func renderEachFunctionAsTemplate(
	functionsMap map[string]functionInfoStruct,
	templateContext *render.Context,
) (err error) {
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
		if functionInfo.SourceCodeAsTemplate || !functionInfo.SourceCodeLoaded {
			continue
		}
		if functionInfo.SourceCode != "" {
			slog.Info("renderEachFunctionAsTemplate", "functionName", functionName)
			newCode, err := myTemplateFunctions.RenderFromTemplateContent(templateContext, functionInfo.SourceCode)
			if err != nil {
				return err
			}
			slog.Info("renderEachFunctionAsTemplate", "functionName", functionName, "code", newCode)
			functionInfo.SourceCode = newCode
		}
		functionInfo.SourceCodeAsTemplate = true
		functionsMap[functionName] = functionInfo
	}
	return nil
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

func generateFunctionCode(functionsMap map[string]functionInfoStruct) (code string, err error) {
	var functionNames []string = utils.MapKeys(functionsMap)
	sort.Strings(functionNames) // ensure to generate functions always in the same order
	var bufferFirst bytes.Buffer
	var bufferMiddle bytes.Buffer
	var bufferLast bytes.Buffer
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
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

func retrieveAllFunctionsContent(functionsMap map[string]functionInfoStruct, binaryModel model.BinaryModel) (
	newFunctionAdded bool, err error,
) {
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
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
		functionInfo.SourceCode = string(fileContent)
		functionInfo.SourceCodeLoaded = true
		functionsMap[functionName] = functionInfo
		newFunctionExtracted := extractUniqueFrameworkFunctions(functionsMap, functionInfo.SourceCode)
		addedFiles, err := retrieveEachFunctionPath(functionsMap, binaryModel.BinFile.SrcDirs)
		newFunctionAdded = newFunctionAdded || addedFiles || newFunctionExtracted
		if err != nil {
			return newFunctionAdded, err
		}
	}
	return newFunctionAdded, nil
}

func retrieveEachFunctionPath(functionsMap map[string]functionInfoStruct, srcDirs []string) (
	addedFiles bool, err error) {
	addedFiles = false
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
		if functionInfo.SrcFile != "" {
			continue
		}
		functionRelativePath := convertFunctionNameToPath(functionName)
		filePath, _, found := findFileInSrcDirs(functionRelativePath, srcDirs)
		if !found {
			return addedFiles, ErrFunctionNotFound(functionName, srcDirs)
		}
		functionInfo.SrcFile = filePath
		functionsMap[functionName] = functionInfo

		// compute relative filepath
		relativeFilePathDir := filepath.Dir(functionRelativePath)

		// check if _.sh in directory of the function is needed to be loaded
		underscoreShFile := filepath.Join(relativeFilePathDir, "_.sh")
		filePath, _, found = findFileInSrcDirs(underscoreShFile, srcDirs)
		if found {
			if _, ok := functionsMap[filePath]; !ok {
				slog.Debug("Adding file", "file", filePath)
				addedFiles = true
				functionsMap[filePath] = functionInfoStruct{
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
		filePath, _, found = findFileInSrcDirs(zzzShFile, srcDirs)
		if found {
			if _, ok := functionsMap[filePath]; !ok {
				addedFiles = true
				slog.Debug("Adding file", "file", filePath)
				functionsMap[filePath] = functionInfoStruct{
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
	slog.Info("Found these", "bashFrameworkFunctions", utils.MapKeys(functionsMap))
	return addedFiles, nil
}

func extractUniqueFrameworkFunctions(functionsMap map[string]functionInfoStruct, code string) (newFunctionAdded bool) {
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
			if _, keyExists := functionsMap[funcName]; !keyExists {
				slog.Debug("Found new", "bashFrameworkFunction", funcName)
				functionsMap[funcName] = functionInfoStruct{
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
func findFileInSrcDirs(relativeFilePath string, srcDirs []string) (
	filePath string, srcDir string, found bool,
) {
	for _, srcDir := range srcDirs {
		srcFile := filepath.Join(srcDir, relativeFilePath)
		srcFileExpanded := os.ExpandEnv(srcFile)
		slog.Debug("Check if file exists", "srcDir", srcDir, "file", srcFile, "fileExpanded", srcFileExpanded)
		err := utils.FileExists(srcFileExpanded)
		if err == nil {
			return srcFileExpanded, srcDir, true
		}
	}
	return "", "", false
}

func convertFunctionNameToPath(functionName string) string {
	return strings.ReplaceAll(functionName, "::", "/") + ".sh"
}
