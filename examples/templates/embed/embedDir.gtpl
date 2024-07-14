{{- define "embedDir" -}}
{{ $targetDir := print "${PERSISTENT_TMPDIR:-/tmp}/" .Data.md5sum "/" .Data.asName }}
Compiler::Embed::extractDirFromBase64 \
  {{ $targetDir | quote }} \
  {{ .Data.base64 | quote }}

declare -gx embed_dir_{{ .Data.asName }}={{ $targetDir | quote }}
{{ end }}
