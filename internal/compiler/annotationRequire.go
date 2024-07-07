package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	myTemplateFunctions "github.com/fchastanet/bash-compiler/internal/render/functions"
	"github.com/fchastanet/bash-compiler/internal/utils"
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
	context                       *compileContext
	kind                          string
	checkRequirementsTemplateName string
	requireTemplateName           string
}

type requireAnnotation struct {
	kind              string
	requiredFunctions []string
	isRequired        bool
	isComputed        bool
}

func NewRequireAnnotationProcessor(context *compileContext) AnnotationProcessorInterface {
	return &requireAnnotationProcessor{
		context: context,
		kind:    annotationRequireKind,
	}
}

func (annotationProcessor *requireAnnotationProcessor) Init() error {
	checkRequirementsTemplateName, err :=
		annotationProcessor.context.config.DynamicConfig.GetStringValue("checkRequirementsTemplate")
	if err != nil {
		return err
	}
	requireTemplateName, err :=
		annotationProcessor.context.config.DynamicConfig.GetStringValue("requireTemplate")
	if err != nil {
		return err
	}

	annotationProcessor.checkRequirementsTemplateName = checkRequirementsTemplateName
	annotationProcessor.requireTemplateName = requireTemplateName

	return nil
}

func (annotationProcessor *requireAnnotationProcessor) ParseFunction(functionStruct *functionInfoStruct) error {
	annotation, ok := functionStruct.AnnotationMap[annotationRequireKind]
	if !ok {
		annotation = requireAnnotation{
			kind:              annotationRequireKind,
			requiredFunctions: []string{},
			isRequired:        false,
			isComputed:        false,
		}
	}
	myAnnotation, ok := annotation.(requireAnnotation)
	if !ok {
		return errAnnotationCastIssue
	}
	requiredFunctions := findRequiredFunctions(
		functionStruct.SourceCode,
	)
	myAnnotation.requiredFunctions = requiredFunctions

	functionStruct.AnnotationMap[annotationRequireKind] = myAnnotation

	return nil
}

func findRequiredFunctions(code string) []string {
	scanner := bufio.NewScanner(strings.NewReader(code))
	requiredFunctions := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		matches := requireRegexp.FindStringSubmatch(line)
		if matches != nil {
			requireIndex := requireRegexp.SubexpIndex(annotationRequireKind)
			requiredFunctions = append(requiredFunctions, strings.Trim(matches[requireIndex], " \t"))
		}
	}
	return requiredFunctions
}

func (annotationProcessor *requireAnnotationProcessor) Process(_ string) error {
	functionsMap := annotationProcessor.context.functionsMap
	var functionNames []string = utils.MapKeys(functionsMap)
	for _, functionName := range functionNames {
		functionStruct := functionsMap[functionName]
		slog.Debug("addCheckRequirementsCodeIfNeeded", "functionName", functionName)
		err := annotationProcessor.addCheckRequirementsCodeIfNeeded(&functionStruct)
		if err != nil {
			return err
		}
		slog.Debug("addRequireCodeToEachRequiredFunctions", "functionName", functionName)
		err = annotationProcessor.addRequireCodeToEachRequiredFunctions(&functionStruct)
		if err != nil {
			return err
		}
		annotationProcessor.context.functionsMap[functionName] = functionStruct
	}
	return nil
}

func (functionStruct *functionInfoStruct) getRequireAnnotation() (*requireAnnotation, error) {
	annotation, ok := functionStruct.AnnotationMap[annotationRequireKind]
	if !ok {
		// shouldn't happen because already computed during ParseFunction
		return nil, nil
	}
	requireAnnotation, ok := annotation.(requireAnnotation)
	if !ok {
		return nil, errAnnotationCastIssue
	}
	return &requireAnnotation, nil
}

func (annotationProcessor *requireAnnotationProcessor) addRequireCodeToEachRequiredFunctions(
	functionStruct *functionInfoStruct,
) error {
	requireAnnotation, err := functionStruct.getRequireAnnotation()
	if err != nil {
		return err
	}

	if len(requireAnnotation.requiredFunctions) > 0 {
		functionsMap := annotationProcessor.context.functionsMap
		for _, requiredFunctionName := range requireAnnotation.requiredFunctions {
			slog.Debug("Check if required function has been imported", "requiredFunctionName", requiredFunctionName)
			requiredFunctionStruct, ok := functionsMap[requiredFunctionName]
			if !ok {
				return ErrRequiredFunctionNotFound(requiredFunctionName)
			}
			err = annotationProcessor.addRequireCode(&requiredFunctionStruct)
			if err != nil {
				return err
			}
			annotationProcessor.context.functionsMap[requiredFunctionName] = requiredFunctionStruct
		}
	}
	return nil
}

func (annotationProcessor *requireAnnotationProcessor) addCheckRequirementsCodeIfNeeded(
	functionStruct *functionInfoStruct,
) error {
	requireAnnotation, err := functionStruct.getRequireAnnotation()
	if err != nil {
		return err
	}
	if len(requireAnnotation.requiredFunctions) > 0 {
		slog.Debug("addCheckRequirementsCode",
			"functionName", functionStruct.FunctionName, "requiredFunctions", requireAnnotation.requiredFunctions)
		err = annotationProcessor.addCheckRequirementsCode(
			functionStruct,
			requireAnnotation.requiredFunctions,
		)
	}
	return err
}

func (annotationProcessor *requireAnnotationProcessor) addRequireCode(
	functionStruct *functionInfoStruct,
) error {
	myRequiredAnnotation, err := functionStruct.getRequireAnnotation()
	if err != nil {
		return err
	}
	if myRequiredAnnotation.isComputed {
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
		*annotationProcessor.context.templateContext,
	)
	if err != nil {
		return err
	}
	functionStruct.SourceCode = sourceCode
	myRequiredAnnotation.isComputed = true
	return nil
}

func (annotationProcessor *requireAnnotationProcessor) addCheckRequirementsCode(
	functionStruct *functionInfoStruct,
	requires []string,
) error {
	err := isCodeContainsFunction(functionStruct.SourceCode, functionStruct.FunctionName)
	if err != nil {
		return err
	}

	functionStruct.SourceCode, err = myTemplateFunctions.MustInclude(
		annotationProcessor.checkRequirementsTemplateName,
		map[string]interface{}{
			"code":         functionStruct.SourceCode,
			"functionName": functionStruct.FunctionName,
			"requires":     requires,
		},
		*annotationProcessor.context.templateContext,
	)
	return err
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
