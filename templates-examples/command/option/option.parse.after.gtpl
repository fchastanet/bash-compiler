{{- define "option.parse.after" -}}
{{- with .Data -}}
{{ if gt .min 0 }}
if ((options_parse_optionParsedCount{{ .variableName | title }} < {{ .min }} )); then
  Log::displayError "Command ${SCRIPT_NAME} - Option '{{ .alts | first}}' should be provided at least {{ .min }} time(s)"
  return 1
fi
{{ end }}
{{ end }}
{{ end }}
