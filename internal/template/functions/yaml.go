// Package functions
package functions

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

// fromYAMLFile loads yaml file and decodes it
func fromYAMLFile(filePath string) interface{} {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("fromYAMLFile %s err #%v ", filePath, err)
	}
	return fromYAML(yamlFile)
}

// fromYAML decodes YAML into a structured value, ignoring errors.
func fromYAML(data []byte) interface{} {
	output, _ := mustFromYAML(data)
	return output
}

// mustFromYAML decodes JSON into a structured value, returning errors.
func mustFromYAML(data []byte) (interface{}, error) {
	var output interface{}
	err := yaml.Unmarshal(data, &output)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
