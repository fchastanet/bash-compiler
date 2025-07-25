package dotenv

import (
	"bufio"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/fchastanet/bash-compiler/internal/utils/logger"
)

var (
	commentRegexp     = regexp.MustCompile(`^[ \t]*#`)
	variableSetRegexp = regexp.MustCompile(`^[ \t]*(?P<name>[A-Za-z_]+)=(?P<value>.*)$`)
)

func LoadEnvFile(confFile string) error {
	confFileContent, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer errors.SafeCloseDeferCallback(confFileContent, &err)

	variables := make(map[string]string)
	scanFile(confFileContent, variables)

	err = interpolateVariables(variables)
	if logger.FancyHandleError(err) {
		return err
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
			slog.Warn("Ignore invalid line",
				logger.LogFieldLineNumber, lineNumber,
				logger.LogFieldLineContent, line,
			)
			continue
		}

		name := matches[variableSetRegexpNameGroupIndex]
		value := matches[variableSetRegexpValueGroupIndex]
		if _, ok := variables[name]; ok {
			slog.Warn("overwriting variable",
				logger.LogFieldLineNumber, lineNumber,
				logger.LogFieldVariableName, name,
			)
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
		// remove " or '
		valueInterpolated = strings.Trim(valueInterpolated, "'\"")
		slog.Debug(
			"Variable interpolated value",
			logger.LogFieldVariableName, name,
			logger.LogFieldVariableValue, valueInterpolated,
		)
		variables[name] = valueInterpolated
		os.Setenv(name, valueInterpolated)
	}

	return nil
}
