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
	templateContext myTemplate.Context, overrideRootData bool) string {
	var output string
	output, _ = mustInclude(template, templateData, templateContext, overrideRootData)
	return output
}

func mustInclude(template string, templateData any,
	templateContext myTemplate.Context, overrideRootData bool) (string, error) {
	var output string
	var err error
	if overrideRootData {
		templateContext.RootData = templateData
	}
	output, err = templateContext.Render(template, templateData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
