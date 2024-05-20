{{ $commandsCount := len .Data.commands }}
{{ if gt $commandsCount 1 }}
local action="$1"
shift || true
case "${action}" in
{{ range $index, $command := .Data.commands }}
{{ $command.commandName -}})
  {{ $command.commandName -}}Parse "$@"
  ;;
{{ end }}
*)
  if Assert::functionExists defaultFacadeAction; then
    # shellcheck disable=SC2016
    defaultFacadeAction "$1" "$@"
  else
    Log::displayError "invalid action requested: ${action}"
    exit 1
  fi
  ;;
esac
exit 0
{{ else }}
{{ (index .Data.commands 0).functionName }}Parse "$@"
{{ end }}
