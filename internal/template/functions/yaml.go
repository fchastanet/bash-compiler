// Package functions
package functions

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

// FromYAMLFile loads yaml file and decodes it
func FromYAMLFile(filePath string) interface{} {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("FromYAMLFile %s err #%v ", filePath, err)
	}
	return fromYAML(yamlFile)
}

// fromYAML decodes YAML into a structured value, ignoring errors.
func fromYAML(data []byte) interface{} {
	output, _ := MustFromYAML(data)
	return output
}

// mustFromYAML decodes JSON into a structured value, returning errors.
func MustFromYAML(data []byte) (interface{}, error) {
	var output interface{}
	err := yaml.Unmarshal(data, &output)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return output, err
}
