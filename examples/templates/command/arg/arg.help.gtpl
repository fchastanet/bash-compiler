{{define "arg.help" -}}
{{ with .Data -}}
{{ if eq .min 0 }}[{{ end -}}
${__HELP_OPTION_COLOR}{{- .name -}}${__HELP_NORMAL}
{{- if eq .max 1 -}} {single}{{- if eq .min 1 }} (mandatory){{ end -}}
{{- else -}}{{/*
  */}} {list}{{/*
  */}}{{ if gt .min 0 }} (at least {{ .min }} times){{ else }} (optional){{ end -}}{{/*
  */}}{{ if gt .max 0 }} (at most {{ .max }} times){{ end -}}{{/*
  */}}
{{- end -}}
{{- if eq .min 0 -}}]{{- end -}}
{{ end -}}
{{ end }}
