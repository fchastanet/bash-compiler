{{- define "arg.parse.arg" -}}
{{- $context := . -}}
{{- with .Data -}}
{{- $Data := . -}}
{{  if .authorizedValues -}}
{{- $sep := "" -}}
if [[ ! "${options_parse_arg}" =~ {{ range .authorizedValues}}{{$sep}}{{.value}}{{$sep = "|"}}{{- end }} ]]; then
  Log::displayError "Command ${SCRIPT_NAME} - Argument {{ .name }} - value '${options_parse_arg}' is not part of authorized values({{-
  $sep := "" -}}{{- range .authorizedValues}}{{$sep}}{{.value}}{{$sep = ", "}}{{- end }})"
  return 1
fi
{{  end -}}
{{  if .regexp -}}
if [[ ! "${options_parse_arg}" =~ {{ .regexp }} ]]; then
  Log::displayError "Command ${SCRIPT_NAME} - Argument {{ .name }} - value '${options_parse_arg}' doesn't match the regular expression({{ .regexp }})"
  return 1
fi
{{  end -}}
{{  if gt .max 0 -}}
if ((options_parse_argParsedCount{{ .variableName | title }} >= {{ .max }} )); then
  Log::displayError "Command ${SCRIPT_NAME} - Argument {{ .name }} - Maximum number of argument occurrences reached({{ .max }})"
  return 1
fi
{{  end -}}
((++options_parse_argParsedCount{{ .variableName | title }}))
# shellcheck disable=SC2034
{{  if eq .type "String" -}}
{{    .variableName }}="${options_parse_arg}"
{{    range .callbacks -}}
{{      . }} "{{ "${" }}{{ $Data.variableName }}{{ "}" }}" -- "${@:2}"
{{    end -}}
{{  else -}}
# shellcheck disable=SC2034
{{    .variableName }}+=("${options_parse_arg}")
{{    range .callbacks -}}
{{      . }} "{{ "${" }}{{ $Data.variableName }}{{ "[@]}" }}" -- "${@:2}"
{{    end -}}
{{  end -}}
{{ end -}}
{{ end -}}
