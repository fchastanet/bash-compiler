{{- define "option.oneLineHelp" -}}
{{- with .Data -}}
# {{ .variableName }} alts {{ .alts | join "|" }}
# type: {{ .type }} min {{ .min }} max {{ .max }}
{{ if .authorizedValuesList -}}
# authorizedValues: {{ .authorizedValuesList | join "|" }}
{{ end }}
{{ if .regexp -}}
# regexp: '{{ .regexp }}'
{{ end }}
{{ end }}
{{ end }}
