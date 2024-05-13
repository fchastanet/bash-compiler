{{define "dataModel"}}
{{ $context := . -}}
binFile:
  commands:
  {{ range .Data.binFile.commands }}
    -
      functionName: {{ .functionName | errorIfEmpty }}
      commandName: {{ .commandName }}
      version: "{{ .version }}"
      help: |
        {{ .help | indent 8 | trim }}
      helpType: {{ coalesce .helpType "string" }}
      longDescription: |
        {{ coalesce .longDescription "" | indent 8 | trim }}
      {{ if .callbacks -}}
      callbacks:
        {{ range .callbacks -}}
        - {{. -}}
        {{end}}
      {{end}}
      {{- if .unknownOptionCallbacks -}}
      unknownOptionCallbacks:
        {{ range .unknownOptionCallbacks -}}
        - {{. -}}
        {{end}}
      {{end}}
      {{- if .options -}}
      options:
        {{ range .options -}}
        -
          {{ include "dataModel.parameter" . $context | indent 10 | trim }}
          {{ if not .variableType -}}
          {{- logWarn "variable type set to Boolean by default" "variableName" .variableName -}}
          {{- end }}
          {{- $variableType := coalesce .variableType "Boolean" -}}
          variableType: {{ $variableType }}
          {{ if .alts }}
          alts:
            {{- range .alts }}
            - {{. }}
            {{- end -}}
          {{ else -}}
          {{ fail (cat "you must provide alts property for option" .variableName) }}
          {{- end }}
          {{- if eq $variableType "Boolean" -}}
          {{- include "dataModel.option.Boolean" . $context | indent 10 -}}
          {{- else if eq $variableType "String" -}}
          {{- include "dataModel.option.StringCommon" . $context | indent 10 -}}
          {{- include "dataModel.option.String" . $context | indent 10 -}}
          {{- else if eq $variableType "StringArray" -}}
          {{- include "dataModel.option.StringCommon" . $context | indent 10 -}}
          {{- include "dataModel.option.StringArray" . $context | indent 10 -}}
          {{- else }}
          {{ fail (cat "invalid variable type" $variableType " for variable " .variableName) }}
          {{ end }}
        {{ end }}
      {{ end }}
      {{- if .args -}}
      args:
        {{ range .args -}}
        -
          {{ include "dataModel.parameter" . $context | indent 10 | trim }}
          {{ include "dataModel.arg.common" . $context | indent 10 | trim }}
        {{ end }}
      {{ end }}
  {{end}}
  {{end}}