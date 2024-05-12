{{define "dataModel.parameter"}}
{{- with .Data -}}
variableName: {{ .variableName | errorIfEmpty }}
{{ if .functionName -}}
functionName: {{ .functionName -}}
{{- else -}}
functionName: {{ .variableName -}}Function
{{- end }}
{{ if .callbacks -}}
callbacks:
  {{ range .callbacks -}}
  - {{ . -}}
  {{end}}
{{end}}
help: |
  {{ .help | indent 2 | trim }}
{{ end }}
{{ end }}
