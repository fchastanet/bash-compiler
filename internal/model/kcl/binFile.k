import regex
import yaml
import .libs
import file

_configFile = option(key="configFile", type='str', required=True, help="load config file")
_globalOptionsGroups = None
schema BinFileSchema:
  compilerConfig?: CompilerConfigSchema
  binData: BinDataSchema
  vars?: VarsSchema

schema BinDataSchema:
  commands: CommandsSchema
  check:
    commands and len([_attr for _attr, _ in commands]) > 0, "At least one command should be provided"
    all _attr, _ in commands {
      AttrRegexpChecker(_attr, "^[a-zA-Z0-9_]+$",'commands')
    }

schema AnnotationsConfigSchema:
  requireTemplateName: str = "require"
  checkRequirementsTemplateName: str = "checkRequirements"
  embedFileTemplateName: str = "embedFile"
  embedDirTemplateName: str = "embedDir"
  [str]: str
  check:
    regex.match(requireTemplateName, "^[a-zA-Z0-9_]+$"), "annotationsConfig - invalid requireTemplateName ${requireTemplateName}"
    regex.match(checkRequirementsTemplateName, "^[a-zA-Z0-9_]+$"), "annotationsConfig - invalid checkRequirementsTemplateName ${checkRequirementsTemplateName}"
    regex.match(embedFileTemplateName, "^[a-zA-Z0-9_]+$"), "annotationsConfig - invalid embedFileTemplateName ${embedFileTemplateName}"
    regex.match(embedDirTemplateName, "^[a-zA-Z0-9_]+$"), "annotationsConfig - invalid embedDirTemplateName ${embedDirTemplateName}"

schema CompilerConfigSchema:
  rootDir: str
  srcDirs: [str] = ["${rootDir}/src"]
  templateDirs: [str] = ["${rootDir}/template"]
  binDir: str = "${rootDir}/bin"
  targetFile: str
  templateFile: str
  relativeRootDirBasedOnTargetDir: str = "."
  annotationsConfig: AnnotationsConfigSchema = {}
  functionsIgnoreRegexpList: [str] = []

  check:
    isunique(functionsIgnoreRegexpList) if functionsIgnoreRegexpList, "functionsIgnoreRegexpList should contains unique regular expressions"

    len(srcDirs) > 0 if srcDirs, "srcDirs - at least directory one should be provided"
    srcDirs and isunique([_x for _, _x in srcDirs]) if srcDirs, \
      "srcDirs - directories should be unique, check for duplicates"

    len(templateDirs) > 0 if templateDirs, "templateDirs - at least one directory should be provided"
    templateDirs and isunique([_x for _, _x in templateDirs]) if templateDirs, \
      "templateDirs - directories should be unique, check for duplicates"

    all _attr, _ in annotationsConfig {
      AttrRegexpChecker(_attr, "^[a-zA-Z0-9_]+$",'annotationsConfig')
    }

schema AttrRegexpChecker[_attr: str, _regexp: str, _propertyName: str]:
  _result = regex.match(_attr, _regexp)
  check:
    _result == True, "${_propertyName} - invalid attribute ${_attr}"

schema CallbackSchema[callback: str, propertyName: str]:
  _result = libs.assertFunctionNamePriority(callback)
  check:
    _result == True, "${propertyName} - invalid callback ${callback}"

schema CallbacksSchema[callbacks: [str], propertyName: str]:
  _callbacks = callbacks
  check:
    all callback in _callbacks {
      CallbackSchema(callback, propertyName)
    }
    isunique([_callback for _, _callback in callbacks]), \
      "${propertyName}: check for duplicates"

schema PropertyDuplicateSchema[myList: DefinitionFilesSchema, listName: str]:
  _myList = myList
  check:
    isunique([ _x for _, _x in _myList]) \
      if _myList, \
        "${listName} list, check for duplicates"
    isunique([ _x for _x, _ in _myList]) \
      if _myList, \
        "${listName}: list, check for index duplicates"

schema OptionGroupSchema:
  title: str
  [attr: ...str]: any
  check:
    regex.match(title, r"^.+$"), "group title cannot be empty"
    len(attr) == 0, "option group invalid attribute ${attr}, only title is authorized"

