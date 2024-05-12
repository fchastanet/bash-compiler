{{define "dataModel.option.StringCommon"}}
{{- with .Data }}
{{ if .helpValueName -}}
helpValueName: {{ .helpValueName }}
{{ else -}}
{{- logWarn "help value name set to variable name as not being set" "variableName" .variableName -}}
helpValueName: {{ .variableName }}
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
{{end}}
{{ end }}
{{end}}
