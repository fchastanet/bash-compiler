{{define "dataModel.parameter"}}
{{- with .Data -}}
variableName: {{ .variableName | errorIfEmpty }}
{{ if .functionName -}}
functionName: {{ .functionName -}}
{{- else -}}
functionName: {{ .variableName -}}Function
{{- end }}
{{ if not .variableType -}}
{{- logWarn "variable type set to Boolean by default" "variableName" .variableName -}}
{{- end }}
{{- $variableType := coalesce .variableType "Boolean" -}}
variableType: {{ $variableType }}
help: |
  {{ .help | indent 2 | trim }}
{{ end }}
{{ end }}