schema OptionGroupsSchema:
  [attr: str]: OptionGroupSchema
  check:
    regex.match(attr, r"^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)$"), "invalid group key ${attr}"

schema ValueSchema:
  value: str
  help?: str
  [attr: ...str]: any
  check:
    value, "authorizedValues missing value key"
    len(value) > 0, "authorizedValues item value cannot be empty"
    len(help) if help, "authorizedValues item help cannot be empty if provided"
    len(attr) == 0, "authorizedValues item invalid attribute ${attr}, only value or help are authorized"

schema ParameterSchema:
  parameterType: str
  variableName: str
  type: str = "Boolean"
  functionName: str = "${variableName}Function"
  help: str = ""
  min: int = 0
  max?: int = 0
  if type == "StringArray":
    max: int = -1
  else:
    max: int = 1
  authorizedValues?: [ValueSchema] = None
  callbacks?: [str] = []
  regexp?: str = None
  defaultValue?: str | int = None
  if type == "String":
    defaultValue = ""

  [...str]: any
  check:
    regex.match(variableName, "(^[a-z][A-Za-z_0-9]+$)|(^[A-Z_][A-Z_0-9]+$)"), "Parameter type ${parameterType}: invalid variable name ${variableName}"
    type in ["Boolean", "String", "StringArray"], "Parameter type ${parameterType} - ${variableName}: type '${type}' is unknown"
    libs.assertFunctionName(functionName)

    # min/max checks
    min >= 0, "Parameter type ${parameterType} - ${variableName}: min value ${min} should be greater or equal to 0"
    min <= max if max != -1, "Parameter type ${parameterType} - ${variableName}: min value ${min} should be less or equal to max value ${max}"
    max == -1 or max > 0, "Parameter type ${parameterType} - ${variableName}: max ${max} should be -1 or greater than 0"
    ((min == 0 or min == 1) and (max == 1)) if type == "String", \
      "Parameter type ${parameterType} - ${variableName}: ${type} type, min can only be 0 or 1 when max is 1"
    ((min == 0 or min == 1) and (max == 1)) if type == "Boolean", \
      "Parameter type ${parameterType} - ${variableName}: ${type} type, min can only be 0 or 1"

    len(authorizedValues) == 0 if authorizedValues != None and type == "Boolean", "Parameter type ${parameterType} - ${variableName}: Boolean type cannot have authorizedValues list"
    CallbacksSchema(callbacks, "Parameter type ${parameterType} - ${variableName} - callbacks") if callbacks
    regexp == None if type == "Boolean", "Parameter type ${parameterType} - ${variableName}: You cannot provide a regexp property on a Boolean argument"
    defaultValue == None if type == "StringArray", "Parameter type ${parameterType} - ${variableName}: defaultValue attribute is not supported for type StringArray"

schema OptionSchema(ParameterSchema):
  parameterType: str = "Option"
  alts: [str]
  helpValueName?: str = ""
  if (type == "Boolean"):
    helpValueName = None
  elif (alts and len(alts) > 0 and alts[0] != None):
    helpValueName = alts[0].lstrip("-")

  onValue?: str | int = None
  offValue?: str | int = None
  if type == "Boolean":
    onValue = 1
    offValue = 0
    defaultValue = offValue
  group?: str
  check:
    alts and len(alts) > 0, "option ${variableName} - at least one alt item is required for alt property"
    isunique(alts), "alts should contains unique alt options"

schema ArgumentSchema(ParameterSchema):
  parameterType: str = "Argument"
  name: str
  check:
    type in ["String", "StringArray"], "type '${type}' of parameter ${variableName} is unknown"
    name and len(name) > 0, "argument ${variableName}, please provide a name property"

schema AssertOptionGroupExists[group: str, optionGroups: OptionGroupsSchema, property: str]:
  _optionGroups = optionGroups
  _group = group
  check:
    _group in optionGroups, "${property} - The group ${group} doesn't exists in optionGroups"

