{{- define "commands" -}}
{{ $context := . }}
{{ range .Data }}
{{- include "command" . $context -}}
{{end}}
{{end}}
