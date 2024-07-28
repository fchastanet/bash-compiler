package compiler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/fchastanet/bash-compiler/internal/utils/structures"
)

var (
	requireRegexp         = regexp.MustCompile(`# @require (?P<require>.*)$`)
	requiredFunctionRegex = regexp.MustCompile(
		`(?m)[ \t]*(function[ \t]+|)(?P<bashFrameworkFunction>([A-Za-z0-9_]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))\(\)[ \t]*\{[ \t]*$`,
	)
)

var errRequiredFunctionNotFound = errors.New("Required function not found")

func ErrRequiredFunctionNotFound(functionName string) error {
	return fmt.Errorf("%w: %s in parsed code", errRequiredFunctionNotFound, functionName)
}

const annotationRequireKind string = "require"

type requireAnnotationProcessor struct {
	compileContextData            *CompileContextData
	checkRequirementsTemplateName string
	requireTemplateName           string
}

type requireAnnotation struct {
	requiredFunctions            []string
	isRequired                   bool
	checkRequirementsCodeAdded   bool
	codeAddedOnRequiredFunctions bool
}

func NewRequireAnnotationProcessor() AnnotationProcessorInterface {
	return &requireAnnotationProcessor{}
}

func (annotationProcessor *requireAnnotationProcessor) Init(
	compileContextData *CompileContextData,
) error {
	annotationProcessor.compileContextData = compileContextData
	checkRequirementsTemplateName, err :=
		annotationProcessor.compileContextData.config.AnnotationsConfig.GetStringValue("checkRequirementsTemplate")
	if err != nil {
		return err
	}
	requireTemplateName, err :=
		annotationProcessor.compileContextData.config.AnnotationsConfig.GetStringValue("requireTemplate")
	if err != nil {
		return err
	}

	annotationProcessor.checkRequirementsTemplateName = checkRequirementsTemplateName
	annotationProcessor.requireTemplateName = requireTemplateName

	return nil
}

func (annotationProcessor *requireAnnotationProcessor) ParseFunction(
	compileContextData *CompileContextData,
	functionStruct *functionInfoStruct,
) error {
	annotation, err := functionStruct.getRequireAnnotation()
	if logger.FancyHandleError(err) {
		return err
	}
	annotation.requiredFunctions, functionStruct.SourceCode = extractRequiredFunctions(
		functionStruct.SourceCode,
	)

	if len(annotation.requiredFunctions) == 0 {
		return nil
	}
	err = isCodeContainsFunction(functionStruct.SourceCode, functionStruct.FunctionName)
	if err != nil {
		return err
	}

	functionStruct.SourceCode, err = myTemplateFunctions.MustInclude(
		annotationProcessor.checkRequirementsTemplateName,
		map[string]interface{}{
			"code":         functionStruct.SourceCode,
			"functionName": functionStruct.FunctionName,
			"requires":     annotation.requiredFunctions,
		},
		*compileContextData.templateContextData,
	)
	if err != nil {
		return err
	}
	annotation.checkRequirementsCodeAdded = true
	functionStruct.AnnotationMap[annotationRequireKind] = *annotation

	return nil
}

func extractRequiredFunctions(code string) ([]string, string) {
	var newCodeBuffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(code))
	requiredFunctions := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		matches := requireRegexp.FindStringSubmatch(line)
		if matches != nil {
			requireIndex := requireRegexp.SubexpIndex(annotationRequireKind)
			requiredFunctions = append(requiredFunctions, strings.Trim(matches[requireIndex], " \t"))
		} else {
			newCodeBuffer.Write([]byte(line))
			newCodeBuffer.WriteByte('\n')
		}
	}
	return requiredFunctions, newCodeBuffer.String()
}

func (annotationProcessor *requireAnnotationProcessor) Process(
	compileContextData *CompileContextData,
) error {
	functionsMap := compileContextData.functionsMap
	var functionNames []string = structures.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionStruct := functionsMap[functionName]
		slog.Debug("addRequireCodeToEachRequiredFunctions", "functionName", functionName)
		err := annotationProcessor.addRequireCodeToEachRequiredFunctions(compileContextData, &functionStruct)
		if err != nil {
			return err
		}
		compileContextData.functionsMap[functionName] = functionStruct
	}
	return nil
}

