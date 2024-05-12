// Package functions
package functions

import (
	"log"

	sprig "github.com/Masterminds/sprig/v3"
)

func fatal(message string) string {
	log.Fatalf("template error: %s", message)
	return message
}

func FuncMap() map[string]interface{} {
	funcMap := sprig.FuncMap()
	funcMap["fatal"] = fatal
	// templates functions
	funcMap["include"] = include
	// YAML functions
	funcMap["fromYAMLFile"] = FromYAMLFile

	return funcMap
}
