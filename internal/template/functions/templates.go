// Package functions
package functions

import (
	"log"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
)

// include allows to include a template
// allowing to use filter
// Eg: {{ include "template.tpl" | indent 4 }}
func include(template string, templateData any) string {
	var output string
	output, _ = mustInclude(template, templateData)
	return output
}

func mustInclude(template string, templateData any) (string, error) {
	var output string
	var err error

	output, err = myTemplate.Render(
		myTemplate.TemplateDir, template, templateData, FuncMap())
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
