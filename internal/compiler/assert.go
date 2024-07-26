package compiler

import (
	"fmt"
	"regexp"
)

const bashFrameworkFunctionRegexpStr = `(?P<funcName>([A-Z]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))`

var (
	functionsDirectiveRegexp        = regexp.MustCompile(`^# FUNCTIONS$`)
	commentRegexp                   = regexp.MustCompile(`^[[:blank:]]*(#.*)?$`)
	bashFrameworkFunctionRegexp     = regexp.MustCompile(bashFrameworkFunctionRegexpStr)
	fullBashFrameworkFunctionRegexp = regexp.MustCompile(
		fmt.Sprintf(`^%s$`, bashFrameworkFunctionRegexpStr),
	)
)

func IsFunctionDirective(line []byte) bool {
	return functionsDirectiveRegexp.Match(line)
}

func IsCommentLine(line []byte) bool {
	return commentRegexp.Match(line)
}

func IsBashFrameworkFunction(line []byte) bool {
	return fullBashFrameworkFunctionRegexp.Match(line)
}
