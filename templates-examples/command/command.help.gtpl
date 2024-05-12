{{define "command.help"}}
{{ .Data.functionName }}Help() {
  {{- if eq .Data.helpType "function" -}}
  Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
    "$({{ .Data.help }})"
  {{ else }}
  Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
    "{{ .Data.help }}"
  {{ end -}}
  echo

  # ------------------------------------------
  # usage section
  # ------------------------------------------
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" "{{- .Data.commandName }} {{/*
    */}}{{- if .Data.options -}} [OPTIONS]{{- end }} {{/*
    */}}{{- if .Data.arguments -}} [ARGUMENTS]{{- end -}}"

  {{ if .Data.options -}}
  # ------------------------------------------
  # usage/options section
  # ------------------------------------------
  {{ $context := . -}}
  optionsAltList=({{ range $index, $option := .Data.options }}
    "{{ include "option.help.alts" $option $context }}"{{/*
    */}}{{ end }}
  )
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" \
    "{{ .Data.commandName }}" "${optionsAltList[@]}"
  {{ end }}

  {{ if .Data.arguments -}}
  # ------------------------------------------
  # usage/arguments section
  # ------------------------------------------
  echo
  echo -e "${__HELP_TITLE_COLOR}ARGUMENTS:${__RESET_COLOR}"
  {{ $context := . -}}
  {{ range $index, $arg := .Data.args }}
    echo -e "{{ include "arg.help" $arg $context }}"{{/*
  */}}{{ end }}
  {{ end }}

  {{ if .Data.options -}}
  # ------------------------------------------
  # options section
  # ------------------------------------------
  {{ $context := . -}}
  {{ $previousGroupId := "" -}}
  {{ range $index, $option := .Data.options -}}
  {{ $groupId := default "__default" $option.groupId -}}
  {{ if ne $groupId $previousGroupId -}}
  echo
  {{- if eq $groupId "__default" }}
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"{{ end -}}
  {{ else -}}
  echo -e "${__HELP_TITLE_COLOR}{{ .title }}${__RESET_COLOR}"
  {{ if .help }}echo "{{ .help }}"{{ end -}}
  {{ end }}
  echo -e "{{- include "option.help" $option $context -}}"
  {{ $previousGroupId := $groupId }}
  {{ end -}}
  {{ end -}}

  {{ if .Data.longDescription -}}
  # ------------------------------------------
  # longDescription section
  # ------------------------------------------
  {{ if eq .Data.longDescriptionType "function" -}}
  Array::wrap2 ' ' 76 0 "$({{ .Data.longDescription }})"
  {{ else -}}
  Array::wrap2 ' ' 76 0 "$(cat <<EOF
{{ .Data.longDescription }}
EOF
)"{{ end -}}
  {{ end }}

  {{ if .Data.version -}}
  # ------------------------------------------
  # version section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}VERSION: ${__RESET_COLOR}"
  echo '{{ .Data.version }}'
  {{ end -}}

  {{ if .Data.author -}}
  # ------------------------------------------
  # author section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}AUTHOR: ${__RESET_COLOR}"
  echo '{{ .Data.author }}'
  {{ end -}}

  {{ if .Data.sourceFile -}}
  # ------------------------------------------
  # sourceFile section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}SOURCE FILE: ${__RESET_COLOR}"
  echo '{{ .Data.sourceFile }}'
  {{ end -}}

  {{ if .Data.license -}}
  # ------------------------------------------
  # license section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}LICENSE: ${__RESET_COLOR}"
  echo '{{ .Data.license }}'
  {{ end -}}

  {{ if .Data.copyright -}}
  # ------------------------------------------
  # copyright section
  # ------------------------------------------
  {{ if eq .Data.copyrightType "function" }}
  Array::wrap2 ' ' 76 0 "$({{ .Data.copyright }})"
  {{ else -}}
  Array::wrap2 ' ' 76 0 """{{ .Data.copyright }}"""
  {{ end -}}
  {{ end }}
}
{{end}}
