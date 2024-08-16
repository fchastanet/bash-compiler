{{- define "arg.parse.args" -}}
{{- $context := . -}}
{{- with .Data -}}
{{- $Data := . -}}
{{- $minParsedArgIndex := 0 }}
{{- $maxParsedArgIndex := 0 }}
{{- $argCount := len .args }}
{{- $incrementArg := 1 }}
{{ if .args }}
if ((0)); then
  # Technical if - never reached
  :
{{ range $index, $arg := .args }}
# Argument {{ add $index 1 }}/{{ $argCount }}
{{ include "arg.oneLineHelp" $arg $context | trim }}
{{ $minParsedArgIndex := add $minParsedArgIndex .max }}
{{ if eq .max -1 }}
elif (( options_parse_parsedArgIndex >= {{ $maxParsedArgIndex }} )); then
{{ else }}
elif (( options_parse_parsedArgIndex >= {{ $maxParsedArgIndex }} &&
  options_parse_parsedArgIndex < {{ $minParsedArgIndex }} )); then
{{ end }}
  {{ include "arg.parse.arg" $arg $context | indent 2 | trim }}
  {{ range .everyArgumentCallbacks }}
  # shellcheck disable=SC2317
  {{ . }} "{{ .variableName }}" "${options_parse_arg}"|| true
  {{ end -}}
  {{ $maxParsedArgIndex := add $maxParsedArgIndex .max }}
{{ end -}}{{/* end range args */}}
# else too much args
else
  {{ range .everyArgumentCallbacks }}
  # shellcheck disable=SC2317
  {{ . }} "" "${options_parse_arg}"|| argOptDefaultBehavior=$?
  {{ end -}}
  {{ if .unknownArgumentCallbacks }}
  # no arg configured, call unknownArgumentCallback
  {{ range .unknownArgumentCallbacks }}
  # shellcheck disable=SC2317
  {{ . }} "${options_parse_arg}"
  {{ end -}}
  {{ else }}
  if [[ "${argOptDefaultBehavior}" = "0" ]]; then
    # too much args and no unknownArgumentCallbacks configured
    Log::displayError "Command ${SCRIPT_NAME} - Argument - too much arguments provided: $*"
    return 1
  fi
  {{end}}
fi
{{ else }}{{/* No args declared */}}
{{ if .unknownArgumentCallbacks }}
# no arg configured, call unknownArgumentCallback
{{ range .unknownArgumentCallbacks }}
# shellcheck disable=SC2317
{{ . }} "${options_parse_arg}"
{{ end -}}
{{ range .everyArgumentCallbacks }}
# shellcheck disable=SC2317
{{ . }} "${options_parse_arg}"
{{ end -}}
{{ else }}{{/* No args declared and no unknownArgumentCallbacks set */}}
{{ range .everyArgumentCallbacks }}
# shellcheck disable=SC2317
{{ . }} "${options_parse_arg}" || argOptDefaultBehavior=\$?
{{ end -}}
if [[ "${argOptDefaultBehavior}" = "0" ]]; then
  # too much args and no unknownArgumentCallbacks configured
  Log::displayError "Command ${SCRIPT_NAME} - Argument - too much arguments provided: $*"
  return 1
fi
{{ $incrementArg := 0 }}{{/* to avoid parse error after return */}}
{{end}}{{/* if .unknownArgumentCallbacks */}}
{{ if eq $incrementArg 1 }}
((++options_parse_parsedArgIndex))
{{ end }}

{{ end }}
{{ end -}}
{{ end -}}
