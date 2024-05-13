{{define "dataModel"}}
{{ $context := . -}}
{{ with .Data.binFile }}
binFile:
  {{ if .optionGroups -}}
  optionGroups:
    {{ range .optionGroups }}
    {{ if index . "include" }}
    {{ $data := fromYAMLFile .include }}
    {{ $data.options.optionGroups | toYAML | indent 2 | trim }}
    {{ else }}
    -
      title: "{{ .title }}"
      functionName: "{{ .functionName }}"
    {{ end }}
    {{ end }}
  {{ end }}{{/* end if .optionsGroups */}}
  commands:
  {{ range .commands }}
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
        {{ if index . "include" }}
        {{ $data := fromYAMLFile .include }}
        {{ $data.options.options | toYAML | indent 6 | trim }}
        {{ else }}
        -
          {{ include "dataModel.parameter" . $context | indent 10 | trim }}
          {{ if not .type -}}
          {{- logWarn "variable type set to Boolean by default" "variableName" .variableName -}}
          {{- end }}
          {{- $type := coalesce .type "Boolean" -}}
          type: {{ $type }}
          {{ if .alts }}
          alts:
            {{- range .alts }}
            - {{. }}
            {{- end -}}
          {{ else -}}
          {{ fail (cat "you must provide alts property for option" .variableName) }}
          {{- end }}
          {{- if eq $type "Boolean" -}}
          {{- include "dataModel.option.Boolean" . $context | indent 10 -}}
          {{- else if eq $type "String" -}}
          {{- include "dataModel.option.StringCommon" . $context | indent 10 -}}
          {{- include "dataModel.option.String" . $context | indent 10 -}}
          {{- else if eq $type "StringArray" -}}
          {{- include "dataModel.option.StringCommon" . $context | indent 10 -}}
          {{- include "dataModel.option.StringArray" . $context | indent 10 -}}
          {{- else }}
          {{ fail (cat "invalid variable type" $type " for variable " .variableName) }}
          {{ end }}
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
  {{end}}