{{define "option.help.alts" -}}
  {{- $variableType := coalesce .Data.variableType "Boolean" -}}
  {{- if eq $variableType "Boolean" -}}
    {{ .Data.alts | join "|" }}
  {{- else -}}
    {{- if eq .Data.min 0 -}}[{{ end -}}{{/*
    */}}{{- .Data.alts | join "|" -}}{{/*
    */}} <{{- .Data.helpValueName | default "value" -}}>{{/*
    */}}{{ if eq .Data.min 0 }}]{{ end }}
  {{- end -}}
{{ end }}