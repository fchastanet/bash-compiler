package model

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

func TestInvalidFiles(t *testing.T) {
	t.Run("invalidYamlFile", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/invalidYamlFile.yaml")
		assert.ErrorContains(t, err, "expect BinFileSchema, got str")
	})
	t.Run("BinData-missing", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-missing.yaml")
		assert.ErrorContains(t, err, "attribute 'binData' of BinFileSchema is required and can't be None or Undefined")
	})
	t.Run("CompilerConfig-targetFile-missing", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-targetFile-missing.yaml")
		assert.ErrorContains(t, err, "attribute 'targetFile' of CompilerConfigSchema is required and can't be None or Undefined")
	})
	t.Run("CompilerConfig-templateFile-missing", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-templateFile-missing.yaml")
		assert.ErrorContains(t, err,
			"attribute 'templateFile' of CompilerConfigSchema is required and can't be None or Undefined",
		)
	})
	t.Run("CompilerConfig-rootDir-missing", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-rootDir-missing.yaml")
		assert.ErrorContains(t, err, "attribute 'rootDir' of CompilerConfigSchema is required and can't be None or Undefined")
	})
	t.Run("CompilerConfig-srcDirs-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-srcDirs-invalid.yaml")
		assert.ErrorContains(t, err, "srcDirs: [str] = [\"${rootDir}/src\"]\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("CompilerConfig-srcDirs-duplicate", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-srcDirs-duplicate.yaml")
		assert.ErrorContains(t, err, "srcDirs - directories should be unique, check for duplicates")
	})
	t.Run("CompilerConfig-templateDirs-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-templateDirs-invalid.yaml")
		assert.ErrorContains(t, err, "templateDirs: [str] = [\"${rootDir}/template\"]\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("CompilerConfig-templateDirs-duplicate", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-templateDirs-duplicate.yaml")
		assert.ErrorContains(t, err, "templateDirs - directories should be unique, check for duplicates")
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-functionsIgnoreRegexpList-invalid.yaml")
		assert.ErrorContains(t, err, "functionsIgnoreRegexpList: [str] = []\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-duplicate", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-functionsIgnoreRegexpList-duplicate.yaml")
		assert.ErrorContains(t, err, "functionsIgnoreRegexpList should contains unique regular expressions")
	})
	t.Run("CompilerConfig-annotationsConfig-invalidKey", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-invalidKey.yaml")
		assert.ErrorContains(t, err, "annotationsConfig - invalid attribute invàlidKey")
	})
	t.Run("CompilerConfig-annotationsConfig-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-invalidValue.yaml")
		assert.ErrorContains(t, err, "embedDirTemplateName: str = \"embedDir\"\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("CompilerConfig-annotationsConfig-requireTemplateName-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-requireTemplateName-invalid.yaml")
		assert.ErrorContains(t, err, "annotationsConfig - invalid requireTemplateName templateé")
	})
	t.Run("CompilerConfig-annotationsConfig-checkRequirementsTemplateName-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-checkRequirementsTemplateName-invalid.yaml")
		assert.ErrorContains(t, err, "annotationsConfig - invalid checkRequirementsTemplateName invalidé")
	})
	t.Run("CompilerConfig-annotationsConfig-embedFileTemplateName-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-embedFileTemplateName-invalid.yaml")
		assert.ErrorContains(t, err, "annotationsConfig - invalid embedFileTemplateName invalidé")
	})
	t.Run("CompilerConfig-annotationsConfig-embedDirTemplateName-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/CompilerConfig-annotationsConfig-embedDirTemplateName-invalid.yaml")
		assert.ErrorContains(t, err, "annotationsConfig - invalid embedDirTemplateName invalidé")
	})
	t.Run("Vars-invalidKey", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/Vars-invalidKey.yaml")
		assert.ErrorContains(t, err, "vars - invalid key invàlidKey")
	})
	t.Run("Vars-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/Vars-invalidValue.yaml")
		assert.ErrorContains(t, err, "vars?: VarsSchema\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-empty", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-empty.yaml")
		assert.ErrorContains(t, err, "attribute 'binData' of BinFileSchema is required and can't be None or Undefined")
	})
	t.Run("BinData-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-invalid.yaml")
		assert.ErrorContains(t, err, "expect BinDataSchema, got dict")
	})
	t.Run("BinData-commands-empty", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-empty.yaml")
		assert.ErrorContains(t, err, "At least one command should be provided")
	})
	t.Run("BinData-commands-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-invalid.yaml")
		assert.ErrorContains(t, err, "commands - invalid attribute invàlid")
	})
	t.Run("BinData-commands-default-empty", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-empty.yaml")
		assert.ErrorContains(t, err, "attribute 'default' of CommandsSchema is required and can't be None or Undefined")
	})
	t.Run("BinData-commands-default-missing.yaml", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-missing.yaml")
		assert.ErrorContains(t, err, "attribute 'default' of CommandsSchema is required and can't be None or Undefined")
	})

	t.Run("BinData-commands-default-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-invalid.yaml")
		assert.ErrorContains(t, err, "expect DefaultCommandSchema, got str")
	})
	t.Run("BinData-commands-default-mainFile-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-mainFile-invalid.yaml")
		assert.ErrorContains(t, err, "mainFile?: str\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-definitionFiles-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-definitionFiles-invalid.yaml")
		assert.ErrorContains(t, err, "definitionFiles?: DefinitionFilesSchema\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect DefinitionFilesSchema, got str")
	})
	t.Run("BinData-commands-default-definitionFiles-duplicateValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-definitionFiles-duplicateValue.yaml")
		assert.ErrorContains(t, err, "Check failed on the condition: definitionFiles list, check for duplicates")
	})
	t.Run("BinData-commands-default-definitionFiles-invalidKey", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-definitionFiles-invalidKey.yaml")
		assert.ErrorContains(t, err, "definitionFiles - must be indexed with int keys")
	})
	t.Run("BinData-commands-default-definitionFiles-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-definitionFiles-invalidValue.yaml")
		assert.ErrorContains(t, err, "definitionFiles?: DefinitionFilesSchema\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-invalidValue.yaml")
		assert.ErrorContains(t, err, "options?: [OptionSchema] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [OptionSchema], got str")
	})
	t.Run("BinData-commands-default-args-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-args-invalidValue.yaml")
		assert.ErrorContains(t, err, "args?: [ArgumentSchema] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [ArgumentSchema], got str")
	})
	t.Run("BinData-commands-default-author-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-author-invalidValue.yaml")
		assert.ErrorContains(t, err, "author: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-sourceFile-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-sourceFile-invalidValue.yaml")
		assert.ErrorContains(t, err, "sourceFile: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-license-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-license-invalidValue.yaml")
		assert.ErrorContains(t, err, "license: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-copyright-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-copyright-invalidValue.yaml")
		assert.ErrorContains(t, err, "copyright: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-copyrightBeginYear-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-copyrightBeginYear-invalidValue.yaml")
		assert.ErrorContains(t, err, "copyrightBeginYear should be empty or a valid 4-digits year")
	})
	t.Run("BinData-commands-default-optionGroups-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-optionGroups-invalidValue.yaml")
		assert.ErrorContains(t, err, "optionGroups?: OptionGroupsSchema\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect OptionGroupsSchema, got str")
	})
	t.Run("BinData-commands-default-functionName-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-functionName-invalidValue.yaml")
		assert.ErrorContains(t, err, "functionName: str = \"${commandName}Function\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-commandName-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-commandName-invalidValue.yaml")
		assert.ErrorContains(t, err, "commandName: str = \"default\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-version-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-version-invalidValue.yaml")
		assert.ErrorContains(t, err, "Check failed on the condition: invalid version format, should be x.y.z")
	})
	t.Run("BinData-commands-default-help-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-help-invalidValue.yaml")
		assert.ErrorContains(t, err, "help: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	// callbacks
	t.Run("BinData-commands-default-callbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-callbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "callbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-callbacks-invalidCallbackName", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-callbacks-invalidCallbackName.yaml")
		assert.ErrorContains(t, err, "Command default - callbacks - invalid callback invalidé")
	})
	t.Run("BinData-commands-default-callbacks-invalidCallbackPriority", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-callbacks-invalidCallbackPriority.yaml")
		assert.ErrorContains(t, err, "Command default - callbacks - invalid callback invalid@invalidPriority")
	})
	t.Run("BinData-commands-default-unknownOptionCallbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-unknownOptionCallbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "unknownOptionCallbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-unknownArgumentCallbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-unknownArgumentCallbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "unknownArgumentCallbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-everyArgumentCallbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-everyArgumentCallbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "everyArgumentCallbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-beforeParseCallbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-beforeParseCallbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "beforeParseCallbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-afterParseCallbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-afterParseCallbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "afterParseCallbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})

	// optionGroups
	t.Run("BinData-commands-default-optionGroups-default-title-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-optionGroups-default-title-invalidValue.yaml")
		assert.ErrorContains(t, err, "title: str\n\x1b[1;38;5;12m   |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-optionGroups-invalidGroupKey", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-optionGroups-invalidGroupKey.yaml")
		assert.ErrorContains(t, err, "Check failed on the condition: invalid group key invàlidKey")
	})

	// options
	t.Run("BinData-commands-default-options-min-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-min-invalidValue.yaml")
		assert.ErrorContains(t, err, "min: int = 0\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect int, got str")
	})
	t.Run("BinData-commands-default-options-max-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-max-invalidValue.yaml")
		assert.ErrorContains(t, err, "max?: int = 0\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect int, got str")
	})
	t.Run("BinData-commands-default-options-min-gt-max", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-min-gt-max.yaml")
		assert.ErrorContains(t, err, "Parameter type Option - test: min value 2 should be less or equal to max value 1")
	})
	t.Run("BinData-commands-default-options-min-negative", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-min-negative.yaml")
		assert.ErrorContains(t, err, "Parameter type Option - test: min value -1 should be greater or equal to 0")
	})
	t.Run("BinData-commands-default-options-variableName-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-variableName-invalid.yaml")
		assert.ErrorContains(t, err, "Parameter type Option: invalid variable name invàlid")
	})
	t.Run("BinData-commands-default-options-type-invalid", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-type-invalid.yaml")
		assert.ErrorContains(t, err, "Parameter type Option - test: type 'notKnown' is unknown")
	})
	t.Run("BinData-commands-default-options-type-Boolean-invalidMax", func(t *testing.T) {
		err := checkFile(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-type-Boolean-invalidMax.yaml",
		)
		assert.ErrorContains(t, err, "Parameter type Option - var: Boolean type, min can only be 0 or 1")
	})
	t.Run("BinData-commands-default-options-help-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-help-invalidValue.yaml")
		assert.ErrorContains(t, err, "help: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-functionName-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-functionName-invalidValue.yaml")
		assert.ErrorContains(t, err, "functionName: str = \"${variableName}Function\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-callbacks-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-callbacks-invalidValue.yaml")
		assert.ErrorContains(t, err, "callbacks?: [str] = []\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-options-authorizedValues-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-authorizedValues-invalidValue.yaml")
		assert.ErrorContains(t, err, "authorizedValues?: [ValueSchema] = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [ValueSchema], got str")
	})
	t.Run("BinData-commands-default-options-authorizedValues-missingValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-authorizedValues-missingValue.yaml")
		assert.ErrorContains(t, err, "authorizedValues?: [ValueSchema] = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mattribute 'value' of ValueSchema is required and can't be None or Undefined")
	})
	t.Run("BinData-commands-default-options-authorizedValues-invalidField", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-authorizedValues-invalidField.yaml")
		assert.ErrorContains(t, err, "authorizedValues missing value key")
	})
	t.Run("BinData-commands-default-options-regexp-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-regexp-invalidValue.yaml")
		assert.ErrorContains(t, err, "regexp?: str = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-onValue-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-onValue-invalidValue.yaml")
		assert.ErrorContains(t, err, "onValue?: str | int = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str | int, got list")
	})
	t.Run("BinData-commands-default-options-offValue-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-offValue-invalidValue.yaml")
		assert.ErrorContains(t, err, "offValue?: str | int = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str | int, got list")
	})
	t.Run("BinData-commands-default-options-defaultValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-defaultValue.yaml")
		assert.ErrorContains(t, err, "defaultValue?: str | int = None\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str | int, got list")
	})
	t.Run("BinData-commands-default-options-defaultValue-StringArray", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-defaultValue-StringArray.yaml")
		assert.ErrorContains(t, err, "defaultValue attribute is not supported for type StringArray")
	})
	t.Run("BinData-commands-default-options-helpValueName-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-helpValueName-invalidValue.yaml")
		assert.ErrorContains(t, err, "helpValueName?: str = \"\"\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-group-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-group-invalidValue.yaml")
		assert.ErrorContains(t, err, "group?: str\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect str, got list")
	})
	t.Run("BinData-commands-default-options-alts-invalidValue", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-alts-invalidValue.yaml")
		assert.ErrorContains(t, err, "alts: [str]\n\x1b[1;38;5;12m    |\x1b[0m\x1b[1;38;5;9m \x1b[0m \x1b[1;38;5;9mexpect [str], got str")
	})
	t.Run("BinData-commands-default-options-alts-duplicate", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-options-alts-duplicate.yaml")
		assert.ErrorContains(t, err, "Check failed on the condition: alts should contains unique alt options")
	})
	t.Run("BinData-commands-default-optionGroups-group-undefined", func(t *testing.T) {
		err := checkFile(t, "testsData/transformModel-error/BinData-commands-default-optionGroups-group-undefined.yaml")
		assert.ErrorContains(t, err, "The group missingGroup doesn't exists in optionGroups")
	})

	// TODO args
}

