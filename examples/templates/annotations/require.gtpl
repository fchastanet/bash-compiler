{{- define "require" -}}
{{ $functionNameUpper := .Data.functionName | snakeCase -}}
{{ $replace := print
   .Data.functionName "() {" "\n"
   "  export REQUIRE_FUNCTION_" $functionNameUpper "_LOADED=1" "\n"
}}
{{- $regexp := print "(?m)[ \t]*(function[ \t]+|)(" .Data.functionName ")\\(\\)[ \t]*\\{[ \t]*$" -}}
{{ $replace | regexReplaceAll $regexp .Data.code }}
{{- end }}
