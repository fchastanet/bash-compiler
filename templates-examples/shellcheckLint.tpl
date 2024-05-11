{{- $command := fromYAMLFile "templates-examples/shellcheckLint.yaml" -}}
  {{include "command" $command}}
