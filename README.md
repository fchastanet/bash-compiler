# Bash Compiler

> **_NOTE:_** Documentation is best viewed on [github-pages](https://bash-compiler.devlab.top/)

<!-- markdownlint-capture -->

<!-- markdownlint-disable MD013 -->

[![GoTemplate](https://img.shields.io/badge/go/template-black?logo=go)](https://github.com/SchwarzIT/go-template)
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

This tool allows to detect all the framework functions used inside a given sh file. The framework functions matches the
pattern Namespace::functionName (we can have several namespaces separated by the characters ::). These framework
functions will be injected inside a compiled file. The process is recursive so that every framework functions used by
imported framework functions will be imported as well (of course only once).

> **_TIP:_** Checkout related projects of this suite
>
> - [My documents](https://devlab.top/)
> - [Bash Tools Framework](https://bash-tools-framework.devlab.top/)
> - [Bash Tools](https://bash-tools.devlab.top/)
> - [Bash Dev Env](https://bash-dev-env.devlab.top/)
> - **[Bash Compiler](https://bash-compiler.devlab.top/)**

## 1. Technical architecture

[Development](content/docs/Development.md) and [Technical architecture](content/docs/TechnicalArchitecture.md)
documentation are available for more details about the project.
