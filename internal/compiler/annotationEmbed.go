package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

var embedRegexp = regexp.MustCompile(
	`(?m)# @embed[ \t]+["']?(?P<resource>[^ \t"']+)["']?[ \t]+(AS|as|As)[ \t]+(?P<asName>[^ \t]+)$`,
)

type duplicatedAsNameError struct {
	error
	lineNumber int
	asName     string
	resource   string
}

func (e *duplicatedAsNameError) Error() string {
	return fmt.Sprintf(
		"Embedded resource '%s' - name '%s' is duplicated on line %d",
		e.resource, e.asName, e.lineNumber,
	)
}

type embedAnnotationProcessor struct {
	annotationProcessor
	annotationEmbedGenerate annotationEmbedGenerateInterface
	embedMap                map[string]string
}

func NewEmbedAnnotationProcessor() AnnotationProcessorInterface {
	return &embedAnnotationProcessor{} //nolint:exhaustruct // Check Init method
}

func validationError(fieldName string, fieldValue any) error {
	return &errors.ValidationError{
		InnerError: nil,
		Context:    "annotationEmbed",
		FieldName:  fieldName,
		FieldValue: fieldValue,
	}
}

func (annotationProcessor *embedAnnotationProcessor) GetTitle() string {
	return "EmbedAnnotationProcessor"
}

func (annotationProcessor *embedAnnotationProcessor) Init(
	compileContextData *CompileContextData,
) error {
	if compileContextData == nil {
		return validationError("compileContextData", nil)
	}
	err := compileContextData.Validate()
	if logger.FancyHandleError(err) {
		return err
	}
	annotationProcessor.embedMap = make(map[string]string)

	embedFileTemplateName, err := compileContextData.config.AnnotationsConfig.GetStringValue("embedFileTemplateName")
	if logger.FancyHandleError(err) {
		return &errors.ValidationError{
			InnerError: err,
			Context:    "compileContextData.config.AnnotationsConfig",
			FieldName:  "embedFileTemplateName",
			FieldValue: nil,
		}
	}

	embedDirTemplateName, err := compileContextData.config.AnnotationsConfig.GetStringValue("embedDirTemplateName")
	if logger.FancyHandleError(err) {
		return &errors.ValidationError{
			InnerError: err,
			Context:    "compileContextData.config.AnnotationsConfig",
			FieldName:  "embedDirTemplateName",
			FieldValue: nil,
		}
	}
	annotationProcessor.annotationEmbedGenerate = &annotationEmbedGenerate{
		embedDirTemplateName:  embedDirTemplateName,
		embedFileTemplateName: embedFileTemplateName,
		templateContextData:   compileContextData.templateContextData,
	}

	return nil
}

func (annotationProcessor *embedAnnotationProcessor) Reset() {
	annotationProcessor.embedMap = make(map[string]string)
}

func (annotationProcessor *embedAnnotationProcessor) ParseFunction(
	_ *CompileContextData,
	_ *functionInfoStruct,
) error {
	return nil
}

func (annotationProcessor *embedAnnotationProcessor) Process(
	_ *CompileContextData,
) error {
	return nil
}

func (annotationProcessor *embedAnnotationProcessor) PostProcess(
	_ *CompileContextData,
	code string,
) (string, error) {
	var bufferOutput bytes.Buffer
	embedRegexpResourceGroupIndex := embedRegexp.SubexpIndex("resource")
	embedRegexpAsNameGroupIndex := embedRegexp.SubexpIndex("asName")
	scanner := bufio.NewScanner(strings.NewReader(code))
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		matches := embedRegexp.FindStringSubmatch(line)
		if matches != nil {
			resource := os.ExpandEnv(strings.Trim(matches[embedRegexpResourceGroupIndex], " \t"))
			asName := strings.Trim(matches[embedRegexpAsNameGroupIndex], " \t")
			if _, exists := annotationProcessor.embedMap[asName]; exists {
				return "", &duplicatedAsNameError{nil, lineNumber, asName, resource}
			}
			annotationProcessor.embedMap[asName] = resource
			embedCode, err := annotationProcessor.annotationEmbedGenerate.RenderResource(
				asName, resource, lineNumber,
			)

			if logger.FancyHandleError(err) {
				return "", err
			}
			bufferOutput.Write([]byte(embedCode))
		} else {
			bufferOutput.Write([]byte(line))
		}
		bufferOutput.WriteByte('\n')
	}

	return bufferOutput.String(), nil
}
