{{- define "option.parse.option" -}}
{{- $context := . -}}
{{- with .Data -}}
{{- $Data := . -}}
{{ .alts | join " | " }})
  {{ if eq .type "Boolean" }}
  # shellcheck disable=SC2034
  {{ .variableName }}="{{ .onValue }}"
  {{ else }}
  shift
  if (($# == 0)); then
    Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
    return 1
  fi
  {{ if .authorizedValuesList }}
  if [[ ! "$1" =~ {{ .authorizedValuesList | join "|" }} ]]; then
    Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - value '$1' is not part of authorized values({{ .authorizedValuesList }})"
    return 1
  fi
  {{ end }}
  {{ end }}
  {{ if gt .max 0 }}
  if ((options_parse_optionParsedCount{{ .variableName | title }} >= {{ .max }} )); then
    Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached({{ .max }})"
    return 1
  fi
  {{ end }}
  ((++options_parse_optionParsedCount{{ .variableName | title }}))
  {{ if eq .type "String" -}}
  # shellcheck disable=SC2034
  {{ .variableName }}="$1"
  {{ else if eq .type "StringArray" -}}
  {{ .variableName }}+=("$1")
  {{ end }}
  {{ range .callbacks }}
  {{ if eq $Data.type "StringArray" -}}
    {{ . }} "${options_parse_arg}" "${ {{- $Data.variableName }}[@]}"
  {{ else -}}
    {{ . }} "${options_parse_arg}" "${ {{- $Data.variableName }}}"
  {{ end }}
  {{ end }}
  {{ range .everyOptionCallbacks }}
  {{ if eq $Data.type "StringArray" -}}
    {{ . }} "${options_parse_arg}" "${ {{- $Data.variableName }}[@]}"
  {{ else -}}
    {{ . }} "${options_parse_arg}" "${ {{- $Data.variableName }}}"
  {{ end -}}
  {{ end -}}
  ;;
{{ end -}}
{{ end -}}
