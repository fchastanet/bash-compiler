{{- define "option.parse.before" -}}
{{- with .Data -}}
{{- if eq .variableType "Boolean" -}}
{{ .variableName }}="{{ .offValue }}"
{{- else if eq .variableType "String" -}}
{{ .variableName }}="{{ .defaultValue }}"
{{- end -}}
{{- if or (gt .min 0) (gt .max 0) }}
local -i options_parse_optionParsedCount{{ .variableName | title}}
((options_parse_optionParsedCount{{ .variableName | title}} = 0)) || true
{{ end }}
{{ end }}
{{ end }}
