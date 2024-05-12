{{define "option.values" -}}
{{- $variableName := .Data.variableName -}}
{{- $variableType := coalesce .Data.variableType "Boolean" -}}
{{- if eq $variableType "Boolean" -}}
  {{- $offValue := coalesce .Data.offValue "0" -}}
  {{- $min := 0 -}}
  {{- $max := 1 -}}
{{- else -}}
  {{- $defaultValue := coalesce .Data.defaultValue "0" -}}
  {{- $min := coalesce .Data.min 0 -}}
  {{- $max := coalesce .Data.max -1 -}}
  {{ if and (ne $max -1) (gt $min $max) }}
    {{ fatal "option {{ $variableName }} --max value should be greater than --min value" }}
  {{ end }}
{{ end }}
{{ end }}