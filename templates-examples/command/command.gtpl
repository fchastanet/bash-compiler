{{- define "command" -}}
#!/usr/bin/env bash
{{ .Data.functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" .Data . | indent 2 | trim}}
}

{{ if and .Data.longDescription (ne .Data.longDescriptionType "function") -}}
{{ .Data.functionName }}LongDescription="$(cat <<'EOF'
{{ .Data.longDescription }}
EOF
)"
{{ end -}}


{{ .Data.functionName }}Help() {
  {{ include "command.help" .Data .| indent 2 | trim}}
}
{{end}}
