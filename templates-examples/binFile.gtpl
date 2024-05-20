{{- define "binFile" -}}
#!/usr/bin/env bash
###############################################################################
{{ if (index .Data.commands 0).sourceFile -}}
# GENERATED FROM {{ (index .Data.commands 0).sourceFile }}
{{ else -}}
# GENERATED
{{ end -}}
# DO NOT EDIT IT
# @generated
###############################################################################
# shellcheck disable=SC2288,SC2034

{{ include "binFile.headers.sh" .Data . -}}
{{ include "binFile.initDirs.gtpl" .Data . -}}

# FUNCTIONS

{{- include "commands" .Data.commands . -}}

MAIN_FUNCTION_NAME="{{ .Data.vars.MAIN_FUNCTION_NAME -}}"
{{ .Data.vars.MAIN_FUNCTION_NAME -}}() {
  {{ include "binFile.hook.mainEntry.gtpl" .Data . | indent 2 | trim }}
  {{ include "binFile.facade.gtpl" .Data . | indent 2 | trim }}
}

# if file is sourced avoid calling main function
# shellcheck disable=SC2178
BASH_SOURCE=".$0" # cannot be changed in bash
# shellcheck disable=SC2128
test ".$0" != ".${BASH_SOURCE}" || {{ .Data.vars.MAIN_FUNCTION_NAME }} "$@"
{{- end -}}
