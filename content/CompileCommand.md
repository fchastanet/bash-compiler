---
title: Compile Command
description: Command line reference for the bash-compiler binary - arguments, flags, configuration file and template variables
weight: 20
type: docs
categories: [documentation]
tags: [compile, cli, configuration, templates]
date: '2026-02-14T21:20:53+02:00'
lastmod: '2026-09-17'
version: '2.0'
---

{{% pageinfo %}}

This page is the command line reference. If you are looking for the annotations you write **inside** your bash source
files, see {{% mdlink text="Directives" path="/directives" %}} and {{% mdlink text="Requirements" path="/requirements"
%}} instead. instead.

{{% /pageinfo %}}

<!-- markdownlint-capture -->

<!-- markdownlint-disable MD033 -->

<a name="compileCommandHelp"></a>

<!-- markdownlint-restore -->

## 1. Synopsis

The compiler takes one or more yaml files describing the bash application, interprets the templates and imports the
necessary bash functions.

```text
bash-compiler [<yaml-files> ...] [flags]
```

**Arguments:**

- `[<yaml-files> ...]` the relative or absolute paths of the yaml files describing the binaries to compile. Each file
  must exist or the command fails.

**Flags:**

- `-h, --help` show context-sensitive help and exit.
- `-t, --intermediate-files-dir=<dir>` directory that will contain the generated intermediate files. When this flag is
  not provided, the intermediate files are not saved. Useful to debug the successive compilation passes.
- `--binary-files-extension="-binary.yaml"` the extension used for the automatic search of binary files.
- `-v, --version` print version information and quit.
- `-d, --debug` set the log level to debug.

**Prerequisite:** the current directory must contain a `.bash-compiler` file. It marks the root directory of the project
being compiled, and the compiler exits with an error if it is missing.

## 2. Examples

Compile a single binary from its yaml description:

```bash
cd /path/to/your/project
bash-compiler examples/configReference/shellcheckLint.yaml
```

Keep the intermediate files to inspect what each compilation pass produced:

```bash
bash-compiler examples/configReference/shellcheckLint.yaml \
  --intermediate-files-dir examples/generated --debug
```

Compile several binaries in one call:

```bash
bash-compiler \
  src/_binaries/commandDefinitions/shellcheckLint-binary.yaml \
  src/_binaries/commandDefinitions/awkLint-binary.yaml
```

When running from the sources instead of an installed binary, `go run` requires the root directory to be passed
explicitly through the otherwise hidden `--rootDir` flag:

```bash
go run ./cmd/bash-compiler examples/configReference/shellcheckLint.yaml \
  --rootDir /path/to/bash-tools-framework \
  -t examples/generated -d
```

## 3. Configuration file

Everything that is not a run-time concern lives in the yaml file rather than on the command line. A configuration file
can inherit from others through the `extends` key, which keeps the framework-wide settings in one place:

```yaml
extends:
  - defaultCommand.yaml
  - frameworkConfig.yaml

vars:
  SRC_FILE_PATH: src/_binaries/commandDefinitions/shellcheckLint-binary.yaml

compilerConfig:
  targetFile: ${FRAMEWORK_ROOT_DIR}/bin/shellcheckLint
  relativeRootDirBasedOnTargetDir: ..
  templateFile: binFile.gtpl
```

### 3.1. compilerConfig keys

- `rootDir` the directory used to compute the src file relative path.

- `srcDirs` the list of directories where the framework functions are searched.

  You can add as many entries as needed to define other source dirs. The functions are searched in the order defined,
  which allows function redefinition.

  _Example:_

  ```yaml
  srcDirs:
    - ${FRAMEWORK_ROOT_DIR}/src
    - ${FRAMEWORK_ROOT_DIR}/vendor/bash-tools-framework/src
  ```

  `Functions::myFunction` will be searched in

  - `${FRAMEWORK_ROOT_DIR}/src/Functions/myFunction.sh`
  - `${FRAMEWORK_ROOT_DIR}/vendor/bash-tools-framework/src/Functions/myFunction.sh`

  **Important note:** if you provide your own `srcDirs` and you also need the functions defined by this project, think
  about adding an entry for this project too.

- `binDir` the fallback bin directory used when the `BIN_FILE` directive has not been provided.

- `targetFile` the path of the generated single-file script. It overrides the `BIN_FILE` directive. See {{% mdlink
  text="Directives" path="/directives" %}} for more information about directives.

- `templateDirs` the list of directories from which the go templates are searched. Listing your own directory first
  allows you to override some of the default template includes.

- `templateFile` the entry point template used to render the binary, eg: `binFile.gtpl`.

- `relativeRootDirBasedOnTargetDir` how to reach the root directory from the generated file's directory.

- `functionsIgnoreRegexpList` the list of patterns describing the functions that must be skipped from being imported.

- `annotationsConfig` the templates used to generate the annotation code, eg: `requireTemplate` and
  `checkRequirementsTemplate`. See {{% mdlink text="Requirements" path="/requirements" %}}.

## 4. Template variables

Some variables are automatically generated to be used in your templates:

- `ORIGINAL_TEMPLATE_DIR` allowing you to include the template relative to the script being interpreted
- `TEMPLATE_DIR` the template directory in which you can override the templates defined in `ORIGINAL_TEMPLATE_DIR`

The following variables depend upon the parameters passed to this script:

- `SRC_FILE_PATH` the src file you want to show at the top of generated file to indicate from which source file the
  binary has been generated.
- `SRC_ABSOLUTE_PATH` is the path of the file being compiled, it can be useful if you need to access a path relative to
  this file during compilation.

## 5. Next steps

- {{% mdlink text="Directives" path="/directives" %}} — the annotations to put in your bash source files
- {{% mdlink text="Requirements" path="/requirements" %}} — declaring and ordering runtime preconditions
- {{% mdlink text="How the compiler works" path="/howitworks" %}} — the compilation passes and the diagrams
