# Bash Compiler

[![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/SchwarzIT/go-template)

<!-- markdownlint-capture -->

<!-- markdownlint-disable MD013 -->

[![GitHub release (latest SemVer)](https://img.shields.io/github/release/fchastanet/bash-compiler?logo=github&sort=semver)](https://github.com/fchastanet/bash-compiler/releases)
[![GitHubLicense](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/fchastanet/bash-compiler/blob/master/LICENSE)
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
    - [2.1.1. Template system](#211-template-system)
- [3. Development](#3-development)
  - [3.1. Pre-commit hook](#31-pre-commit-hook)
  - [3.2. Build/run/clean](#32-buildrunclean)
    - [3.2.1. Build](#321-build)
    - [3.2.2. Tests](#322-tests)
    - [3.2.3. Coverage](#323-coverage)
    - [3.2.4. run the binary](#324-run-the-binary)
    - [3.2.5. Clean](#325-clean)

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

#### 2.1.1. Template system

[template system](https://pkg.go.dev/text/template@go1.22.3)
[doc 1](https://lets-go.alexedwards.net/sample/02.08-html-templating-and-inheritance.html)

There is the choice between Go template/text or template/html libraries I chosen
template/text to avoid some escaping that are not needed in bash.

Go template/text or template/html don't provide any execution context to the
filters (FuncMap).

I'm not using Template.ParseGlob because I have to call it twice to include
files of root directory and sub directories with 2 glob patterns. But a bug in
text/template makes the template be initialized again after each calls to
ParseGlob function.

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
- Name if the name of the first template that has been rendered
- RootData are the data that have been sent at the start of the rendering
- Data are the data sent to the sub template (possibly a part of RootData or the
  whole RootData)

Then each filter has to be called with the right context. The special filter
`include` allows to include a sub template overriding context Data.

Template filter functions: my current template filter functions are inspired by
[Sprig](https://github.com/Masterminds/sprig)

- I'm not using it because it is not maintained anymore
- a possible alternate fork is [sprout](https://github.com/go-sprout/sprout) but
  it misses a lot of functions.

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
