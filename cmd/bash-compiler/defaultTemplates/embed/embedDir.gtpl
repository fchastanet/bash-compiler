{{- define "embedDir" -}}
{{- $targetDir := print "${PERSISTENT_TMPDIR:-/tmp}/" .Data.md5sum "/" .Data.asName -}}
Linux::requireTarCommand
Compiler::Embed::extractDirFromBase64 \
  {{ $targetDir | quote }} \
  "{{ .Data.base64 | chunkBase64 }}"

declare -gx embed_dir_{{ .Data.asName }}={{ $targetDir | quote }}
{{ end }}
