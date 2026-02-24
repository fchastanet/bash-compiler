---
title: Development
description: Guidelines for developing and contributing to bash-compiler
weight: 40
categories: [documentation]
tags: [development, contribution, guidelines]
creationDate: 2025-04-09
lastUpdated: 2026-02-24
version: '1.0'
---

## 3.1. Pre-commit hook

This repository uses pre-commit software to ensure every commits respects a set of rules specified by the
`.pre-commit-config.yaml` file. It supposes pre-commit software is [installed](https://pre-commit.com/#install) in your
environment.

You also have to execute the following command to enable it:

```bash
pre-commit install --hook-type pre-commit --hook-type pre-push
```

Now each time you commit or push, some linters/compilation tools are launched automatically

## 3.2. Build/run/clean

Formatting is managed exclusively by pre-commit hooks.

### 3.2.1. Build

```bash
build/build-docker.sh
```

```bash
build/build-local.sh
```

### 3.2.2. Tests

```bash
build/test.sh
```

### 3.2.3. Coverage

```bash
build/coverage.sh
```

### 3.2.4. Run the binary

```bash
build/run.sh
```

### 3.2.5. Clean

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
