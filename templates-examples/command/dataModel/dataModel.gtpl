{{define "dataModel"}}
{{ $context := . -}}
command:
  functionName: {{ .Data.command.functionName | errorIfEmpty }}
  commandName: {{ .Data.command.commandName }}
  version: "{{ .Data.command.version }}"
  help: |
    {{ .Data.command.help | indent 4 | trim }}
  helpType: {{ coalesce .Data.command.helpType "string" }}
  longDescription: |
    {{ coalesce .Data.command.longDescription "" | indent 4 | trim }}
  {{ if .Data.command.callbacks -}}
  callbacks:
    {{ range .Data.command.callbacks -}}
    - {{. -}}
    {{end}}
  {{end}}
  {{- if .Data.command.unknownOptionCallbacks -}}
  unknownOptionCallbacks:
    {{ range .Data.command.unknownOptionCallbacks -}}
    - {{. -}}
    {{end}}
  {{end}}
  {{- if .Data.command.options -}}
  options:
    {{ range .Data.command.options -}}
    - variableName: {{ .variableName | errorIfEmpty }}
      {{ if .functionName -}}
      functionName: {{ .functionName -}}
      {{- else -}}
      functionName: {{ .variableName -}}Function
      {{- end }}
      {{ if not .variableType -}}
      {{- logWarn "variable type set to Boolean by default" "variableName" .variableName -}}
      {{- end }}
      {{- $variableType := coalesce .variableType "Boolean" -}}
      variableType: {{ $variableType }}
      help: |
        {{ .help | indent 8 | trim }}
      {{ if .alts -}}
      alts:
        {{- range .alts }}
        - {{. }}
        {{- end -}}
      {{ else -}}
      {{ fail (cat "you must provide alts property for option" .variableName) }}
      {{- end }}
      {{- if eq $variableType "Boolean" -}}
      {{- include "dataModel.option.Boolean" . $context | indent 6 -}}
      {{- else if eq $variableType "String" -}}
      {{- include "dataModel.option.StringCommon" . $context | indent 6 -}}
      {{- include "dataModel.option.String" . $context | indent 6 -}}
      {{- else if eq $variableType "StringArray" -}}
      {{- include "dataModel.option.StringCommon" . $context | indent 6 -}}
      {{- include "dataModel.option.StringArray" . $context | indent 6 -}}
      {{- else }}
      {{ fail (cat "invalid variable type" $variableType " for variable " .variableName) }}
      {{ end }}
    {{ end }}
  {{ end }}
  {{- if .Data.command.args -}}
  args:
    {{ range .Data.command.args -}}
    - variableName: {{ .variableName }}
    {{ end }}
  {{ end }}
{{end}}
