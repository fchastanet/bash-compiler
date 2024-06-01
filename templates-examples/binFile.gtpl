{{- define "binFile" -}}
#!/usr/bin/env bash
{{ $context := . -}}
###############################################################################
{{ if (index .Data.binData.commands 0).sourceFile -}}
# GENERATED FROM {{ (index .Data.binData.commands 0).sourceFile }}
{{ else -}}
# GENERATED
{{ end -}}
# DO NOT EDIT IT
# @generated
###############################################################################
# shellcheck disable=SC2288,SC2034

{{ include "binFile.headers.sh" .Data.binData $context -}}
{{ include "binFile.initDirs.gtpl" .Data.binData $context -}}

# FUNCTIONS
{{ range $file := .Data.binFile.CommandDefinitionFiles }}
{{- includeFileAsTemplate $file $context }}
{{ end }}
{{- include "commands" .Data.binData.commands $context -}}

{{ $mainFunction := .Data.vars.MAIN_FUNCTION_NAME | default "main" -}}
MAIN_FUNCTION_NAME="{{ $mainFunction -}}"
{{ $mainFunction -}}() {
  {{ include "binFile.hook.main.in.gtpl" . . | indent 2 | trim }}
  {{ include "binFile.facade.gtpl" .Data.binData . | indent 2 | trim }}
  {{ include "binFile.hook.main.out.gtpl" . . | indent 2 | trim }}
}

# if file is sourced avoid calling main function
# shellcheck disable=SC2178
BASH_SOURCE=".$0" # cannot be changed in bash
# shellcheck disable=SC2128
test ".$0" != ".${BASH_SOURCE}" || {{ $mainFunction }} "$@"
{{- end -}}
