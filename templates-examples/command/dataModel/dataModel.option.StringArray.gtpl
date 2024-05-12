{{define "dataModel.option.StringArray"}}
{{ with .Data }}
{{ $min := coalesce .min 0 }}
{{ if .mandatory -}}
{{ $min := 1 }}
{{ end }}
{{ $max := coalesce .max -1 }}
{{ if (ne $max -1) and (gt $min $max) }}
  {{ fail (cat "max value " $min " should be greater than min value " $max) }}
{{ end }}
min: {{ $min }}
max: {{ $max }}
{{ end }}
{{end}}
