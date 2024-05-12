{{define "option.parse.before" -}}
{{/*
{{- $variableType := coalesce .Data.variableType "Boolean" -}}
{{- if eq $variableType "Boolean" -}}
  {{- $offValue := coalesce .Data.offValue "0" -}}
  {{- $min := 0 -}}
  {{- $max := 1 -}}
  {{- .Data.variableName }="{{ coalesce .Data.offValue ""}}"
{{- else if eq $variableType "String" -}}
{{ else }}
{{ .Data.variableName }="{{ coalesce .Data.defaultValue "" }}"
{{ end }}
{{- $min := coalesce .Data.min "Boolean" -}}
{{- $max := coalesce .Data.max "Boolean" -}}
{{ if ((min > 0 || max > 0)); then
local -i options_parse_optionParsedCount<% ${variableName^} %>
((options_parse_optionParsedCount<% ${variableName^} %> = 0)) || true
% fi
 */}}
{{ end }}
