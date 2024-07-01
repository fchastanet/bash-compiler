package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/fchastanet/bash-compiler/internal/render"
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
const annotationRequireRequires string = "requires"
const annotationRequireRequired string = "required"
const annotationRequireComputed string = "computed"

func NewRequireAnnotationProcessor() AnnotationProcessorInterface {
	return &Annotation{
		kind: annotationRequireKind,
	}
}

func (annotation *Annotation) Process(compileContext *compileContext, _ string) error {
	err := compileContext.computeRequires()
	if err != nil {
		return err
	}
	err = compileContext.computeRequired()
	if err != nil {
		return err
	}
	return nil
}

func (context *compileContext) computeRequires() (err error) {
	checkRequirementsTemplateName, err :=
		context.config.DynamicConfig.GetStringValue("checkRequirementsTemplate")
	if err != nil {
		return err
	}
	var functionNames []string = utils.MapKeys(context.functionsMap)
	for _, functionName := range functionNames {
		functionStruct := context.functionsMap[functionName]
		annotation, ok := functionStruct.AnnotationMap[annotationRequireKind]
		if !ok {
			annotation = Annotation{
				kind: annotationRequireKind,
				properties: map[string]interface{}{
					annotationRequireRequires: []string{},
					annotationRequireRequired: false,
					annotationRequireComputed: false,
				},
			}
			functionStruct.AnnotationMap = map[string]Annotation{
				annotationRequireKind: annotation,
			}
		}
		requiredFunctions := findRequiredFunctions(
			functionStruct.SourceCode,
		)
		annotation.properties[annotationRequireRequires] = requiredFunctions
		if len(requiredFunctions) > 0 {
			functionStruct.SourceCode, err = addCheckRequirementsCode(
				functionStruct.SourceCode,
				functionName,
				requiredFunctions,
				checkRequirementsTemplateName,
				*context.templateContext,
			)
			if err != nil {
				return err
			}
		}

		context.functionsMap[functionName] = functionStruct
	}
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

func (context *compileContext) computeRequired() (err error) {
	var functionNames []string = utils.MapKeys(context.functionsMap)
	requireTemplateName, err := context.config.DynamicConfig.GetStringValue("requireTemplate")
	if err != nil {
		return err
	}
	for _, functionName := range functionNames {
		functionStruct := context.functionsMap[functionName]
		annotation, ok := functionStruct.AnnotationMap[annotationRequireKind]
		if !ok {
			// shouldn't happen
			continue
		}
		requires := annotation.properties[annotationRequireRequires].([]string)
		for _, requiredFunctionName := range requires {
			requiredFunctionStruct, ok := context.functionsMap[requiredFunctionName]
			if !ok {
				return ErrRequiredFunctionNotFound(requiredFunctionName)
			}
			computedValue, ok := requiredFunctionStruct.AnnotationMap[annotationRequireKind].
				properties[annotationRequireComputed]
			computed := computedValue.(bool)
			if ok && computed {
				continue
			}
			requiredFunctionStruct.AnnotationMap[annotationRequireKind].
				properties[annotationRequireRequired] = true
			requiredFunctionStruct.SourceCode, err = addRequireCode(
				requiredFunctionStruct.SourceCode,
				requiredFunctionName,
				requireTemplateName,
				*context.templateContext,
			)
			requiredFunctionStruct.AnnotationMap[annotationRequireKind].
				properties[annotationRequireComputed] = true
			if err != nil {
				return err
			}
			context.functionsMap[requiredFunctionName] = requiredFunctionStruct
		}
	}
	return nil
}

func addRequireCode(
	code string,
	functionName string,
	requireTemplateName string,
	templateContext render.Context,
) (string, error) {
	err := isCodeContainsFunction(code, functionName)
	if err != nil {
		return "", err
	}

	return myTemplateFunctions.MustInclude(
		requireTemplateName,
		map[string]interface{}{
			"functionName": functionName,
			"code":         code,
		},
		templateContext,
	)
}

func addCheckRequirementsCode(
	code string,
	functionName string,
	requires []string,
	checkRequirementsTemplateName string,
	templateContext render.Context,
) (string, error) {
	err := isCodeContainsFunction(code, functionName)
	if err != nil {
		return "", err
	}

	return myTemplateFunctions.MustInclude(
		checkRequirementsTemplateName,
		map[string]interface{}{
			"code":         code,
			"functionName": functionName,
			"requires":     requires,
		},
		templateContext,
	)
}

func isCodeContainsFunction(code string, functionName string) error {
	matches := requiredFunctionRegex.FindAllStringSubmatch(code, -1)
	if matches == nil {
		return ErrRequiredFunctionNotFound(functionName)
	}
	bashFrameworkFunctionGroupIndex := requiredFunctionRegex.SubexpIndex("bashFrameworkFunction")
	for _, match := range matches {
		if match[bashFrameworkFunctionGroupIndex] == functionName {
			return nil
		}
	}
	return ErrRequiredFunctionNotFound(functionName)
}
