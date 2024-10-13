{{- define "binFile" -}}
#!/usr/bin/env bash
{{ $context := . -}}
###############################################################################
{{ if .Data.binData.commands.default.sourceFile -}}
# GENERATED FROM {{ .Data.binData.commands.default.sourceFile | expandenv }}
{{ else -}}
# GENERATED
{{ end -}}
# DO NOT EDIT IT
# @generated
###############################################################################
# shellcheck disable=SC2288,SC2034

{{ include "binFile.headers.gtpl" .Data.binData $context -}}
{{ include "binFile.initDirs.gtpl" .Data.binData $context -}}

# FUNCTIONS
{{- $sortedDefinitionFiles := .Data.binData.commands.default.definitionFiles | sortByKeys -}}
{{ range $file := $sortedDefinitionFiles }}
{{- includeFileAsTemplate $file $context }}
{{ end }}
{{- include "commands" .Data.binData.commands $context -}}

{{ include "binFile.facade.gtpl" .Data.binData . | trim }}

{{- $mainFunction := .Data.vars.MAIN_FUNCTION_NAME | default "main" }}
MAIN_FUNCTION_NAME="{{ $mainFunction -}}"
{{ $mainFunction -}}() {
{{ include "binFile.hook.main.in.gtpl" . . | trim }}
{{ if .Data.binData.commands.default.mainFile -}}
{{ includeFileAsTemplate .Data.binData.commands.default.mainFile $context | removeFirstShebangLineIfAny | trim }}
{{ end -}}
{{ include "binFile.hook.main.out.gtpl" . . | trim }}
}

# if file is sourced avoid calling main function
# shellcheck disable=SC2178
BASH_SOURCE=".$0" # cannot be changed in bash
# shellcheck disable=SC2128
if test ".$0" == ".${BASH_SOURCE}"; then
  if [[ "${BASH_FRAMEWORK_QUIET_MODE:-0}" = "1" ]]; then
    {{ $mainFunction }} "$@" &>/dev/null
  else
    {{ $mainFunction }} "$@"
  fi
fi
{{- end -}}
