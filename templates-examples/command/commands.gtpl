{{- define "commands" -}}
{{ $context := . }}
{{ range .Data }}
# ------------------------------------------
# Command {{ .functionName }}
# ------------------------------------------

# options variables initialization
{{ range $index, $option := .options -}}
{{- include "option.init" $option $context -}}
{{ end -}}
# arguments variables initialization
{{ range $index, $arg := .args -}}
{{- include "arg.init" $arg $context -}}
{{ end -}}

# @description parse command options and arguments for {{ .functionName }}
{{ .functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" . $context | indent 2 | trim}} || return $?
  {{ range $callback := .callbacks -}}
  {{ $callback }}
  {{ end }}
}

# @description display command options and arguments help for {{ .functionName }}
{{ .functionName }}Help() {
  {{ include "command.help" . $context | indent 2 | trim}}
}
{{end}}
{{end}}
