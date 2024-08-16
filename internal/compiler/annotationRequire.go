package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"

	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

var (
	requireRegexp         = regexp.MustCompile(`# @require (?P<require>.*)$`)
	requiredFunctionRegex = regexp.MustCompile(
		`(?m)[ \t]*(function[ \t]+|)(?P<bashFrameworkFunction>([A-Za-z0-9_]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))\(\)[ \t]*\{[ \t]*$`,
	)
)

type requiredFunctionNotFoundError struct {
	error
	functionName string
}

func (e *requiredFunctionNotFoundError) Error() string {
	msg := "required function not found in parsed code: " + e.functionName
	if e.error != nil {
		msg = fmt.Sprintf("%s - inner error:\n%v", msg, e.error)
	}
	return msg
}

const annotationRequireKind string = "require"

type requireAnnotationProcessor struct {
	annotationProcessor
	compileContextData            *CompileContextData
	checkRequirementsTemplateName string
	requireTemplateName           string
}

type requireAnnotation struct {
	annotation
	requiredFunctions            []string
	isRequired                   bool
	checkRequirementsCodeAdded   bool
	codeAddedOnRequiredFunctions bool
}

func NewRequireAnnotationProcessor() AnnotationProcessorInterface {
	return &requireAnnotationProcessor{} //nolint:exhaustruct // Check Init method
}

func (annotationProcessor *requireAnnotationProcessor) Init(
	compileContextData *CompileContextData,
) error {
	if compileContextData == nil {
		return validationError("compileContextData", nil)
	}
	err := compileContextData.Validate()
	if logger.FancyHandleError(err) {
		return err
	}
	annotationProcessor.compileContextData = compileContextData
	checkRequirementsTemplateName, err := annotationProcessor.compileContextData.config.AnnotationsConfig.GetStringValue("checkRequirementsTemplateName")
	if err != nil {
		return &errors.ValidationError{
			InnerError: err,
			Context:    "compileContextData.config.AnnotationsConfig",
			FieldName:  "checkRequirementsTemplateName",
			FieldValue: nil,
		}
	}
	requireTemplateName, err := annotationProcessor.compileContextData.config.AnnotationsConfig.GetStringValue("requireTemplateName")
	if err != nil {
		return &errors.ValidationError{
			InnerError: err,
			Context:    "compileContextData.config.AnnotationsConfig",
			FieldName:  "requireTemplateName",
			FieldValue: nil,
		}
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

func extractRequiredFunctions(code string) (requiredFunctions []string, newCode string) {
	var newCodeBuffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(code))
	requiredFunctions = []string{}
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
	newCode = newCodeBuffer.String()
	return requiredFunctions, newCode
}

func (annotationProcessor *requireAnnotationProcessor) Process(
	compileContextData *CompileContextData,
) error {
	functionsMap := compileContextData.functionsMap
	functionNames := getSortedFunctionNamesFromMap(functionsMap)
	sort.Strings(functionNames)
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
	annotationObj, ok := functionStruct.AnnotationMap[annotationRequireKind]
	if annotationObj == nil || !ok {
		newAnnotation := requireAnnotation{
			annotation:                   annotation{},
			requiredFunctions:            []string{},
			isRequired:                   false,
			checkRequirementsCodeAdded:   false,
			codeAddedOnRequiredFunctions: false,
		}
		functionStruct.AnnotationMap[annotationRequireKind] = annotationObj
		return &newAnnotation, nil
	}
	castedAnnotation, ok := annotationObj.(requireAnnotation) //nolint:gocritic
	if !ok {
		return nil, &annotationCastError{nil, functionStruct.FunctionName}
	}
	return &castedAnnotation, nil
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
				return &requiredFunctionNotFoundError{nil, requiredFunctionName}
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
		return &requiredFunctionNotFoundError{nil, functionName}
	}
	bashFrameworkFunctionGroupIndex := requiredFunctionRegex.SubexpIndex("bashFrameworkFunction")
	for _, match := range matches {
		if match[bashFrameworkFunctionGroupIndex] == functionName {
			return nil
		}
	}
	slog.Error("isCodeContainsFunction function does not match", "functionName", functionName)
	return &requiredFunctionNotFoundError{nil, functionName}
}

func (annotationProcessor *requireAnnotationProcessor) PostProcess(
	_ *CompileContextData, code string,
) (string, error) {
	return code, nil
}
