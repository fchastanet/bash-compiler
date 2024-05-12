{{define "dataModel.arg.common"}}
{{ with .Data }}
{{ $min := coalesce .min 1 }}
{{ $max := coalesce .max 1 }}
{{ $variableType := "StringArray" }}
{{ if and (eq $max 1) (le $min 1) }}
{{ $variableType := "String" }}
{{ end }}
variableType: {{ $variableType }}
{{ if .helpValueName -}}
helpValueName: {{ .helpValueName }}
{{ else -}}
{{ logWarn "help value name set to variable name as not being set" "variableName" .variableName -}}
helpValueName: {{ .variableName }}
{{ end }}
{{ if .name -}}
name: {{ .name }}
{{ else -}}
{{- logWarn "argument name set to variable name as not being set" "variableName" .variableName -}}
name: {{ .variableName }}
{{ end }}
{{ if .regexp -}}
regexp: {{ .regexp }}
{{ end }}
{{ if .group -}}
group: {{ .group }}
{{ end }}
{{- if .authorizedValues -}}
authorizedValues:
  {{- range .authorizedValues }}
  {{- if eq (kindOf .) "string" -}}
  - value: {{ . }}
    help: ""
  {{ else }}
  - value: {{.value | errorIfEmpty }}
    help: {{ coalesce .help "\"\"" }}
  {{- end -}}
  {{- end -}}
{{end}}{{/* if .authorizedValues */}}
min: {{ $min }}
max: {{ $max }}
{{ if and (ne $max -1) (gt $min $max) }}
  {{ fail (cat "max value " $min " should be greater than min value " $max) }}
{{ end }}
{{end}}{{/* with .Data */}}
{{end}}{{/* define */}}
