package model

import (
	"bytes"
	"os"
	"testing"

	"github.com/fchastanet/bash-compiler/internal/utils/errors"
	"gotest.tools/v3/assert"
)

func TestInvalidFiles(t *testing.T) {
	t.Run("invalidYamlFile", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/invalidYamlFile.yaml",
			"mismatched types string and struct",
		)
	})
	t.Run("CompilerConfig-targetFile-missing", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-targetFile-missing.yaml",
			"output.compilerConfig.targetFile: incomplete value string",
		)
	})
	t.Run("CompilerConfig-templateFile-missing", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-templateFile-missing.yaml",
			"output.compilerConfig.templateFile: incomplete value string",
		)
	})
	t.Run("CompilerConfig-rootDir-missing", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-rootDir-missing.yaml",
			"output.compilerConfig.rootDir: incomplete value string",
		)
	})
	t.Run("CompilerConfig-srcDirs-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-srcDirs-empty.yaml",
			"output.compilerConfig.srcDirs: 2 errors in empty disjunction",
		)
	})
	t.Run("CompilerConfig-srcDirs-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-srcDirs-invalid.yaml",
			"output.compilerConfig.srcDirs: 2 errors in empty disjunction",
		)
	})
	t.Run("CompilerConfig-srcDirs-duplicate", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-srcDirs-duplicate.yaml",
			"output.compilerConfig.srcDirs: 2 errors in empty disjunction",
		)
	})
	t.Run("CompilerConfig-templateDirs-missing", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-templateDirs-missing.yaml",
			"output.compilerConfig.templateDirs.0: incomplete value string",
		)
	})
	t.Run("CompilerConfig-templateDirs-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-templateDirs-empty.yaml",
			"output.compilerConfig.templateDirs: conflicting values null and list.UniqueItems",
		)
	})
	t.Run("CompilerConfig-templateDirs-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-templateDirs-invalid.yaml",
			"output.compilerConfig.templateDirs: conflicting values \"invalid\" and list.UniqueItems() (mismatched types string and list)",
		)
	})
	t.Run("CompilerConfig-templateDirs-duplicate", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-templateDirs-duplicate.yaml",
			"output.compilerConfig.templateDirs: invalid value [\"dir1\",\"dir1\"] (does not satisfy list.UniqueItems)",
		)
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-functionsIgnoreRegexpList-empty.yaml",
			"output.compilerConfig.functionsIgnoreRegexpList: conflicting values null and list.UniqueItems() (mismatched types null and list)",
		)
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-functionsIgnoreRegexpList-invalid.yaml",
			"output.compilerConfig.functionsIgnoreRegexpList: conflicting values \"invalid\" and list.UniqueItems() (mismatched types string and list)",
		)
	})
	t.Run("CompilerConfig-functionsIgnoreRegexpList-duplicate", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-functionsIgnoreRegexpList-duplicate.yaml",
			"output.compilerConfig.functionsIgnoreRegexpList: invalid value [\"regexp1\",\"regexp1\"] (does not satisfy list.UniqueItems) (and 1 more errors)",
		)
	})
	t.Run("CompilerConfig-binDir-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-binDir-invalid.yaml",
			"output.compilerConfig.binDir: 2 errors in empty disjunction",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-invalid.yaml",
			"output.compilerConfig.annotationsConfig: conflicting values null",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-invalidKey", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-invalidKey.yaml",
			"output.compilerConfig.annotationsConfig.invàlidKey: field not allowed",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-invalidValue.yaml",
			"output.compilerConfig.annotationsConfig.key: conflicting values string and [\"invalidValue\"] (mismatched types string and list)",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-requireTemplateName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-requireTemplateName-invalid.yaml",
			"output.compilerConfig.annotationsConfig.requireTemplateName: conflicting values null and string (mismatched types null and string)",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-requireTemplateName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-requireTemplateName-invalid.yaml",
			"output.compilerConfig.annotationsConfig.requireTemplateName: conflicting values null and string (mismatched types null and string)",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-checkRequirementsTemplateName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-checkRequirementsTemplateName-invalid.yaml",
			"output.compilerConfig.annotationsConfig.checkRequirementsTemplateName: conflicting values null and string (mismatched types null and string)",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-embedFileTemplateName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-embedFileTemplateName-invalid.yaml",
			"output.compilerConfig.annotationsConfig.embedFileTemplateName: conflicting values null and string (mismatched types null and string)",
		)
	})
	t.Run("CompilerConfig-annotationsConfig-embedDirTemplateName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/CompilerConfig-annotationsConfig-embedDirTemplateName-invalid.yaml",
			"output.compilerConfig.annotationsConfig.embedDirTemplateName: conflicting values null and string (mismatched types null and string)",
		)
	})
	t.Run("Vars-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/Vars-empty.yaml",
			"output.vars: conflicting values null",
		)
	})
	t.Run("Vars-invalidKey", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/Vars-invalidKey.yaml",
			"output.vars.invàlidKey: field not allowed",
		)
	})
	t.Run("Vars-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/Vars-invalidValue.yaml",
			"output.vars.KEY: conflicting values string and [\"value\"] (mismatched types string and list)",
		)
	})
	t.Run("BinData-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-empty.yaml",
			"output.binData: conflicting values null and {commands:#CommandsSchema} (mismatched types null and struct)",
		)
	})
	t.Run("BinData-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-invalid.yaml",
			"output.binData.invalid: field not allowed",
		)
	})
	t.Run("BinData-commands-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-empty.yaml",
			"output.binData.commands: conflicting values null and {[\"default\"]:#CommandSchema,[(!=\"default\" & #functionName)]:#CommandSchema} (mismatched types null and struct)",
		)
	})
	t.Run("BinData-commands-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-invalid.yaml",
			"output.binData.commands.invàlid: field not allowed",
		)
	})
	t.Run("BinData-commands-default-empty", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-empty.yaml",
			"output.binData.commands.default: conflicting values null and",
		)
	})
	t.Run("BinData-commands-default-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-invalid.yaml",
			"output.binData.commands.default: conflicting values \"invalid\" and",
		)
	})
	t.Run("BinData-commands-default-mainFile-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-mainFile-invalid.yaml",
			"output.binData.commands.default.mainFile: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-definitionFiles-invalid.yaml",
			"output.binData.commands.default.definitionFiles: conflicting values \"invalid\" and {[=~\"^[0-9]+$\"]:string} (mismatched types string and struct",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-duplicateKey", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-definitionFiles-duplicateKey.yaml",
			"binData.commands.default.definitionFiles.\"1\": conflicting values \"item2\" and \"item1\"",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-duplicateValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-definitionFiles-duplicateValue.yaml",
			"output.binData.commands.default._uniqueDefinitionFileName: conflicting values false and true",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-invalidKey", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-definitionFiles-invalidKey.yaml",
			"output.binData.commands.default.definitionFiles.invalidKey: field not allowed",
		)
	})
	t.Run("BinData-commands-default-definitionFiles-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-definitionFiles-invalidValue.yaml",
			"output.binData.commands.default.definitionFiles.\"1\": conflicting values string and [\"invalid\"] (mismatched types string and list)",
		)
	})
	t.Run("BinData-commands-default-options-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-invalidValue.yaml",
			"output.binData.commands.default.options: conflicting values \"invalid\" and [...#OptionSchema]",
		)
	})
	t.Run("BinData-commands-default-args-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-args-invalidValue.yaml",
			"output.binData.commands.default.args: conflicting values \"invalid\" and [...#ArgSchema] (mismatched types string and list)",
		)
	})
	t.Run("BinData-commands-default-author-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-author-invalidValue.yaml",
			"output.binData.commands.default.author: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-sourceFile-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-sourceFile-invalidValue.yaml",
			"output.binData.commands.default.sourceFile: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-license-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-license-invalidValue.yaml",
			"output.binData.commands.default.license: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-copyright-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-copyright-invalidValue.yaml",
			"output.binData.commands.default.copyright: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-copyrightBeginYear-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-copyrightBeginYear-invalidValue.yaml",
			"output.binData.commands.default.copyrightBeginYear: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-optionGroups-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-optionGroups-invalidValue.yaml",
			"output.binData.commands.default.optionGroups: conflicting values \"invalid\" and {[(=~\"^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)$\" & !~\"^()$\")]:{title:(=~\"^.*$\"|*\"\")}} (mismatched types string and struct)",
		)
	})
	t.Run("BinData-commands-default-functionName-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-functionName-invalidValue.yaml",
			"output.binData.commands.default.functionName: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-commandName-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-commandName-invalidValue.yaml",
			"output.binData.commands.default.commandName: conflicting values =~\"^[a-zA-Z0-9_-]+$\" and [\"invalid\"] (mismatched types string and list",
		)
	})
	t.Run("BinData-commands-default-version-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-version-invalidValue.yaml",
			"output.binData.commands.default.version: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-help-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-help-invalidValue.yaml",
			"output.binData.commands.default.help: 2 errors in empty disjunction",
		)
	})
	// callbacks
	t.Run("BinData-commands-default-callbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-callbacks-invalidValue.yaml",
			"output.binData.commands.default.callbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-unknownOptionCallbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-unknownOptionCallbacks-invalidValue.yaml",
			"output.binData.commands.default.unknownOptionCallbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-unknownArgumentCallbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-unknownArgumentCallbacks-invalidValue.yaml",
			"output.binData.commands.default.unknownArgumentCallbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-everyArgumentCallbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-everyArgumentCallbacks-invalidValue.yaml",
			"output.binData.commands.default.everyArgumentCallbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-beforeParseCallbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-beforeParseCallbacks-invalidValue.yaml",
			"output.binData.commands.default.beforeParseCallbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-afterParseCallbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-afterParseCallbacks-invalidValue.yaml",
			"output.binData.commands.default.afterParseCallbacks: 2 errors in empty disjunction",
		)
	})

	// optionGroups
	t.Run("BinData-commands-default-optionGroups-default-title-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-optionGroups-default-title-invalidValue.yaml",
			"output.binData.commands.default.optionGroups.default.title: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-optionGroups-invalidGroupKey", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-optionGroups-invalidGroupKey.yaml",
			"output.binData.commands.default.optionGroups.invàlidKey: field not allowed",
		)
	})

	// options
	t.Run("BinData-commands-default-options-min-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-min-invalidValue.yaml",
			"output.binData.commands.default.options.0.min: 2 errors in empty disjunction:",
		)
	})
	t.Run("BinData-commands-default-options-max-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-max-invalidValue.yaml",
			"output.binData.commands.default.options.0.max: 4 errors in empty disjunction:",
		)
	})
	t.Run("BinData-commands-default-options-min-gt-max", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-min-gt-max.yaml",
			"output.binData.commands.default.options.0.min: 3 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-min-negative", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-min-negative.yaml",
			"output.binData.commands.default.options.0.min: 3 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-variableName-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-variableName-invalid.yaml",
			"output.binData.commands.default.options.0.variableName: invalid value \"invàlid\" (out of bound =~\"(^[a-z][A-Za-z_0-9]+$)|(^[A-Z_][A-Z_0-9]+$)\")",
		)
	})
	t.Run("BinData-commands-default-options-type-invalid", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-type-invalid.yaml",
			"output.binData.commands.default.options.0.type: 3 errors in empty disjunction",
		)
	})
	// not working
	// t.Run("BinData-commands-default-options-type-Boolean-invalidMax", func(t *testing.T) {
	// 	checkFileError(
	// 		t,
	// 		"testsData/transformModel-error/BinData-commands-default-options-type-Boolean-invalidMax.yaml",
	// 		"output.binData.commands.default.options.0.type: 3 errors in empty disjunction",
	// 	)
	// })
	t.Run("BinData-commands-default-options-help-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-help-invalidValue.yaml",
			"output.binData.commands.default.options.0.help: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-functionName-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-functionName-invalidValue.yaml",
			"output.binData.commands.default.options.0.functionName: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-callbacks-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-callbacks-invalidValue.yaml",
			"output.binData.commands.default.options.0.callbacks: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-authorizedValues-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-authorizedValues-invalidValue.yaml",
			"output.binData.commands.default.options.0.authorizedValues: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-authorizedValues-missingValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-authorizedValues-missingValue.yaml",
			"output.binData.commands.default.options.0.authorizedValues.0.value: incomplete value strings.MinRunes(1)",
		)
	})
	t.Run("BinData-commands-default-options-authorizedValues-invalidField", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-authorizedValues-invalidField.yaml",
			"output.binData.commands.default.options.0.authorizedValues: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-regexp-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-regexp-invalidValue.yaml",
			"output.binData.commands.default.options.0.regexp: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-onValue-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-onValue-invalidValue.yaml",
			"output.binData.commands.default.options.0.onValue: 4 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-offValue-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-offValue-invalidValue.yaml",
			"output.binData.commands.default.options.0.offValue: 4 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-defaultValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-defaultValue.yaml",
			"output.binData.commands.default.options.0.defaultValue: 4 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-helpValueName-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-helpValueName-invalidValue.yaml",
			"output.binData.commands.default.options.0.helpValueName: 2 errors in empty disjunction",
		)
	})
	t.Run("BinData-commands-default-options-group-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-group-invalidValue.yaml",
			"output.binData.commands.default.options.0.group: conflicting values string and [\"invalid\"] (mismatched types string and list",
		)
	})
	t.Run("BinData-commands-default-options-alts-invalidValue", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-alts-invalidValue.yaml",
			"output.binData.commands.default.options.0.alts: conflicting values \"invalid\" and list.UniqueItems() (mismatched types string and list",
		)
	})
	t.Run("BinData-commands-default-options-alts-duplicate", func(t *testing.T) {
		checkFileError(
			t,
			"testsData/transformModel-error/BinData-commands-default-options-alts-duplicate.yaml",
			"output.binData.commands.default.options.0.alts: invalid value [\"-v\",\"-v\"] (does not satisfy list.UniqueItems",
		)
	})

	// TODO args
}

func TestMinimalWorkingFile(t *testing.T) {
	file, err := os.OpenFile("testsData/transformModel-ok/minimalWorkingYamlFile.yaml", os.O_RDONLY, os.ModePerm)
	defer errors.SafeCloseDeferCallback(file, &err)
	assert.NilError(t, err)
	var resultWriter bytes.Buffer
	err = transformModel(*file, &resultWriter)
	assert.NilError(t, err)
}

func checkFileError(t *testing.T, fileName string, expectedError string) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	defer errors.SafeCloseDeferCallback(file, &err)
	assert.NilError(t, err)
	var resultWriter bytes.Buffer
	err = transformModel(*file, &resultWriter)
	assert.ErrorContains(t, err, expectedError)
}
