{{- define "arg.parse.after" -}}
{{- with .Data -}}
{{ if gt .min 0 }}
if ((options_parse_argParsedCount{{ .variableName | title }} < {{ .min }} )); then
  Log::displayError "Command ${SCRIPT_NAME} - Argument '{{ .name }}' should be provided at least {{ .min }} time(s)"
  return 1
fi
{{ end }}
{{ end }}
{{ end }}
