{{define "command.dataModel.Boolean"}}
onValue: {{ coalesce .Data.onValue 1 }}
offValue: {{ coalesce .Data.offValue "0" }}
min: 0
max: 1
{{end}}
