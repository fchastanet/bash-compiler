// Package render
package render

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	replacePattern = "${1}_${2}"
)

var (
	matchFirstCap  = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap    = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchAllColons = regexp.MustCompile("::")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, replacePattern)
	snake = matchAllCap.ReplaceAllString(snake, replacePattern)
	snake = matchAllColons.ReplaceAllString(snake, replacePattern)
	return strings.ToUpper(snake)
}

func FirstCharacterTitle(str string) string {
	titleCases := cases.Title(language.AmericanEnglish)
	return titleCases.String(str)
}
