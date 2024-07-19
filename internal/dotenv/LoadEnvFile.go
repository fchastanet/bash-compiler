package dotenv

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/a8m/envsubst"
	"github.com/fchastanet/bash-compiler/internal/logger"
)

var (
	commentRegexp     = regexp.MustCompile(`^[ \t]*#`)
	variableSetRegexp = regexp.MustCompile(`^[ \t]*(?P<name>[A-Za-z_]+)=(?P<value>.*)$`)
)

func LoadEnvFile(confFile string) error {
	confFileContent, err := os.OpenFile(confFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer confFileContent.Close()

	variables := make(map[string]string)
	scanFile(confFileContent, variables)

	err = interpolateVariables(variables)
	if logger.FancyHandleError(err) {
		return err
	}

	for name, value := range variables {
		os.Setenv(name, value)
	}
	return nil
}

func scanFile(confFileContent *os.File, variables map[string]string) {
	variableSetRegexpNameGroupIndex := variableSetRegexp.SubexpIndex("name")
	variableSetRegexpValueGroupIndex := variableSetRegexp.SubexpIndex("value")

	scanner := bufio.NewScanner(confFileContent)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Bytes()
		if commentRegexp.Match(line) {
			continue
		}
		matches := variableSetRegexp.FindStringSubmatch(string(line))
		if matches == nil {
			slog.Warn(fmt.Sprintf("Ignore invalid line %d : %s", lineNumber, line))
			continue
		}

		name := matches[variableSetRegexpNameGroupIndex]
		value := matches[variableSetRegexpValueGroupIndex]
		if _, ok := variables[name]; ok {
			slog.Warn(fmt.Sprintf("overwriting variable %s on line %d", name, lineNumber))
		}
		variables[name] = value
		lineNumber++
	}
}

func interpolateVariables(variables map[string]string) error {
	for name, value := range variables {
		valueInterpolated, err := envsubst.String(value)
		if logger.FancyHandleError(err) {
			return err
		}
		slog.Info(fmt.Sprintf("Variable %s interpolated value %s", name, valueInterpolated))
		variables[name] = valueInterpolated
	}

	return nil
}