{{- define "embedFile" -}}
{{- $targetFile := print "${PERSISTENT_TMPDIR:-/tmp}/" .Data.md5sum "/" .Data.asName -}}
Linux::requireTarCommand
Compiler::Embed::extractFileFromBase64 \
  {{ $targetFile | quote }} \
  "{{ .Data.base64 | chunkBase64 }}" \
  {{ .Data.fileMode | quote }}

declare -gx embed_file_{{ .Data.asName }}={{ $targetFile | quote }}
{{ end }}
