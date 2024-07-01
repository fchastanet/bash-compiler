{{- define "checkRequirementsCode" -}}
{{ .Data.functionName -}}() {
{{ range $index, $require := .Data.requires }}
  {{ $functionNameUpper := $require | snakeCase -}}
  if [[ "${REQUIRE_FUNCTION_{{- $functionNameUpper -}}_LOADED:-0}" != 1 ]]; then
    echo >&2 "Requirement {{ $require }} has not been loaded"
    exit 1
  fi
{{ end }}
{{- end }}
