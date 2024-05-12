{{define "command.dataModel.String" -}}
{{- with .Data -}}
{{ $min := int (coalesce .min 0) }}
{{- if .mandatory -}}
{{ $min := 1 }}
{{- end -}}
{{- if gt $min 1 -}}
  {{ fail "min cannot be greater than 1" }}
{{ end }}
{{ if .defaultValue -}}
defaultValue: {{ .defaultValue }}
{{ else -}}
{{- logWarn "default value set to empty string" "variableName" .variableName -}}
defaultValue: ""
{{ end -}}
min: {{ $min }}
max: 1
{{ end }}
{{end}}
