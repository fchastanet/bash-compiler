// Package functions
package functions

import "strings"

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.ReplaceAll(v, "\n", "\n"+pad)
}

func nindent(spaces int, v string) string {
	return "\n" + indent(spaces, v)
}