schema CommandSchema:
  commandName: str = "default"
  help: str = ""
  longDescription: str = ""
  functionName: str = "${commandName}Function"
  callbacks?: [str] = []
  unknownOptionCallbacks?: [str] = []
  unknownArgumentCallbacks?: [str] = []
  everyArgumentCallbacks?: [str] = []
  beforeParseCallbacks?: [str] = []
  afterParseCallbacks?: [str] = []
  optionGroups?: OptionGroupsSchema
  options?: [OptionSchema] = []
  args?: [ArgumentSchema] = []
  if typeof(options) != "UndefinedType":
    _optionGroups: [str] = [_x.group for _, _x in options]
  [...str]: any
  check:
    regex.match(commandName, r"^[a-zA-Z0-9_-]+$") if commandName, "invalid command name"
    CallbacksSchema(callbacks, "Command ${commandName} - callbacks") if callbacks
    CallbacksSchema(unknownOptionCallbacks, "Command ${commandName} - unknownOptionCallbacks") if unknownOptionCallbacks
    CallbacksSchema(unknownArgumentCallbacks, "Command ${commandName} - unknownArgumentCallbacks") if unknownArgumentCallbacks
    CallbacksSchema(everyArgumentCallbacks, "Command ${commandName} - everyArgumentCallbacks") if everyArgumentCallbacks
    CallbacksSchema(beforeParseCallbacks, "Command ${commandName} - beforeParseCallbacks") if beforeParseCallbacks
    CallbacksSchema(afterParseCallbacks, "Command ${commandName} - afterParseCallbacks") if afterParseCallbacks

    isunique([_x.title for _, _x in optionGroups]) if optionGroups, \
      "Command ${commandName} - optionGroups - title should be unique, check for duplicates"
    libs.assertFunctionName(functionName)
    isunique([_x.name for _, _x in args]) if args, \
      "Command ${commandName} - args - name should be unique, check for duplicates"
    isunique([_x.variableName for _, _x in args]) if args, \
      "Command ${commandName} - args - variableName should be unique, check for duplicates"
    isunique([_x.variableName for _, _x in options]) if options, \
      "Command ${commandName} - options - variableName should be unique, check for duplicates"
    isunique(
      [_x.variableName for _, _x in args] + [_x.variableName for _, _x in options]
    ) if options and args, \
      "Command ${commandName} - variableName should be unique across options and args, check for duplicates"
    isunique(
      [_x.functionName for _, _x in args if args] + [_x.functionName for _, _x in options if options]
    ), \
      "Command ${commandName} - functionName should be unique across options and args, check for duplicates"
    isunique(
      [_alt for _, _x in options for _, _alt in _x.alts]
    ) if options, \
      "Command ${commandName} - alts should be unique across options, check for duplicates"
    all _group in _optionGroups {
      AssertOptionGroupExists(_group, optionGroups, "Command ${commandName}")
    } if typeof(optionGroups) == "OptionGroupsSchema"

schema DefaultCommandSchema(CommandSchema):
  mainFile?: str
  definitionFiles?: DefinitionFilesSchema
  author: str = ""
  sourceFile: str = ""
  license: str = ""
  copyright: str = ""
  copyrightBeginYear: str|int = ""
  version: str = "1.0.0"
  check:
    regex.match(version, r"^([0-9]+\.)?([0-9]+\.)?([0-9]+)$"), \
      "invalid version format, should be x.y.z"
    regex.match(str(copyrightBeginYear), "^[0-9]{4}$") if copyrightBeginYear != "", \
      "copyrightBeginYear should be empty or a valid 4-digits year"
    PropertyDuplicateSchema(definitionFiles, "definitionFiles")

schema DefinitionFilesSchema:
  [order:str]: str
  check:
    typeof(order) == "int" or int(order) > 0, "definitionFiles - must be indexed with int keys %s"
    int(order) > 0 if typeof(order) == "int", "definitionFiles - int key must be >0"

schema CommandsSchema:
  default: DefaultCommandSchema
  [attr: str]: CommandSchema

schema VarsSchema:
  [attr: str]: str
  check:
    regex.match(attr, r'^[A-Z0-9_]+$'), "vars - invalid key ${attr}"

configYaml:BinFileSchema = yaml.decode(file.read(_configFile))

yaml.encode(configYaml, True, True)
