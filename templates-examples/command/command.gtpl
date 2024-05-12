{{- define "command" -}}
#!/usr/bin/env bash
{{ .Data.functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" .Data . | indent 2}}
}

{{ .Data.functionName }}Help() {
  {{- include "command.help" .Data .| indent 2 }}
}
{{end}}
