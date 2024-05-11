// Package functions
package functions

import (
	"log"

	myTemplate "github.com/fchastanet/bash-compiler/internal/template"
)

// include allows to include a template
// allowing to use filter
// Eg: {{ include "template.tpl" | indent 4 }}
func include(
	template string, templateData any,
	templateContext myTemplate.Context) string {
	var output string
	output, _ = mustInclude(template, templateData, templateContext)
	return output
}

func mustInclude(template string, templateData any,
	templateContext myTemplate.Context) (string, error) {
	var output string
	var err error

	templateContext.Data = templateData
	output, err = templateContext.Render(template)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
