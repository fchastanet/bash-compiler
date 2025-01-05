# Bash Compiler

[![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/SchwarzIT/go-template)

> **_TIP:_** Checkout related projects of this suite
>
> - [My documents](https://fchastanet.github.io/my-documents/)
> - [Bash Tools Framework](https://fchastanet.github.io/bash-tools-framework/)
> - [Bash Tools](https://fchastanet.github.io/bash-tools/)
> - [Bash Dev Env](https://fchastanet.github.io/bash-dev-env/)
> - **[Bash Compiler](https://fchastanet.github.io/bash-compiler/)**

<!-- markdownlint-capture -->

<!-- markdownlint-disable MD013 -->

[![GitHub release (latest SemVer)](https://img.shields.io/github/release/fchastanet/bash-compiler?logo=github&sort=semver)](https://github.com/fchastanet/bash-compiler/releases)
[![GitHubLicense](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/fchastanet/bash-compiler/blob/master/LICENSE)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
[![CI/CD](https://github.com/fchastanet/bash-compiler/actions/workflows/main.yml/badge.svg)](https://github.com/fchastanet/bash-compiler/actions?query=workflow%3A%22Lint+and+test%22+branch%3Amaster)
[![ProjectStatus](http://opensource.box.com/badges/active.svg)](http://opensource.box.com/badges "Project Status")
[![DeepSource](https://deepsource.io/gh/fchastanet/bash-compiler.svg/?label=active+issues&show_trend=true)](https://deepsource.io/gh/fchastanet/bash-compiler/?ref=repository-badge)
[![DeepSource](https://deepsource.io/gh/fchastanet/bash-compiler.svg/?label=resolved+issues&show_trend=true)](https://deepsource.io/gh/fchastanet/bash-compiler/?ref=repository-badge)
[![AverageTimeToResolveAnIssue](http://isitmaintained.com/badge/resolution/fchastanet/bash-compiler.svg)](http://isitmaintained.com/project/fchastanet/bash-compiler "Average time to resolve an issue")
[![PercentageOfIssuesStillOpen](http://isitmaintained.com/badge/open/fchastanet/bash-compiler.svg)](http://isitmaintained.com/project/fchastanet/bash-compiler "Percentage of issues still open")

<!-- markdownlint-restore -->

- [1. Excerpt](#1-excerpt)
- [2. Documentation](#2-documentation)
  - [2.1. Go Libraries used](#21-go-libraries-used)
  - [2.2. Template system](#22-template-system)
  - [2.3. Compiler](#23-compiler)
- [3. Development](#3-development)
  - [3.1. Pre-commit hook](#31-pre-commit-hook)
    - [3.1.1. @embed](#311-embed)
  - [3.2. Build/run/clean](#32-buildrunclean)
    - [3.2.1. Build](#321-build)
    - [3.2.2. Tests](#322-tests)
    - [3.2.3. Coverage](#323-coverage)
    - [3.2.4. run the binary](#324-run-the-binary)
    - [3.2.5. Clean](#325-clean)
- [4. Commands](#4-commands)
- [5. KCL](#5-kcl)
- [6. Alternatives](#6-alternatives)

## 1. Excerpt

This tool allows to detect all the framework functions used inside a given sh
file. The framework functions matches the pattern Namespace::functionName (we
can have several namespaces separated by the characters ::). These framework
functions will be injected inside a compiled file. The process is recursive so
that every framework functions used by imported framework functions will be
imported as well (of course only once).

## 2. Documentation

### 2.1. Go Libraries used

- [slog](https://pkg.go.dev/golang.org/x/exp/slog) is logging system
  - [slog tutorial](https://betterstack.com/community/guides/logging/logging-in-go/#customizing-the-default-logger)
- [Yaml parser](https://github.com/goccy/go-yaml) is used to load template data
  from yaml file
- [Kong](https://github.com/alecthomas/kong) used for command arguments parsing
- [cuelang](https://github.com/cue-lang/cue) allows to transform yaml file in
  another one

### 2.2. Template system

[template system](https://pkg.go.dev/text/template@go1.22.3)
[doc 1](https://lets-go.alexedwards.net/sample/02.08-html-templating-and-inheritance.html)

There is the choice between Go template/text or template/html libraries. I
chosen template/text to avoid some escaping that are not needed in bash.

Go template/text or template/html don't provide any execution context to the
filters (FuncMap).

I'm not using Template.ParseGlob because I have to call it twice to include
files of root directory and sub directories with 2 glob patterns. But a bug in
text/template makes the template be initialized again after each calls to
ParseGlob function. So I compute manually list of templates in
internal/render/render.go NewTemplate function.

I simulated a context by pushing the context to the render function. So the data
associated to the template has the following structure:

```go
type Context struct {
 Template *template.Template
 Name     string
 RootData any
 Data     any
}
```

- Template points to the first template that has been rendered
- Name is the name of the first template that has been rendered
- RootData are the data that have been sent at the start of the rendering
- Data are the data sent to the sub template (possibly a part of RootData or the
  whole RootData)

Then each filter has to be called with the right context. The special filter
`include` allows to include a sub template overriding context Data.

Template filter functions, `internal/render/functions/index.go` includes:

- [Sprig filter functions](https://github.com/Masterminds/sprig)
  - Sprig is not maintained anymore, a possible alternate fork is
    [sprout](https://github.com/go-sprout/sprout) but it misses a lot of
    functions.
- my own templates functions
  - string functions
    - stringLength
    - format allow to format string like in this example
      - `{{ format "${%sLongDescription[@]}" .functionName }}`
  - templates functions
    - include: allows to include a template by template name allowing to use
      filter
    - includeFile: allows to include a template by filename
    - includeFileAsTemplate: same as includeFile but interpreting the file as a
      template
    - dynamicFile: resolve first matching filepath in paths provided as argument

### 2.3. Compiler

see [Compile command](CompileCommand.md).

## 3. Development

### 3.1. Pre-commit hook

This repository uses pre-commit software to ensure every commits respects a set
of rules specified by the `.pre-commit-config.yaml` file. It supposes pre-commit
software is [installed](https://pre-commit.com/#install) in your environment.

You also have to execute the following command to enable it:

```bash
pre-commit install --hook-type pre-commit --hook-type pre-push
```

Now each time you commit or push, some linters/compilation tools are launched
automatically

#### 3.1.1. @embed

Allows to embed files, directories or a framework function. The following syntax
can be used:

_Syntax:_ `# @embed "srcFile" AS "targetFile"`

_Syntax:_ `# @embed "srcDir" AS "targetDir"`

if `@embed` annotation is provided, the file/dir provided will be added inside
the resulting bin file as a tar gz file(base64 encoded) and automatically
extracted when executed.

The compiler's embed annotation offers the ability to embed files or
directories. `annotationEmbed` allows to:

- **include a file**(binary or not) as base64 encoded, the file can then be
  extracted using the automatically generated method
  `Compiler::Embed::extractFile_asName` where asName is the name chosen using
  annotation explained above. The original file mode will be restored after
  extraction. The variable `embed_file_asName` contains the targeted filepath.
- **include a directory**, the directory will be tar gz and added to the
  compiled file as base64 encoded string. The directory can then be extracted
  using the automatically generated method `Compiler::Embed::extractDir_asName`
  where asName is the name chosen using annotation explained above. The variable
  embed_dir_asName contains the targeted directory path.
- **include a bash framework function**, a special binary file that simply calls
  this function will be automatically generated. This binary file will be added
  to the compiled file as base64 encoded string. Then it will be automatically
  extracted to temporary directory and is callable directly using `asName`
  chosen above because path of the temporary directory has been added into the
  PATH variable.

### 3.2. Build/run/clean

Formatting is managed exclusively by pre-commit hooks.

#### 3.2.1. Build

```bash
build/build-docker.sh
```

```bash
build/build-local.sh
```

#### 3.2.2. Tests

```bash
build/test.sh
```

#### 3.2.3. Coverage

```bash
build/coverage.sh
```

#### 3.2.4. run the binary

```bash
build/run.sh
```

#### 3.2.5. Clean

```bash
build/clean.sh
```

## 4. Commands

Compile bin file

```bash
go run ./cmd/bash-compiler examples/configReference/shellcheckLint.yaml \
  --root-dir /home/wsl/fchastanet/bash-dev-env/vendor/bash-tools-framework \
  -t examples/generated -k -d
```

for debugging purpose, manually Transform and validate yaml file using cue

```bash
cue export \
  -l input: examples/generated/shellcheckLint-merged.yaml \
  internal/model/binFile.cue --out yaml \
  -e output >examples/generated/shellcheckLint-cue-transformed.yaml
```

## 5. KCL

<https://www.kcl-lang.io/docs/user_docs/getting-started/install>

```bash
cd internal/model/kcl
kcl -D configFile=testsKcl/example.yaml
```

## 6. Alternatives

- Convert ecmascript to bash
  - <https://github.com/Ph0enixKM/Amber> alpha version - 2024-05-25
