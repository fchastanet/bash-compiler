{{- define "arg.oneLineHelp" -}}
{{- with .Data -}}
# Argument {{ .variableName }} min {{ .min }} max {{ .max }}
{{ if .authorizedValues -}}
# Argument {{ .variableName }} authorizedValues: {{ .authorizedValuesList | join "|" }}
{{ end -}}
{{ if .regexp -}}
# Argument {{ .variableName }} regexp: '{{ .regexp }}'
{{ end -}}
{{ end -}}
{{ end }}
