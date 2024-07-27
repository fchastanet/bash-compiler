package model

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

func writeYamlFile(modelMap map[string]interface{}, tempYamlFile os.File) (err error) {
	yamlResult, err := yaml.Marshal(modelMap)
	if err != nil {
		return err
	}
	_, err = tempYamlFile.Write(yamlResult)
	if err != nil {
		return err
	}
	log.Printf("Temp file containing resulting yaml file : %s\n", tempYamlFile.Name())
	return nil
}
