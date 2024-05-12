// Package functions
package functions

import (
	"log"
	"log/slog"

	render "github.com/fchastanet/bash-compiler/internal/template"
)

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
	templateContext.Data = &templateData
	output, err = templateContext.Render(templateName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
