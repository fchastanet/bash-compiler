{{define "option.help.alts" -}}
  {{- $type := coalesce .Data.type "Boolean" -}}
  {{- if eq $type "Boolean" -}}
    [{{ .Data.alts | join "|" }}]
  {{- else -}}
    {{- if eq .Data.min 0 -}}[{{ end -}}{{/*
    */}}{{- .Data.alts | join "|" -}}{{/*
    */}} <{{- .Data.helpValueName | default "value" -}}>{{/*
    */}}{{ if eq .Data.min 0 }}]{{ end -}}
  {{- end -}}
{{- end -}}
