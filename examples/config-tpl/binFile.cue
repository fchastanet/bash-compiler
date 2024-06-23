// validation command:
//   cue vet    -l input: examples/config/shellcheckLint-generated.yaml examples/config/binFile.cue -E -c --trace --strict
// export command:
//   cue export -l input: examples/config/shellcheckLint-generated.yaml examples/config/binFile.cue --out yaml -e output
package config

import "list"
import "strings"

// place the yaml input here with "-l"
input: _

// validate the input against a schema
input: #Schema
#CompileConfigSchema: {
  FRAMEWORK_ROOT_DIR: string | "${FRAMEWORK_ROOT_DIR}"
  FRAMEWORK_SRC_DIR: string | "${FRAMEWORK_SRC_DIR:-${FRAMEWORK_ROOT_DIR}/src}"
  FRAMEWORK_SRC_DIRS: list.UniqueItems() & [string, ...string] | *[ "${FRAMEWORK_ROOT_DIR}/src" ]
  FRAMEWORK_VENDOR_BIN_DIR: string | "${FRAMEWORK_VENDOR_BIN_DIR:-${FRAMEWORK_ROOT_DIR}/vendor/bin}"
  FRAMEWORK_VENDOR_DIR: string | "${FRAMEWORK_VENDOR_DIR:-${FRAMEWORK_ROOT_DIR}/vendor}"
  BASH_FRAMEWORK_DISPLAY_LEVEL: string | "${BASH_FRAMEWORK_DISPLAY_LEVEL:-3}"
  BASH_FRAMEWORK_LOG_FILE: string | "${BASH_FRAMEWORK_LOG_FILE:-${FRAMEWORK_ROOT_DIR}/logs/${0##*/}.log}"
  BASH_FRAMEWORK_LOG_FILE_MAX_ROTATION: string | "${BASH_FRAMEWORK_LOG_FILE_MAX_ROTATION:-5}"
  BASH_FRAMEWORK_LOG_LEVEL: string | "${BASH_FRAMEWORK_LOG_LEVEL:-0}"
  BASH_FRAMEWORK_THEME: string | "${BASH_FRAMEWORK_THEME:-default}"
  BATS_FILE_NOT_NEEDED_REGEXP: list.UniqueItems() & [...string]
  COMPILE_PARAMETERS: {
    binDir: string | "${FRAMEWORK_BIN_DIR}"
    rootDir: string | "${FRAMEWORK_ROOT_DIR}"
    srcDir: string | "${FRAMEWORK_SRC_DIR}"
    templateDir: string | "${FRAMEWORK_SRC_DIR}"
  }
  DISPLAY_DURATION: string | "${DISPLAY_DURATION:-1}"
  FRAMEWORK_BIN_DIR: string | "${FRAMEWORK_BIN_DIR:-${FRAMEWORK_ROOT_DIR}/bin}"
  FRAMEWORK_FILES_FUNCTION_MATCHING_IGNORE_REGEXP: list.UniqueItems() & [...string]
  FRAMEWORK_FUNCTIONS_IGNORE_REGEXP: list.UniqueItems() & [...string]
  NON_FRAMEWORK_FILES_REGEXP: list.UniqueItems() & [...string]
}

#BinFileSchema: {
  targetFile:                      string
  relativeRootDirBasedOnTargetDir: string | *"."
  templateFile:                    string
  templateName:                    string
  templateDirs: list.UniqueItems() & [string, ...string]
  srcDirs: list.UniqueItems() & [string, ...string]
}

#defaultValueSchema: null | int | string

#functionName: =~"^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)$"

#callbacks: list.UniqueItems() & [...#functionName] | *[]

#ArgSchema: {
  // Variable Name
  variableName: =~"(^[a-z][A-Za-z_0-9]+$)|(^[A-Z_][A-Z_0-9]+$)"

  // Type
  type: "String" | "StringArray"

  // Help
  help!:         string | *""

  functionName: #functionName | *"\(variableName)Function"

  // Authorized values (if null ignore)
  authorizedValues: *([
    if (authorizedValuesList != _|_) {
      for opt in authorizedValuesList {
        {
          value: "\(opt)"
          help: "\(opt)"
        }
      }
    }
  ]) | null | [...{
    value: strings.MinRunes(1)
    help:  strings.MinRunes(0)
  }]

  authorizedValuesList: *[] | null | [...=~"^.+$"]

  // Default Value of the argument (if null default value following
  // the type, eg: [] for StringArray)
  defaultValue: *([
    if (type == "StringArray") { [] }
    ""
  ][0]) | #defaultValueSchema

  // Min
  min!: int & >0 | *0

  // Max
  max!: int & (>0| -1) | *1
  _checkMaxValue: max == -1 | max>=min

  callbacks: #callbacks

  regexp: *([
    if type == "Boolean" {null}
    ""
  ][0]) | null | string

  // Name
  name: =~"^.+$"
}

