{{- define "arg.oneLineHelp" -}}
{{- with .Data -}}
# {{ .variableName }} min {{ .min }} max {{ .max }}
{{ if .authorizedValues -}}
# authorizedValues: {{ .authorizedValuesList | join "|" }}
{{ end -}}
{{ if .regexp -}}
# regexp: '{{ .regexp }}'
{{ end -}}
{{ end -}}
{{ end }}
