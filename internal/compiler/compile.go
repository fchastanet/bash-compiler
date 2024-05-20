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
	"strings"

	"github.com/fchastanet/bash-compiler/internal/model"
	"github.com/fchastanet/bash-compiler/internal/utils"
)

var (
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
	FunctionName     string
	SrcFile          string // "" if not computed yet
	Inserted         bool
	InsertPosition   InsertPosition
	SourceCode       string // the src file content
	SourceCodeLoaded bool
}

var errFunctionNotFound = errors.New("Function not found")

func ErrFunctionNotFound(functionName string, srcDirs []string) error {
	return fmt.Errorf("%w: %s in any srcDirs %v", errFunctionNotFound, functionName, srcDirs)
}

// Compile generates code from given model
func Compile(code string, binaryModel *model.BinaryModel) (codeCompiled string, err error) {
	functionsMap := make(map[string]functionInfoStruct)
	extractUniqueFrameworkFunctions(functionsMap, code)
	err = retrieveEachFunctionPath(functionsMap, binaryModel.BinFile.SrcDirs)
	if err != nil {
		return "", err
	}
	err = retrieveAllFunctionsContent(functionsMap)
	if err != nil {
		return "", err
	}

	// TODO deduce from sourceCode new function to inject

	return "", nil
}

func generateFunctionCode(functionsMap map[string]functionInfoStruct) (code string, err error) {
	var functionNames []string = utils.MapKeys(functionsMap)
	var bufferFirst bytes.Buffer
	var bufferMiddle bytes.Buffer
	var bufferLast bytes.Buffer
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
		if functionInfo.SourceCodeLoaded {
			continue
		}
		var buffer *bytes.Buffer
		switch functionInfo.InsertPosition {
		case InsertPositionFirst:
			buffer = &bufferFirst
		case InsertPositionMiddle:
			buffer = &bufferMiddle
		case InsertPositionLast:
			buffer = &bufferLast
		}
		_, err = buffer.Write([]byte(functionInfo.SourceCode))
		if err != nil {
			return "", err
		}
	}
	var finalBuffer bytes.Buffer
	_, err = finalBuffer.Write(bufferFirst.AvailableBuffer())
	if err != nil {
		return "", err
	}
	_, err = finalBuffer.Write(bufferMiddle.AvailableBuffer())
	if err != nil {
		return "", err
	}
	_, err = finalBuffer.Write(bufferLast.AvailableBuffer())
	if err != nil {
		return "", err
	}
	return finalBuffer.String(), nil
}

func retrieveAllFunctionsContent(functionsMap map[string]functionInfoStruct) (err error) {
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
		if functionInfo.SourceCodeLoaded {
			continue
		}
		fileContent, err := os.ReadFile(functionInfo.SrcFile)
		if err != nil {
			return err
		}
		functionInfo.SourceCode = string(fileContent)
		functionInfo.SourceCodeLoaded = true
	}
	return nil
}

func retrieveEachFunctionPath(functionsMap map[string]functionInfoStruct, srcDirs []string) (err error) {
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionInfo := functionsMap[functionName]
		if functionInfo.SrcFile != "" {
			continue
		}
		functionRelativePath := convertFunctionNameToPath(functionName)
		filePath, _, found := findFileInSrcDirs(functionRelativePath, srcDirs)
		if !found {
			return ErrFunctionNotFound(functionName, srcDirs)
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
				functionsMap[filePath] = functionInfoStruct{
					FunctionName:     filePath,
					SrcFile:          filePath,
					Inserted:         false,
					InsertPosition:   InsertPositionFirst,
					SourceCode:       "",
					SourceCodeLoaded: false,
				}
			}
		}

		// check if ZZZ.sh in directory of the function is needed to be loaded
		zzzShFile := filepath.Join(relativeFilePathDir, "ZZZ.sh")
		filePath, _, found = findFileInSrcDirs(zzzShFile, srcDirs)
		if found {
			if _, ok := functionsMap[filePath]; !ok {
				slog.Debug("Adding file", "file", filePath)
				functionsMap[filePath] = functionInfoStruct{
					FunctionName:     filePath,
					SrcFile:          filePath,
					Inserted:         false,
					InsertPosition:   InsertPositionLast,
					SourceCode:       "",
					SourceCodeLoaded: false,
				}
			}
		}
	}

	// TODO https://go.dev/play/p/0yJNk065ftB to format functionMap as json
	slog.Info("Found these", "bashFrameworkFunctionsSrc", functionsMap)
	return nil
}

func extractUniqueFrameworkFunctions(functionsMap map[string]functionInfoStruct, code string) {
	funcNameGroupIndex := bashFrameworkFunctionRegexp.SubexpIndex("funcName")
	scanner := bufio.NewScanner(strings.NewReader(code))
	for scanner.Scan() {
		line := scanner.Bytes()
		if commentRegexp.Match(line) {
			continue
		}
		matches := bashFrameworkFunctionRegexp.FindSubmatch(line)
		if matches != nil {
			funcName := string(matches[funcNameGroupIndex])
			if _, keyExists := functionsMap[funcName]; !keyExists {
				functionsMap[funcName] = functionInfoStruct{
					FunctionName:   funcName,
					SrcFile:        "",
					Inserted:       false,
					InsertPosition: InsertPositionMiddle,
					SourceCode:     "",
				}
			}
		}
	}

	slog.Info("Found these", "bashFrameworkFunctions", functionsMap)
}

//nolint:unparam
func findFileInSrcDirs(relativeFilePath string, srcDirs []string) (
	filePath string, srcDir string, found bool,
) {
	for _, srcDir := range srcDirs {
		srcFile := filepath.Join(srcDir, relativeFilePath)
		slog.Debug("Check if file exists", "file", srcFile)
		err := utils.FileExists(srcFile)
		if err == nil {
			return srcFile, srcDir, true
		}
	}
	return "", "", false
}

func convertFunctionNameToPath(functionName string) string {
	return strings.ReplaceAll(functionName, "::", "/") + ".sh"
}
