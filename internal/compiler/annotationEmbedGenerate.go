package compiler

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/render"
	"github.com/fchastanet/bash-compiler/internal/utils/encoding"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
	"github.com/fchastanet/bash-compiler/internal/utils/tar"
)

type annotationEmbedGenerateInterface interface {
	RenderResource(asName string, resource string, lineNumber int) (string, error)
}

type annotationEmbedGenerate struct {
	embedDirTemplateName  string
	embedFileTemplateName string
	templateContextData   *render.TemplateContextData
}

type unsupportedEmbeddedResourceError struct {
	error
	asName     string
	resource   string
	lineNumber int
}

func (e *unsupportedEmbeddedResourceError) Error() string {
	msg := fmt.Sprintf(
		"Embedded resource '%s' - name '%s' on line %d cannot be embedded",
		e.resource, e.asName, e.lineNumber,
	)
	if e.error != nil {
		msg = fmt.Sprintf("%s - inner error: %v", msg, e.error)
	}
	return msg
}

func (annotationEmbedGenerate *annotationEmbedGenerate) RenderResource(
	asName string,
	resource string,
	lineNumber int,
) (string, error) {
	fi, err := os.Stat(resource)
	if err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			return annotationEmbedGenerate.renderDir(asName, resource)
		case mode.IsRegular():
			return annotationEmbedGenerate.renderFile(asName, resource, mode)
		}
	}

	return "", &unsupportedEmbeddedResourceError{
		err, asName, resource, lineNumber,
	}
}

func (annotationEmbedGenerate *annotationEmbedGenerate) renderFile(
	asName string,
	resource string,
	fileMode os.FileMode,
) (string, error) {
	file, err := os.Open(resource)
	if logger.FancyHandleError(err) {
		return "", err
	}
	defer file.Close()

	md5sum, err := encoding.ChecksumFromFile(file)
	if logger.FancyHandleError(err) {
		return "", err
	}
	file.Seek(0, 0)
	base64, err := encoding.Base64FromFile(file)
	if logger.FancyHandleError(err) {
		return "", err
	}

	data := map[string]string{
		"asName":   asName,
		"fileMode": fmt.Sprintf("%o", fileMode.Perm()),
		"base64":   base64,
		"md5sum":   md5sum,
	}
	code, err := annotationEmbedGenerate.renderTemplate(
		data, annotationEmbedGenerate.embedFileTemplateName,
	)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return code, nil
}

func (annotationEmbedGenerate *annotationEmbedGenerate) renderDir(
	asName string,
	resource string,
) (string, error) {
	directoryArchive, err := os.CreateTemp("", "directoryArchive*.tgz")
	if logger.FancyHandleError(err) {
		return "", err
	}
	slog.Info("Create directory archive", "sourceDir", resource, "targetFile", directoryArchive)
	err = createDirectoryArchive(resource, directoryArchive)
	if logger.FancyHandleError(err) {
		return "", err
	}
	defer directoryArchive.Close()

	_, err = directoryArchive.Seek(0, 0)
	if err != nil {
		return "", err
	}
	md5sum, err := encoding.ChecksumFromFile(directoryArchive)
	if err != nil {
		return "", err
	}
	_, err = directoryArchive.Seek(0, 0)
	if err != nil {
		return "", err
	}
	base64, err := encoding.Base64FromFile(directoryArchive)
	if err != nil {
		return "", err
	}

	data := map[string]string{
		"asName": asName,
		"base64": base64,
		"md5sum": md5sum,
	}
	code, err := annotationEmbedGenerate.renderTemplate(
		data, annotationEmbedGenerate.embedDirTemplateName,
	)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return code, nil
}

func createDirectoryArchive(directory string, buf io.Writer) error {
	filesList, err := files.MatchPatterns(directory, "**/*")
	if logger.FancyHandleError(err) {
		return err
	}
	files.SortFilesByPath(filesList)

	err = tar.CreateArchive(
		filesList,
		directory,
		buf,
		tar.ReproducibleTarOptions,
	)
	if err != nil {
		return err
	}

	return nil
}

func (annotationEmbedGenerate *annotationEmbedGenerate) renderTemplate(
	data map[string]string,
	templateName string,
) (string, error) {
	annotationEmbedGenerate.templateContextData.Data = data
	annotationEmbedGenerate.templateContextData.RootData = data
	return annotationEmbedGenerate.templateContextData.TemplateContext.Render(
		annotationEmbedGenerate.templateContextData,
		templateName,
	)
}
