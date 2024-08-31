{{- define "arg.parse.args" -}}
{{- $context := . -}}
{{- with .Data -}}
{{-   $Data := . -}}
{{    if .args }}
{{-     $argCount := len .args }}
((minParsedArgIndex0 = 0)) || true
((maxParsedArgIndex0 = 0)) || true
{{    range $index, $arg := .args -}}
((minParsedArgIndex{{ add $index 1}} = minParsedArgIndex{{$index}} + {{.min}})) || true
{{      if eq .max -1 -}}
((maxParsedArgIndex{{ add $index 1}} = maxParsedArgIndex{{$index}})) || true
{{      else -}}
((maxParsedArgIndex{{ add $index 1}} = maxParsedArgIndex{{$index}} + {{.max}})) || true
{{      end -}}
{{    end -}}
((incrementArg = 1 ))
if ((0)); then
  # Technical if - never reached
  :
{{      range $index, $arg := .args }}
# Argument {{ add $index 1 }}/{{ $argCount }} - {{ $arg.variableName }}
{{        include "arg.oneLineHelp" $arg $context | trim }}
{{        if eq .max -1 -}}
elif (( options_parse_parsedArgIndex >= minParsedArgIndex{{$index}} )); then
{{        else -}}
elif (( options_parse_parsedArgIndex >= minParsedArgIndex{{$index}} &&
  options_parse_parsedArgIndex < maxParsedArgIndex{{add $index 1}} )); then
{{        end -}}
{{        include "arg.parse.arg" $arg $context | indent 2 | trimSuffix " " }}
{{        range .everyArgumentCallbacks }}
  # shellcheck disable=SC2317
  {{ . }} "{{ .variableName }}" "${options_parse_arg}"|| true
{{        end -}}
{{      end -}}{{/* end range args */}}
# else too much args
else
{{      range .everyArgumentCallbacks }}
  # shellcheck disable=SC2317
{{        . }} "" "${options_parse_arg}"|| argOptDefaultBehavior=$?
{{      end -}}{{/* end range .everyArgumentCallbacks */}}
{{      if .unknownArgumentCallbacks }}
  # no arg configured, call unknownArgumentCallback
{{        range .unknownArgumentCallbacks }}
  # shellcheck disable=SC2317
{{          . }} "${options_parse_arg}"
{{        end -}}
{{      else }}{{/* else if .unknownArgumentCallbacks */}}
  if [[ "${argOptDefaultBehavior}" = "0" ]]; then
    # too much args and no unknownArgumentCallbacks configured
    Log::displayError "Command ${SCRIPT_NAME} - Argument - too much arguments provided: $*"
    return 1
  fi
{{      end }}{{/* end if .unknownArgumentCallbacks */}}
fi
{{    else }}{{/* else if .args - No args declared */}}
{{      if .unknownArgumentCallbacks }}
# no arg configured, call unknownArgumentCallback
{{        range .unknownArgumentCallbacks }}
# shellcheck disable=SC2317
{{          . }} "${options_parse_arg}"
{{        end -}}{{/* range .unknownArgumentCallbacks */}}
{{        range .everyArgumentCallbacks }}
# shellcheck disable=SC2317
{{          . }} "${options_parse_arg}"
{{        end -}}{{/* range .everyArgumentCallbacks */}}
{{      else }}{{/* if .unknownArgumentCallbacks - No args declared and no unknownArgumentCallbacks set */}}
{{        range .everyArgumentCallbacks }}
# shellcheck disable=SC2317
{{          . }} "${options_parse_arg}" || argOptDefaultBehavior=\$?
{{        end -}}{{/* range .everyArgumentCallbacks */}}
if [[ "${argOptDefaultBehavior}" = "0" ]]; then
  # too much args and no unknownArgumentCallbacks configured
  Log::displayError "Command ${SCRIPT_NAME} - Argument - too much arguments provided: $*"
  return 1
fi
((incrementArg = 0)){{/* to avoid parse error after return */}}
{{      end -}}{{/* if .unknownArgumentCallbacks */ -}}
{{    end -}}{{/* end if .args */ -}}
if ((incrementArg == 1)); then
  ((++options_parse_parsedArgIndex))
fi
{{   end -}}{{/* end with .Data */}}
{{ end -}}
