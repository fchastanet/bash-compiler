{{- define "binFile" -}}
#!/usr/bin/env bash

# FUNCTIONS

{{- include "commands" .Data.commands . -}}
{{- end -}}
