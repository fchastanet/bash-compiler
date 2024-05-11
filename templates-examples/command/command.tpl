{{define "command"}}
{{ .Data.command.functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" . . false | indent 2}}
}

{{ .Data.command.functionName }}Help() {
  {{- include "command.help" . . false| indent 2 }}
}
{{end}}
