{{- define "arg.parse.before" -}}
{{- with .Data -}}
{{- if eq .type "StringArray" -}}
{{   .variableName }}=()
{{- else -}}
{{   .variableName }}=""
{{- end -}}
{{- if or (gt .min 0) (gt .max 0) }}
local -i options_parse_argParsedCount{{ .variableName | title}}
((options_parse_argParsedCount{{ .variableName | title}} = 0)) || true
{{ end }}
{{ end }}
{{ end }}
