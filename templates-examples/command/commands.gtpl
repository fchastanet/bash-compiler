{{- define "commands" -}}
#!/usr/bin/env bash
{{ $context := . }}
{{ range .Data.commands }}
# ------------------------------------------
# Command {{ .functionName }}
# ------------------------------------------
# @description parse command options and arguments for {{ .functionName }}
{{ .functionName }}Parse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  {{include "command.parse" . $context | indent 2 | trim}}
}

{{ if and .longDescription (ne .longDescriptionType "function") -}}
{{ .functionName }}LongDescription="$(cat <<'EOF'
{{ .longDescription }}
EOF
)"
{{ end -}}

# @description display command options and arguments help for {{ .functionName }}
{{ .functionName }}Help() {
  {{ include "command.help" . $context | indent 2 | trim}}
}
{{end}}
{{end}}
