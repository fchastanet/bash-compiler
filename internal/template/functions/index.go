// Package functions
package functions

func FuncMap() map[string]interface{} {
	funcMap := map[string]interface{}{
		// templates functions
		"include": include,
		// string functions
		"indent":  indent,
		"nindent": nindent,
		// YAML functions
		"fromYAMLFile": FromYAMLFile,
	}
	return funcMap
}
