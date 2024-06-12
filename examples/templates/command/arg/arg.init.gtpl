{{define "arg.init" -}}
{{- with .Data -}}
{{ if eq .type "StringArray" -}}
declare -a {{ .variableName }}=()
{{ else -}}
{{ if ne .defaultValue nil -}}
declare {{ .variableName }}={{ .defaultValue | quote }}
{{ else -}}
declare {{ .variableName }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
