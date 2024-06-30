{{define "option.help" -}}
${__HELP_OPTION_COLOR}{{/*
*/}}{{ .Data.alts | join "${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}" -}}{{/*
*/}}{{- $type := coalesce .Data.type "Boolean" -}}{{/*
*/}}{{ if ne $type "Boolean" }} {{ .Data.helpValueName }}{{end -}}{{/*
*/}}${__HELP_NORMAL}{{/*
*/}}{{ if eq .Data.max 1 }} {single}{{ if eq .Data.min 1 }} (mandatory){{ end -}}{{/*
*/}}{{ else }}{{/*
  */}} {list}{{/*
  */}}{{- $min := default 0 .Data.min -}}{{/*
  */}}{{- $max := default -1 .Data.max -}}{{/*
  */}}{{ if gt $min 0 }} (at least {{ $min }} times){{ else }} (optional){{ end }}{{/*
  */}}{{ if gt $max 0 }} (at most {{ $max }} times) {{ end }}{{/*
*/}}{{- end }}
{{- end }}
