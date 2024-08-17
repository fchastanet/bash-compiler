{{define "command.parse"}}
{{- $context := . -}}
{{- with .Data -}}

{{- range $index, $option := .options -}}
{{include "option.parse.before" $option $context | trim}}
{{ end -}}
{{ range $index, $arg := .args }}
{{include "arg.parse.before" $arg $context | trim}}
{{ end }}

# shellcheck disable=SC2034
local -i options_parse_parsedArgIndex=0
while (($# > 0)); do
  local options_parse_arg="$1"
  local argOptDefaultBehavior=0
  case "${options_parse_arg}" in
    {{ $optionsCount := len .options }}
    {{- range $index, $option := .options -}}
    # Option {{ add $index 1 }}/{{ $optionsCount }}
    {{ include "option.oneLineHelp" $option $context | indent 4 | trim }}
    {{ include "option.parse.option" $option $context | indent 4 | trim }}
    {{ end }}
    -*)
      {{ range .everyOptionCallbacks }}
      # shellcheck disable=SC2317
      {{ . }} "${options_parse_arg}" || argOptDefaultBehavior=$?
      {{ end -}}
      {{ if .unknownOptionCallbacks }}
      {{ range .unknownOptionCallbacks }}
      {{ . }} "${options_parse_arg}" || argOptDefaultBehavior=$?
      {{ end -}}
      {{ else }}
      if [[ "${argOptDefaultBehavior}" = "0" ]]; then
        Log::displayError "Command ${SCRIPT_NAME} - Invalid option ${options_parse_arg}"
        return 1
      fi
      {{ end }}
      ;;
    *)
      {{ include "arg.parse.args" . $context | indent 6 | trim }}
      ;;
  esac
  shift || true
done
{{- range $index, $option := .options -}}
{{include "option.parse.after" $option $context | trim}}
{{ end -}}
{{ range $index, $arg := .args }}
{{include "arg.parse.after" $arg $context | trim}}
{{ end }}
{{ range .commandCallbacks }}
# shellcheck disable=SC2317
{{ . }}
{{ end -}}
{{ end }}
{{end}}
