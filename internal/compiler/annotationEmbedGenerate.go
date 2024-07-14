package compiler

import (
	"fmt"
	"io"
	"os"

	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/tar"
	"github.com/fchastanet/bash-compiler/internal/utils"
)

func (annotationProcessor *embedAnnotationProcessor) RenderResource(
	asName string,
	resource string,
	lineNumber int,
) (string, error) {
	fi, err := os.Stat(resource)
	if err == nil {
		switch mode := fi.Mode(); {
		case mode.IsDir():
			return annotationProcessor.renderDir(asName, resource)
		case mode.IsRegular():
			return annotationProcessor.renderFile(asName, resource, mode)
		}
	}

	return "", &unsupportedEmbeddedResourceError{
		asName, resource, lineNumber, err,
	}
}

func (annotationProcessor *embedAnnotationProcessor) renderFile(
	asName string,
	resource string,
	fileMode os.FileMode,
) (string, error) {
	file, err := os.Open(resource)
	if logger.FancyHandleError(err) {
		return "", err
	}
	defer file.Close()

	md5sum, err := utils.Md5SumFromFile(file)
	if logger.FancyHandleError(err) {
		return "", err
	}
	base64, err := utils.Base64FromFile(file)
	if logger.FancyHandleError(err) {
		return "", err
	}

	data := map[string]string{
		"asName":   asName,
		"fileMode": fmt.Sprintf("%o", fileMode.Perm()),
		"base64":   base64,
		"md5sum":   md5sum,
	}
	code, err := annotationProcessor.renderTemplate(
		data, annotationProcessor.embedFileTemplateName,
	)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return code, nil
}

func (annotationProcessor *embedAnnotationProcessor) renderDir(
	asName string,
	resource string,
) (string, error) {
	directoryArchive, err := os.CreateTemp("", "directoryArchive*.tgz")
	if logger.FancyHandleError(err) {
		return "", err
	}
	err = createDirectoryArchive(resource, directoryArchive)
	if logger.FancyHandleError(err) {
		return "", err
	}
	defer directoryArchive.Close()

	_, err = directoryArchive.Seek(0, 0)
	if err != nil {
		return "", err
	}
	md5sum, err := utils.Md5SumFromFile(directoryArchive)
	if err != nil {
		return "", err
	}
	_, err = directoryArchive.Seek(0, 0)
	if err != nil {
		return "", err
	}
	base64, err := utils.Base64FromFile(directoryArchive)
	if err != nil {
		return "", err
	}

	data := map[string]string{
		"asName": asName,
		"base64": base64,
		"md5sum": md5sum,
	}
	code, err := annotationProcessor.renderTemplate(
		data, annotationProcessor.embedDirTemplateName,
	)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return code, nil
}

func createDirectoryArchive(directory string, buf io.Writer) error {
	filesList, err := files.MatchFullDirectory(directory)
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

func (annotationProcessor *embedAnnotationProcessor) renderTemplate(
	data map[string]string,
	templateName string,
) (string, error) {
	annotationProcessor.context.templateContext.Data = data
	annotationProcessor.context.templateContext.RootData = data
	return annotationProcessor.context.templateContext.Render(
		templateName,
	)
}
