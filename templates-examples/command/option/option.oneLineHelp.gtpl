{{- define "option.oneLineHelp" -}}
{{- with .Data -}}
# {{ .variableName }} alts {{ .alts | join "|" }}
# variableType: {{ .variableType }} min {{ .min }} max {{ .max }}
{{ if .authorizedValues -}}
# authorizedValues: {{ .authorizedValuesList | join "|" }}
{{ end }}
{{ if .regexp -}}
# regexp: '{{ .regexp }}'
{{ end }}
{{ end }}
{{ end }}
