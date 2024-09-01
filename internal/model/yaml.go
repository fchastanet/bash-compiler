package model

import (
	"log/slog"
	"os"

	"github.com/fchastanet/bash-compiler/internal/utils/logger"
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
	slog.Debug("Temp file containing resulting yaml file", logger.LogFieldFilePath, tempYamlFile.Name())
	return nil
}
