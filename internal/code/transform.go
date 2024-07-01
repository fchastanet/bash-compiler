package code

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	shebangRegexp = regexp.MustCompile(`^(#!/.*)?$`)
)

var errRequiredFunctionNotFound = errors.New("Required function not found")

func ErrRequiredFunctionNotFound(functionName string) error {
	return fmt.Errorf("%w: %s", errRequiredFunctionNotFound, functionName)
}

func RemoveFirstShebangLineIfAny(code string) string {
	scanner := bufio.NewScanner(strings.NewReader(code))
	var rewrittenCode bytes.Buffer
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Bytes()
		if lineNumber == 1 && shebangRegexp.Match(line) {
			rewrittenCode.WriteByte('\n')
			continue
		}
		rewrittenCode.Write(line)
		rewrittenCode.WriteByte('\n')
		lineNumber++
	}
	return rewrittenCode.String()
}
