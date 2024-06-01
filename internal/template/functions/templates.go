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

	"github.com/fchastanet/bash-compiler/internal/template"
	render "github.com/fchastanet/bash-compiler/internal/template"
	"github.com/fchastanet/bash-compiler/internal/utils"
)

var errFileNotFound = errors.New("File does not exist")

func ErrFileNotFound(file string, srcDirs []string) error {
	return fmt.Errorf("%w: %s in any srcDirs %v", errFileNotFound, file, srcDirs)
}

// include allows to include a template
// allowing to use filter
// Eg: {{ include "template.tpl" | indent 4 }}
func include(
	template string, templateData any,
	templateContext render.Context) string {
	var output string
	output, _ = mustInclude(template, templateData, templateContext)
	return output
}

func mustInclude(templateName string, templateData any,
	templateContext render.Context) (output string, err error) {
	slog.Info("mustInclude", "templateName", templateName, "templateData", templateData)
	templateContext.Data = templateData
	output, err = templateContext.Render(templateName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}

func includeFile(filePath string) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Info("includeFile", "filePath", filePath, "filePathExpanded", filePathExpanded)

	file, err := os.ReadFile(filePathExpanded)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(file)
}

func RenderFromTemplateContent(
	templateContext *template.Context, templateContent string,
) (code string, err error) {
	template, err := templateContext.Template.Parse(templateContent)
	if err != nil {
		return "", err
	}
	var tplWriter bytes.Buffer
	err = template.Execute(&tplWriter, templateContext)
	if err != nil {
		return "", err
	}

	return tplWriter.String(), err
}

func includeFileAsTemplate(filePath string, templateContext render.Context) string {
	filePathExpanded := os.ExpandEnv(filePath)
	slog.Info("includeFileAsTemplate", "filePath", filePath, "filePathExpanded", filePathExpanded)

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
	slog.Info("dynamicFile", "filePath", filePath, "filePathExpanded", filePathExpanded)
	err := utils.FileExists(filePathExpanded)
	if err == nil {
		return filePathExpanded
	}
	for _, dir := range paths {
		dirExpanded := os.ExpandEnv(dir)
		currentPath := path.Join(dirExpanded, filePathExpanded)
		slog.Info("dynamicFile", "filePath", filePath, "dir", dir, "dirExpanded", dirExpanded, "currentPath", currentPath)
		if err := utils.FileExists(currentPath); err == nil {
			return currentPath
		}
	}

	log.Fatalf("error: %v", ErrFileNotFound(filePathExpanded, paths))
	return ""
}
