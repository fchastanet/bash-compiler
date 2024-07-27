// Package render
package render

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	matchFirstCap  = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap    = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchAllColons = regexp.MustCompile("([A-Za-z0-9])::([a-zA-Z0-9])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchAllColons.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func FirstCharacterTitle(str string) string {
	titleCases := cases.Title(language.AmericanEnglish)
	return titleCases.String(str)
}
