// Package functions
package functions

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/fchastanet/bash-compiler/internal/code"
	"github.com/fchastanet/bash-compiler/internal/files"
	"github.com/fchastanet/bash-compiler/internal/logger"
	"github.com/fchastanet/bash-compiler/internal/render"
)

var errFileNotFound = errors.New("File does not exist")

func ErrFileNotFound(file string, srcDirs []string) error {
	return fmt.Errorf("%w: %s in any srcDirs %v", errFileNotFound, file, srcDirs)
}

// include allows to include a template
// allowing to use filter
// Eg: {{ include "template.tpl" | indent 4 }}
func Include(
	template string, templateData any,
	templateContext render.Context,
) string {
	var output string
	output, _ = MustInclude(template, templateData, templateContext)
	return output
}

func MustInclude(
	templateName string,
	templateData any,
	templateContext render.Context,
) (output string, err error) {
	slog.Debug("MustInclude",
		logger.LogFieldTemplateName, templateName,
		logger.LogFieldTemplateData, templateData,
	)
	templateContext.Data = templateData
	output, err = templateContext.Render(templateName)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return code.RemoveFirstShebangLineIfAny(output), err
}

func includeFile(filePath string) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Debug(
		"includeFile",
		logger.LogFieldFilePath, filePath,
		logger.LogFieldFilePathExpanded, filePathExpanded,
	)

	file, err := os.ReadFile(filePathExpanded)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(file)
}

func RenderFromTemplateContent(
	templateContext *render.Context, templateContent string,
) (codeStr string, err error) {
	template, err := templateContext.Template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	var tplWriter bytes.Buffer
	err = template.Execute(&tplWriter, templateContext)
	if err != nil {
		return "", err
	}

	return code.RemoveFirstShebangLineIfAny(tplWriter.String()), err
}

func includeFileAsTemplate(filePath string, templateContext render.Context) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Info(
		"includeFileAsTemplate",
		logger.LogFieldFilePath, filePath,
		logger.LogFieldFilePathExpanded, filePathExpanded,
	)

	file, err := os.ReadFile(filePathExpanded)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	code, err := RenderFromTemplateContent(&templateContext, string(file))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return code
}

func dynamicFile(filePath string, paths []string) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Info(
		"dynamicFile",
		logger.LogFieldFilePath, filePath,
		logger.LogFieldFilePathExpanded, filePathExpanded,
	)
	err := files.FileExists(filePathExpanded)
	if err == nil {
		return filePathExpanded
	}
	for _, dir := range paths {
		dirExpanded := os.ExpandEnv(dir)
		currentPath := path.Join(dirExpanded, filePathExpanded)
		slog.Info(
			"dynamicFile",
			logger.LogFieldFilePath, filePath,
			logger.LogFieldDirPath, dir,
			logger.LogFieldDirPathExpanded, dirExpanded,
			logger.LogFieldFilePathExpanded, currentPath,
		)
		if err := files.FileExists(currentPath); err == nil {
			return currentPath
		}
	}

	log.Fatalf("error: %v", ErrFileNotFound(filePathExpanded, paths))
	return ""
}
