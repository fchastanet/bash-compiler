{{- define "checkRequirements" -}}
{{ $context := . -}}
{{ $replace := include "checkRequirementsCode" .Data $context }}
{{ $regexp := print "(?m)[ \t]*(function[ \t]+|)(" .Data.functionName ")\\(\\)[ \t]*\\{[ \t]*$" }}
{{ $replace | regexReplaceAll $regexp .Data.code }}
{{- end }}
