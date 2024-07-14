package compiler

import "regexp"

var (
	functionsDirectiveRegexp    = regexp.MustCompile(`^# FUNCTIONS$`)
	commentRegexp               = regexp.MustCompile(`^[[:blank:]]*(#.*)?$`)
	bashFrameworkFunctionRegexp = regexp.MustCompile(
		`(?P<funcName>([A-Z]+[A-Za-z0-9_-]*::)+([a-zA-Z0-9_-]+))`)
)

func IsFunctionDirective(line []byte) bool {
	return functionsDirectiveRegexp.Match(line)
}

func IsCommentLine(line []byte) bool {
	return commentRegexp.Match(line)
}

func IsBashFrameworkFunction(line []byte) bool {
	return bashFrameworkFunctionRegexp.Match(line)
}