func (functionStruct *functionInfoStruct) getRequireAnnotation() (*requireAnnotation, error) {
	annotation, ok := functionStruct.AnnotationMap[annotationRequireKind]
	if !ok {
		annotation = requireAnnotation{
			requiredFunctions:            []string{},
			isRequired:                   false,
			checkRequirementsCodeAdded:   false,
			codeAddedOnRequiredFunctions: false,
		}
		functionStruct.AnnotationMap[annotationRequireKind] = annotation
	}
	requireAnnotation, ok := annotation.(requireAnnotation)
	if !ok {
		return nil, errAnnotationCastIssue
	}
	return &requireAnnotation, nil
}

func (annotationProcessor *requireAnnotationProcessor) addRequireCodeToEachRequiredFunctions(
	compileContextData *CompileContextData,
	functionStruct *functionInfoStruct,
) error {
	requireAnnotation, err := functionStruct.getRequireAnnotation()
	if err != nil {
		return err
	}

	if len(requireAnnotation.requiredFunctions) > 0 {
		functionsMap := compileContextData.functionsMap
		for _, requiredFunctionName := range requireAnnotation.requiredFunctions {
			slog.Debug("Check if required function has been imported", "requiredFunctionName", requiredFunctionName)
			requiredFunctionStruct, ok := functionsMap[requiredFunctionName]
			if !ok {
				return ErrRequiredFunctionNotFound(requiredFunctionName)
			}
			err = annotationProcessor.addRequireCode(compileContextData, &requiredFunctionStruct)
			if err != nil {
				return err
			}
			compileContextData.functionsMap[requiredFunctionName] = requiredFunctionStruct
		}
		requireAnnotation.codeAddedOnRequiredFunctions = true
	}
	functionStruct.AnnotationMap[annotationRequireKind] = *requireAnnotation
	return nil
}

func (annotationProcessor *requireAnnotationProcessor) addRequireCode(
	compileContextData *CompileContextData,
	functionStruct *functionInfoStruct,
) error {
	myRequiredAnnotation, err := functionStruct.getRequireAnnotation()
	if err != nil {
		return err
	}
	if myRequiredAnnotation.codeAddedOnRequiredFunctions {
		return nil
	}

	err = isCodeContainsFunction(functionStruct.SourceCode, functionStruct.FunctionName)
	if err != nil {
		return err
	}
	myRequiredAnnotation.isRequired = true

	sourceCode, err := myTemplateFunctions.MustInclude(
		annotationProcessor.requireTemplateName,
		map[string]interface{}{
			"code":         functionStruct.SourceCode,
			"functionName": functionStruct.FunctionName,
		},
		*compileContextData.templateContextData,
	)
	if err != nil {
		return err
	}
	functionStruct.SourceCode = sourceCode
	myRequiredAnnotation.codeAddedOnRequiredFunctions = true
	functionStruct.AnnotationMap[annotationRequireKind] = *myRequiredAnnotation
	return nil
}

func isCodeContainsFunction(code string, functionName string) error {
	matches := requiredFunctionRegex.FindAllStringSubmatch(code, -1)
	slog.Debug("isCodeContainsFunction", "functionName", functionName)
	if matches == nil {
		slog.Error("isCodeContainsFunction no function regexp match")
		return ErrRequiredFunctionNotFound(functionName)
	}
	bashFrameworkFunctionGroupIndex := requiredFunctionRegex.SubexpIndex("bashFrameworkFunction")
	for _, match := range matches {
		if match[bashFrameworkFunctionGroupIndex] == functionName {
			return nil
		}
	}
	slog.Error("isCodeContainsFunction function does not match", "functionName", functionName)
	return ErrRequiredFunctionNotFound(functionName)
}

func (annotationProcessor *requireAnnotationProcessor) PostProcess(
	_ *CompileContextData, code string,
) (string, error) {
	return code, nil
}
