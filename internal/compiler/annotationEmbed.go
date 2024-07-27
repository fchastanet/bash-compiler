package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

var (
	embedRegexp = regexp.MustCompile(
		`(?m)# @embed[ \t]+["']?(?P<resource>[^ \t"']+)["']?[ \t]+(AS|as|As)[ \t]+(?P<asName>[^ \t]+)$`,
	)
)

type unsupportedEmbeddedResourceError struct {
	asName     string
	resource   string
	lineNumber int
	innerError error
}

func (e *unsupportedEmbeddedResourceError) Error() string {
	msg := fmt.Sprintf(
		"Embedded resource '%s' - name '%s' on line %d cannot be embedded",
		e.resource, e.asName, e.lineNumber,
	)
	if e.innerError != nil {
		msg = fmt.Sprintf("%s - inner error:\n%v", msg, e.innerError)
	}
	return msg
}

type duplicatedAsNameError struct {
	lineNumber int
	asName     string
	resource   string
}

func (e *duplicatedAsNameError) Error() string {
	return fmt.Sprintf(
		"Embedded resource '%s' - name '%s' is already used on line %d",
		e.resource, e.asName, e.lineNumber,
	)
}

type embedAnnotationProcessor struct {
	context               *CompileContext
	embedFileTemplateName string
	embedDirTemplateName  string
	embedMap              map[string]string
}

func NewEmbedAnnotationProcessor(context *CompileContext) AnnotationProcessorInterface {
	return &embedAnnotationProcessor{
		context:  context,
		embedMap: nil,
	}
}

func (annotationProcessor *embedAnnotationProcessor) Init() error {
	annotationProcessor.embedMap = make(map[string]string)

	embedFileTemplateName, err :=
		annotationProcessor.context.config.AnnotationsConfig.GetStringValue("embedFileTemplateName")
	if logger.FancyHandleError(err) {
		return err
	}
	annotationProcessor.embedFileTemplateName = embedFileTemplateName

	embedDirTemplateName, err :=
		annotationProcessor.context.config.AnnotationsConfig.GetStringValue("embedDirTemplateName")
	if logger.FancyHandleError(err) {
		return err
	}
	annotationProcessor.embedDirTemplateName = embedDirTemplateName

	return nil
}

func (annotationProcessor *embedAnnotationProcessor) ParseFunction(_ *functionInfoStruct) error {
	return nil
}

func (annotationProcessor *embedAnnotationProcessor) Process() error {
	return nil
}

func (annotationProcessor *embedAnnotationProcessor) PostProcess(code string) (string, error) {
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
				return "", &duplicatedAsNameError{
					lineNumber: lineNumber,
					asName:     asName,
					resource:   resource,
				}
			}
			annotationProcessor.embedMap[asName] = resource
			embedCode, err := annotationProcessor.RenderResource(
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
