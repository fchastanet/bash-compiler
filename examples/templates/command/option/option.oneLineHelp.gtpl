{{- define "option.oneLineHelp" -}}
{{- with .Data -}}
# {{ .variableName }} alts {{ .alts | join "|" }}
# type: {{ .type }} min {{ .min }} max {{ .max }}
{{ if .authorizedValues -}}
# authorizedValues: {{ $sep := ""
  -}}{{- range .authorizedValues}}{{$sep}}{{.value}}{{$sep = "|"}}{{- end }}
{{ end }}
{{ if .regexp -}}
# regexp: '{{ .regexp }}'
{{ end }}
{{ end }}
{{ end }}