func TestMinimalWorkingFile(t *testing.T) {
	t.Run("minimal-configFile", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/minimalWorkingYamlFile.yaml",
			"testsData/transformModel-ok/minimalWorkingYamlFile-expected.yaml",
		)
	})
	t.Run("CompilerConfig-srcDirs-empty", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-srcDirs-empty.yaml",
			"testsData/transformModel-ok/CompilerConfig-srcDirs-empty-expected.yaml",
		)
	})

	t.Run("CompilerConfig-templateDirs-missing", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-templateDirs-missing.yaml",
			"testsData/transformModel-ok/CompilerConfig-templateDirs-missing-expected.yaml",
		)
	})
	t.Run("CompilerConfig-templateDirs-empty", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-templateDirs-empty.yaml",
			"testsData/transformModel-ok/CompilerConfig-templateDirs-empty-expected.yaml",
		)
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-empty", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-functionsIgnoreRegexpList-empty.yaml",
			"testsData/transformModel-ok/CompilerConfig-functionsIgnoreRegexpList-empty-expected.yaml",
		)
	})
	t.Run("CompilerConfig-binDir-invalid", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-binDir-invalid.yaml",
			"testsData/transformModel-ok/CompilerConfig-binDir-invalid-expected.yaml",
		)
	})

	t.Run("CompilerConfig-annotationsConfig-invalid", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-annotationsConfig-invalid.yaml",
			"testsData/transformModel-ok/CompilerConfig-annotationsConfig-invalid-expected.yaml",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-emptyFields", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/CompilerConfig-annotationsConfig-emptyFields.yaml",
			"testsData/transformModel-ok/CompilerConfig-annotationsConfig-emptyFields-expected.yaml",
		)
	})
	t.Run("Vars-empty", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/Vars-empty.yaml",
			"testsData/transformModel-ok/Vars-empty-expected.yaml",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-duplicateKey", func(t *testing.T) {
		AssertFileIsWorking(
			t,
			"testsData/transformModel-ok/BinData-commands-default-definitionFiles-duplicateKey.yaml",
			"testsData/transformModel-ok/BinData-commands-default-definitionFiles-duplicateKey-expected.yaml",
		)
	})
}

func AssertFileIsWorking(t *testing.T, filePath string, expectedFilePath string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	defer errors.SafeCloseDeferCallback(file, &err)
	assert.NilError(t, err)
	var resultWriter bytes.Buffer
	err = transformModel(*file, &resultWriter)
	assert.NilError(t, err)

	expectedFileContent, err := os.ReadFile(expectedFilePath)
	assert.NilError(t, err)
	if diff := cmp.Diff(string(expectedFileContent), resultWriter.String()); diff != "" {
		goldenFile, err := os.OpenFile(expectedFilePath, os.O_WRONLY, os.ModePerm)
		defer errors.SafeCloseDeferCallback(goldenFile, &err)
		goldenFile.Write(resultWriter.Bytes())
		goldenFile.Close()
		fmt.Println(diff)
		t.Errorf("mismatch (-want +got):\n%v", diff)
	}
}

func checkFile(t *testing.T, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	defer errors.SafeCloseDeferCallback(file, &err)
	assert.NilError(t, err)
	var resultWriter bytes.Buffer
	return transformModel(*file, &resultWriter)
}
