package bash

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

var shebangRegexp = regexp.MustCompile(`^(#!/.*)?$`)

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
