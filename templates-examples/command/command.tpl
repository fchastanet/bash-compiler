{{define "command"}}
{{ .command.functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" . | indent 2}}
}

{{ .command.functionName }}Help() {
  {{- include "command.help" . | indent 2 }}
}
{{end}}
