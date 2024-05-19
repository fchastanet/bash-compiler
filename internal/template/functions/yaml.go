// Package functions
package functions

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const IndentSize = 2

// FromYAMLFile loads yaml file and decodes it
func FromYAMLFile(filePath string) interface{} {
	model, err := MustFromYAMLFile(filePath)
	if err != nil {
		slog.Error(fmt.Sprintf("FromYAMLFile %s err #%v ", filePath, err))
	}
	return model
}

// FromYAMLFile loads yaml file and decodes it
func MustFromYAMLFile(filePath string) (interface{}, error) {
	yamlFileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return MustFromYAML(yamlFileContent)
}

// mustFromYAML decodes JSON into a structured value, returning errors.
func MustFromYAML(data []byte) (interface{}, error) {
	var output interface{}
	err := yaml.Unmarshal(data, &output)
	if err != nil {
		slog.Error("error during Unmarshalling", "error", err)
	}
	return output, err
}

// ToYAML decodes YAML into a structured value, ignoring errors.
func ToYAML(data []interface{}) string {
	output, _ := yaml.MarshalWithOptions(data, yaml.Indent(IndentSize), yaml.IndentSequence(true))
	slog.Info("-------------------------------------------------------------")
	fmt.Printf("%s\n", string(output))
	return string(output)
}
