---
title: Bash Compiler
linkTitle: Bash Compiler
description: Go tool that inlines every bash framework function your script depends on into a single, self-contained, distributable file
weight: 10
type: docs
categories: [documentation]
tags: [compile, directives, templates, bash, golang]
date: '2026-02-14'
lastmod: '2026-09-17'
version: '2.0'
---

{{% pageinfo %}}

**bash-compiler** takes a bash script that calls framework functions spread over dozens of files and produces **one
self-contained script** with every dependency inlined — ready to copy, paste and run anywhere, with no installer and no
`source` statements.

{{% /pageinfo %}}

{{< img src="assets/bash-compiler-homepage.jpg" alt="Bash Compiler illustration" >}}

## 1. Why this tool exists

Distributing scripts that rely on multiple sourced files is inconvenient, as users must unpack and run them from
specific locations. To improve this, the goal is to create a single script file that embeds all necessary files. This
approach ensures:

- Easy distribution as a single file
- Reliable copy-paste across systems and editors
- Support for embedding binary (non-printable) files
- Reusability of bash functions

To achieve this, the script must store other files' contents in a way that avoids non-printable characters, which can
cause issues during copy-paste or transmission over messaging programs.

## 2. What the compiler does

This tool detects all the framework functions used inside a given sh file. The framework functions match the pattern
`Namespace::functionName` (several namespaces can be separated by the characters `::`). These framework functions are
injected inside a compiled file. The process is recursive, so every framework function used by imported framework
functions is imported as well — each one only once.

On top of function inlining, the compiler can:

- **embed files, directories and binaries** as base64 payloads that are extracted at runtime — see \[@embed
  directive\]({{< relref "Directives.md#3-embed-directive" >}})
- **order runtime requirements** (`@require`) so that preconditions are checked before the code that needs them runs —
  see \[Requirements\]({{< relref "Requirements.md" >}})
- **render the script through Go templates**, so the same source can produce different binaries — see \[Technical
  architecture\]({{< relref "TechnicalArchitecture.md" >}})

You can see several examples of compiled files by checking the
[bash-tools-framework src/\_binaries folder](https://github.com/fchastanet/bash-tools-framework/tree/master/src/_binaries).

## 3. What is a bash framework function ?

The so called `bash framework functions` are the functions defined in this framework that respect the following naming
convention:

- Namespace::Namespace::functionName
  - we can have any number of namespaces
  - each namespace is followed by ::
  - namespace must begin by an uppercase letter `[A-Z]` followed by any of these characters `[A-Za-z0-9_-]`.
  - the function name is traditionally written using camelCase with first letter in small case
  - function name authorized characters are `[a-zA-Z0-9_-]+`
- the function source code using namespace convention will be searched under the source directories declared in the
  compiler configuration
  - each namespace corresponds to a folder
  - the filename of the function is the function name with .sh extension
  - eg: `Filters::camel2snakeCase` source code can be found in `src/Filters/camel2snakeCase.sh`

## 4. Quick start

Build the compiler from source (Go 1.26.5 or later):

```bash
git clone https://github.com/fchastanet/bash-compiler.git
cd bash-compiler
go build -o bin/bash-compiler ./cmd/bash-compiler
```

Compile a binary from its yaml description:

```bash
cd /path/to/your/project # must contain a .bash-compiler file
bash-compiler examples/configReference/shellcheckLint.yaml
```

The compiler reads the yaml file describing the bash application, interprets the templates, imports the necessary bash
functions and writes the resulting single-file script. See \[Compile command\]({{< relref "CompileCommand.md" >}}) for
every option.

## 5. Documentation map

- \[Compile command\]({{< relref "CompileCommand.md" >}}) — CLI arguments, flags, configuration file and template
  variables
- \[Directives\]({{< relref "Directives.md" >}}) — `BIN_FILE`, `FUNCTIONS` and `@embed`, the annotations you write in
  your source file
- \[Requirements\]({{< relref "Requirements.md" >}}) — the `@require` annotation and how the dependency tree is ordered
- \[How the compiler works\]({{< relref "HowItWorks.md" >}}) — compilation passes, activity, class and dependency
  diagrams
- \[Technical architecture\]({{< relref "TechnicalArchitecture.md" >}}) — Go libraries, template engine and rendering
  context
- \[Development\]({{< relref "Development.md" >}}) — build, test, coverage, dependency management and contribution
  workflow

## 6. Alternatives

- Convert ecmascript to bash
  - <https://github.com/Ph0enixKM/Amber> alpha version - 2024-05-25

## 7. Related projects

This tool is part of a suite of bash tooling projects:

- [Bash Tools Framework](https://bash-tools-framework.devlab.top/) — the framework providing the functions being
  compiled
- [Bash Tools](https://bash-tools.devlab.top/) — a collection of ready-to-use binaries built with this compiler
- [Bash Dev Env](https://bash-dev-env.devlab.top/) — automated development environment installation
