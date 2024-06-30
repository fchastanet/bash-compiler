{{define "command.help"}}
{{ $context := . -}}
{{- with .Data }}
{{- if eq .helpType "function" -}}
Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
  "$({{ .help }})"
{{ else }}
Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
  "{{ .help }}"
{{ end -}}
echo

# ------------------------------------------
# usage section
# ------------------------------------------
Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" "{{- .commandName }} {{/*
  */}}{{- if .options -}} [OPTIONS]{{- end }} {{/*
  */}}{{- if .arguments -}} [ARGUMENTS]{{- end }}"
echo
{{ if .options -}}
# ------------------------------------------
# usage/options section
# ------------------------------------------
optionsAltList=({{ range $index, $option := .options -}}"{{- include "option.help.alts" $option $context | trimAll "\n" -}}" {{ end }}
)
Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" \
  "{{ .commandName }}" "${optionsAltList[@]}"
echo
{{ end }}

{{ if .arguments -}}
# ------------------------------------------
# usage/arguments section
# ------------------------------------------
echo
echo -e "${__HELP_TITLE_COLOR}ARGUMENTS:${__RESET_COLOR}"
{{ range $index, $arg := .args }}
  echo -e "{{ include "arg.help" $arg $context }}"{{/*
*/}}{{ end }}
{{ end -}}

{{ if .options -}}
# ------------------------------------------
# options section
# ------------------------------------------
{{ $previousGroupId := "" -}}
{{ $command := . }}
{{ range $index, $option := .options -}}
{{ $groupId := default "__default" $option.group -}}
{{ if ne $groupId $previousGroupId }}
echo
echo -e "${__HELP_TITLE_COLOR}{{ (index $command.optionGroups $groupId).title  }}${__RESET_COLOR}"
{{ end -}}
echo -e "  {{ include "option.help" $option $context | trimAll "\n" -}}"
{{ if $option.help -}}
Array::wrap2 ' ' 76 4 "    " {{ $option.help | quote }}
echo
{{ end }}
{{ $previousGroupId = $groupId }}
{{ end -}}
{{ end -}}

{{ if .longDescription -}}
# ------------------------------------------
# longDescription section
# ------------------------------------------
echo
{{ if hasSuffix "Function" .longDescription -}}
{{ .longDescription }}
{{ else -}}
declare -a {{ .functionName }}LongDescription=(
{{ $longDescriptionList := splitList "\n" .longDescription -}}
{{ range $line := $longDescriptionList -}}
{{ if hasPrefix "$'" $line -}}
{{ $line -}}
{{ else -}}
{{ quote $line -}}
{{ end }}
{{ end -}}
)
Array::wrap2 ' ' 76 0 "{{ format "${%sLongDescription[@]}" .functionName }}"
echo
{{ end -}}
{{ end -}}

{{ if .version -}}
# ------------------------------------------
# version section
# ------------------------------------------
echo
echo -n -e "${__HELP_TITLE_COLOR}VERSION: ${__RESET_COLOR}"
echo {{ .version | quote }}
{{ end -}}

{{ if .author -}}
# ------------------------------------------
# author section
# ------------------------------------------
echo
echo -n -e "${__HELP_TITLE_COLOR}AUTHOR: ${__RESET_COLOR}"
echo {{ .author | quote }}
{{ end -}}

{{ if .sourceFile -}}
# ------------------------------------------
# sourceFile section
# ------------------------------------------
echo
echo -n -e "${__HELP_TITLE_COLOR}SOURCE FILE: ${__RESET_COLOR}"
echo {{ .sourceFile | expandenv | quote }}
{{ end -}}

{{ if .license -}}
# ------------------------------------------
# license section
# ------------------------------------------
echo
echo -n -e "${__HELP_TITLE_COLOR}LICENSE: ${__RESET_COLOR}"
echo {{ .license | quote }}
{{ end -}}

{{ if .copyright -}}
# ------------------------------------------
# copyright section
# ------------------------------------------
{{ if hasSuffix "Callback" .copyright -}}
Array::wrap2 ' ' 76 0 "$({{ .copyright }})"
{{ else -}}
Array::wrap2 ' ' 76 0 {{ .copyright | quote }}
{{ end -}}
{{ end }}
{{end}}
{{end -}}
