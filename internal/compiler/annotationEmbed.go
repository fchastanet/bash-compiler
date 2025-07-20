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

type writeError struct {
	error
	lineNumber int
	message    string
}

func (e *writeError) Error() string {
	return fmt.Sprintf(
		"Failed to write full line from src line %d : %s",
		e.lineNumber, e.message,
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

func (*embedAnnotationProcessor) GetTitle() string {
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

func (*embedAnnotationProcessor) ParseFunction(
	_ *CompileContextData,
	_ *functionInfoStruct,
) error {
	return nil
}

func (*embedAnnotationProcessor) Process(_ *CompileContextData) error {
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
			err := annotationProcessor.processEmbedMatch(
				&bufferOutput,
				matches,
				embedRegexpResourceGroupIndex,
				embedRegexpAsNameGroupIndex,
				lineNumber,
			)
			if err != nil {
				return "", err
			}
		} else {
			err := annotationProcessor.processRegularLine(&bufferOutput, line, lineNumber)
			if err != nil {
				return "", err
			}
		}
		bufferOutput.WriteByte('\n')
	}

	return bufferOutput.String(), nil
}

func (annotationProcessor *embedAnnotationProcessor) processEmbedMatch(
	bufferOutput *bytes.Buffer,
	matches []string,
	resourceIndex, asNameIndex, lineNumber int,
) error {
	embedCode, err := annotationProcessor.generateEmbedCode(
		matches[resourceIndex],
		matches[asNameIndex],
		lineNumber,
	)
	if logger.FancyHandleError(err) {
		return err
	}
	return annotationProcessor.writeToBuffer(bufferOutput, embedCode, lineNumber, "annotationEmbed - failed to write full embed code")
}

func (annotationProcessor *embedAnnotationProcessor) processRegularLine(
	bufferOutput *bytes.Buffer,
	line string,
	lineNumber int,
) error {
	return annotationProcessor.writeToBuffer(bufferOutput, line, lineNumber, "annotationEmbed - failed to write full line")
}

func (*embedAnnotationProcessor) writeToBuffer(
	bufferOutput *bytes.Buffer,
	content string,
	lineNumber int,
	errorMessage string,
) error {
	n, err := bufferOutput.WriteString(content)
	if logger.FancyHandleError(err) {
		return err
	}
	if n != len(content) {
		return &writeError{
			error:      nil,
			lineNumber: lineNumber,
			message:    errorMessage,
		}
	}
	return nil
}

func (annotationProcessor *embedAnnotationProcessor) generateEmbedCode(
	resource string,
	asName string,
	lineNumber int,
) (string, error) {
	resource = os.ExpandEnv(strings.Trim(resource, " \t"))
	asName = strings.Trim(asName, " \t")
	if _, exists := annotationProcessor.embedMap[asName]; exists {
		return "", &duplicatedAsNameError{nil, lineNumber, asName, resource}
	}
	annotationProcessor.embedMap[asName] = resource
	return annotationProcessor.annotationEmbedGenerate.RenderResource(
		asName, resource, lineNumber,
	)
}
