// Package render
package render

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/fchastanet/bash-compiler/internal/utils/bash"
	"github.com/fchastanet/bash-compiler/internal/utils/files"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

type fileNotFoundError struct {
	error
	File    string
	SrcDirs []string
}

func (e *fileNotFoundError) Error() string {
	return fmt.Sprintf("file does not exist: %s in any srcDirs %v", e.File, e.SrcDirs)
}

// include allows to include a template
// allowing to use filter
// Eg: {{ include "template.tpl" | indent 4 }}
func Include(
	template string,
	templateData any,
	templateContextData TemplateContextData,
) string {
	var output string
	output, _ = MustInclude(template, templateData, templateContextData)
	return output
}

func MustInclude(
	templateName string,
	templateData any,
	templateContextData TemplateContextData,
) (output string, err error) {
	slog.Debug("MustInclude",
		logger.LogFieldTemplateName, templateName,
		logger.LogFieldTemplateData, templateData,
	)
	templateContextData.Data = templateData
	output, err = templateContextData.TemplateContext.Render(&templateContextData, templateName)
	if logger.FancyHandleError(err) {
		return "", err
	}
	return bash.RemoveFirstShebangLineIfAny(output), err
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

func includeFileAsTemplate(
	filePath string,
	templateContextData TemplateContextData,
) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Debug(
		"includeFileAsTemplate",
		logger.LogFieldFilePath, filePath,
		logger.LogFieldFilePathExpanded, filePathExpanded,
	)

	fileContent, err := os.ReadFile(filePathExpanded)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	code, err := templateContextData.TemplateContext.RenderFromTemplateContent(&templateContextData, string(fileContent))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return code
}

func dynamicFile(filePath string, paths []string) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Debug(
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
		slog.Debug(
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

	log.Fatalf("error: %v", fileNotFoundError{nil, filePathExpanded, paths})
	return ""
}
