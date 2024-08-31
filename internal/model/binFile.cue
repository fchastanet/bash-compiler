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
#CompilerConfigSchema: {
  targetFile:                      string
  relativeRootDirBasedOnTargetDir: string | *"."
  templateFile: string | input.compilerConfig.rootDir
  annotationsConfig: {
    requireTemplateName: string | *"require"
    checkRequirementsTemplateName: string | *"checkRequirements"
    embedFileTemplateName: string | *"embedFile"
    embedDirTemplateName: string | *"embedDir"
    [=~"^[a-zA-Z0-9_]+$" & !~"^()$"]: string
  }
  rootDir: string
  srcDirs: list.UniqueItems() & [string, ...string] | *[ "\(rootDir)/src" ]
  binDir: string | "\(rootDir)/bin"
  templateDirs: list.UniqueItems() & [string, ...string]

  functionsIgnoreRegexpList: list.UniqueItems() & [...string]
}

#defaultValueSchema: null | int | string

#functionName: =~"^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)(@[0-9]+)?$"

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
  mainFile: string | *null
  // command definition files to include
  definitionFiles:
    [=~"^[a-zA-Z0-9_]+$" & !~"^()$"]: string
  // check for unique definitionFile
  _uniqueDefinitionFileName: true
  _uniqueDefinitionFileName: list.UniqueItems([
    for _, i in definitionFiles {
      if i != _|_ {i}
    }
  ]) // TODO when feature error ready https://github.com/cue-lang/cue/issues/943

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
  unknownArgumentCallbacks: #callbacks
  beforeParseCallbacks: #callbacks
  afterParseCallbacks: #callbacks
}

#CommandsSchema: {
  [CommandName= "default"]: #CommandSchema
  [CommandName= !="default" & #functionName]?: #CommandSchema
}

#Schema: {
	// Commands
  compilerConfig: #CompilerConfigSchema
	vars: {
		{[=~"^[A-Z0-9_]+$" & !~"^()$"]: string}
	}
	binData: {
		commands: #CommandsSchema
	}

}

output: input