#OptionSchema: {
  #minSchema: int & >=0 & (max != -1 | <=max) | *0
  #maxSchema: int | (type == "StringArray" | *-1) | *1
  _checkMaxValue: max == -1 | max>=min

  // Variable Name
  variableName: =~"(^[a-z][A-Za-z_0-9]+$)|(^[A-Z_][A-Z_0-9]+$)"

  // Type
  type: "Boolean" | "String" | "StringArray" | *"Boolean"
  _checkMinMaxCompatibleWithBoolean: type == "Boolean" && (min == 0 || min == 1) && (max == 1)
  _checkMinMaxCompatibleWithString: type == "String" && (min == 0 || min == 1) && (max == 1)
  _checkMinMaxCompatibleWithStringArray: type == "StringArray" && (min >= 0) && (max >= 0)

  // Help
  help:         string | *""

  functionName: #functionName | *"\(variableName)Function"

  // Min
  min: #minSchema

  // Max
  max: #maxSchema

  callbacks: #callbacks

  // Authorized values (if null ignore)
  authorizedValues: *([
    if (authorizedValuesList != _|_) for opt in authorizedValuesList {
      {
        value: "\(opt)"
        help: "\(opt)"
      }
    }
  ]) | null | [...{
    value: strings.MinRunes(1)
    help:  strings.MinRunes(0)
  }]

  authorizedValuesList: *[] | null | [...=~"^.+$"]

  regexp: *([
    if (type == "Boolean") {null}
    ""
  ][0]) | null | string

  // OnValue
  onValue: *([
    if (type == "Boolean") {1}
    if (type != "Boolean") {null}
  ][0]) | null | int | string

  // OffValue
  offValue: *([
    if (type == "Boolean") {0}
    if (type != "Boolean") {null}
  ][0]) | null | int | string
  // Default Value of the option (if null default value following
  // the type, eg: [] for StringArray)
  defaultValue: *([
    if (type == "Boolean") { offValue }
    if (type == "StringArray") { [] }
    ""
  ][0]) | #defaultValueSchema

  // HelpValueName
  helpValueName: *([
    if (type == "Boolean") {null}
    if alts[0] != _|_ {strings.TrimLeft(alts[0], "-")}
    ""
  ][0]) | null | string

  // Group
  group!: string

  // Alts
  alts: list.UniqueItems() & [=~"^.+$", ...=~"^.+$"]
}

#CommandSchema: {
  // command definition files to include
  definitionFiles: list.UniqueItems() & [string, ...string]
  // Options
  options: [...#OptionSchema]

  // Args
  args: [...#ArgSchema]

  // check for unique option or arg function name
  _uniqueFunctionName: true
  _uniqueFunctionName: list.UniqueItems([
    for _, i in options {
      if i.functionName != _|_ {i.functionName}
    }
    for _, j in args {
      if j.functionName != _|_ {j.functionName}
    }
    functionName
  ]) // TODO when feature error ready https://github.com/cue-lang/cue/issues/943

  _uniqueVariableName: true
  _uniqueVariableName: list.UniqueItems([
    for _, i in options {i.variableName},
    for _, j in args {j.variableName}
  ]) // TODO when feature error ready https://github.com/cue-lang/cue/issues/943

  // Author
  author: string | *""

  // Sourcefile
  sourceFile: string | *""

  // License
  license: string | *""

  // Copyright
  copyright: string | *""

  // Option Groups
  optionGroups: {
    default: {
      title: =~"^.+$" | *"Default group"
    }
    {[=~"^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)$" & !~"^()$"]: {
      // Title
      title: =~"^.*$" | *""
    }}

    {[=~"examples" & !~"^()$"]: _}
  }
  functionName: #functionName | *"\(commandName)Function"

  // Command Name
  commandName: =~"^[a-zA-Z0-9_-]+$"

  // Version
  version: =~ "^([0-9]+\\.)?([0-9]+\\.)?([0-9]+)$" | *"1.0.0"

  // Help
  help: string | *""

  // Long Description
  longDescription:        string | *""
  callbacks:              #callbacks
  unknownOptionCallbacks: #callbacks
}

#CommandsSchema: {
  [CommandName= "default"]: #CommandSchema
  [CommandName= !="default" & #functionName]?: #CommandSchema
}

#Schema: {
	// Commands
	binFile: #BinFileSchema
	vars: {
		{[=~"^[A-Z0-9_]+$" & !~"^()$"]: string}
	}
	binData: {
		commands: #CommandsSchema
	}

  compileConfig: #CompileConfigSchema
}

output: input
