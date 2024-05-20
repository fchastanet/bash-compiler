{{- define "binFile" -}}
#!/usr/bin/env bash

{{ include "binFile.headers.sh" .Data . -}}

# FUNCTIONS

{{- include "commands" .Data.commands . -}}
{{- end -}}
