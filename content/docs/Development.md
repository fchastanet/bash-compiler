---
title: Development
description: Guidelines for developing and contributing to bash-compiler
weight: 40
categories: [documentation]
tags: [development, contribution, guidelines]
creationDate: 2025-04-09
lastUpdated: 2026-03-01
version: '2.0'
---

## 1. Requirements

### 1.1. Go Version

This project requires **Go 1.25.7** or later.

To check your current Go version:

```bash
go version
```

To install or upgrade Go, visit [golang.org/dl](https://golang.org/dl).

### 1.2. Pre-commit Hook

This repository uses pre-commit software to ensure every commit respects a set of rules specified by the
`.pre-commit-config.yaml` file. It requires [pre-commit](https://pre-commit.com/#install) to be installed.

Enable pre-commit hooks with:

```bash
pre-commit install --hook-type pre-commit --hook-type pre-push
```

Now linters and compilation tools will run automatically on commit and push.

### 1.3. Pre-commit dependencies

You need to install some dependencies for the pre-commit hooks to work properly. You can do it with:

```bash
.github/scripts/install-dev.sh
```

## 2. Build/Run/Clean

Formatting is managed exclusively by pre-commit hooks.

### 2.1. Build

Build with Docker:

```bash
.github/scripts/build-docker.sh
```

Build locally (requires Go 1.25.7):

```bash
.github/scripts/build-local.sh
```

### 2.2. Tests

Run all tests with race detector:

```bash
.github/scripts/test.sh
```

Run tests for specific package:

```bash
go test -v -race ./internal/compiler/...
```

### 2.3. Coverage

Generate coverage report:

```bash
.github/scripts/coverage.sh
```

Coverage reports are generated in `logs/coverage.log`.

### 2.4. Run the Binary

```bash
.github/scripts/run.sh
```

### 2.5. Clean

```bash
.github/scripts/clean.sh
```

## 3. Dependencies Management

### 3.1. Updating Go Dependencies

The project uses Go modules for dependency management. To update dependencies:

#### 3.1.1. Update All Dependencies

```bash
go get -u ./...
```

Then tidy and verify:

```bash
go mod tidy
go mod verify
```

#### 3.1.2. Update Specific Package

```bash
go get -u github.com/example/package
```

#### 3.1.3. Downgrade a Package

```bash
go get github.com/example/package@v1.0.0
```

#### 3.1.4. Upgrading Go Version

1. **Check current version:**

   ```bash
   go version
   ```

2. **Update go.mod:**

   ```bash
   go get -u golang.org/x/net golang.org/x/crypto golang.org/x/sys golang.org/x/text
   go get -u ./...
   ```

3. **Update go.mod with new version:**

   ```bash
   # Edit go.mod manually or use:
   grep -n "^go " go.mod
   # Then update the version number and toolchain directive
   ```

   Example: changing from 1.24 to 1.25.7:

   ```go
   module github.com/fchastanet/bash-compiler

   go 1.25.7

   toolchain go1.25.7
   ```

4. **Tidy and verify:**

   ```bash
   go mod tidy
   go mod verify
   ```

5. **Run tests:**

   ```bash
   go test ./... -race
   ```

6. **Build:**

   ```bash
   go build ./cmd/bash-compiler
   ```

## 4. Manual Compilation Commands

### 4.1. Compile Binary with Config

```bash
go run ./cmd/bash-compiler examples/configReference/shellcheckLint.yaml \
  --root-dir /path/to/bash-tools-framework \
  -t examples/generated -k -d
```

Flags:

- `-t, --intermediate-files-dir`: Output directory for generated files
- `-k`: Keep intermediate files
- `-d, --debug`: Enable debug logging

### 4.2. Transform and Validate YAML with CUE

```bash
cue export \
  -l input: examples/generated/shellcheckLint-merged.yaml \
  internal/model/binFile.cue --out yaml \
  -e output >examples/generated/shellcheckLint-cue-transformed.yaml
```

## 5. KCL Configuration Language

The project uses KCL for configuration validation. See
[KCL Documentation](https://www.kcl-lang.io/docs/user_docs/getting-started/install).

### 5.1. Test KCL Files

```bash
cd internal/model/kcl
kcl -D configFile=testsKcl/bad-example.yaml
```

## 6. Project Structure

Key directories:

- `.github/scripts/` - Build and test scripts
- `cmd/bash-compiler/` - Main entry point and CLI
- `internal/compiler/` - Core compilation logic
- `internal/model/` - Data structures and YAML models
- `internal/render/` - Template rendering engine
- `internal/utils/` - Utility packages
- `examples/configReference/` - Reference YAML configurations
- `content/docs/` - Documentation files

## 7. Common Workflows

### 7.1. Adding a New Dependency

```bash
go get github.com/example/package
go mod tidy
go test ./...
```

### 7.2. Running Specific Tests

```bash
# Table-driven test specific case
go test -v -run TestFunctionName/caseName ./internal/package

# All tests in package with coverage
go test -v -cover ./internal/package
```

### 7.3. Debugging

Enable debug logging in the compiler:

```bash
bash-compiler yaml-file -d
```

Check intermediate files:

```bash
bash-compiler yaml-file -t /tmp/debug -k
ls -la /tmp/debug/
```

### 7.4. Code Style

- **Indentation:** Tabs for Go files (enforced by `.editorconfig`)
- **Formatting:** Handled by pre-commit hooks (`gofmt`)
- **Linting:** Multiple linters via MegaLinter

To run linters manually:

```bash
pre-commit run --all-files
```

## 8. Testing

### 8.1. Test Organization

Tests use:

- `github.com/stretchr/testify` for assertions
- `gotest.tools/v3` for advanced utilities
- Table-driven test pattern (standard)

Example test pattern:

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {name: "case 1", input: "x", want: "y"},
        {name: "case 2", input: "foo", wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### 8.2. Test Coverage

Minimum coverage recommendations:

- **Warning:** 60% of statements
- **Good:** 80% or higher

View coverage:

```bash
go test ./... -cover
```

Generate detailed report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 9. CI/CD

GitHub Actions workflows:

- **main.yml** - Build Docker images, run MegaLinter, and tests
- Runs on push to `master`, pull requests, and manual dispatch
- Excludes changes to `docs/**`

Check workflow logs in: `.github/workflows/main.yml`

## 10. Troubleshooting

### 10.1. Build Fails

```bash
# Check Go version
go version

# Clean and rebuild
.github/scripts/clean.sh
.github/scripts/build-local.sh

# Check module cache
go clean -modcache
go mod tidy
```

### 10.2. Tests Fail

```bash
# Run with verbose output
go test -v -race ./...

# Run specific test
go test -v -run TestName ./package
```

### 10.3. Pre-commit Hooks Fail

```bash
# Run hooks manually
pre-commit run --all-files

# Run specific hook
pre-commit run hook-id --all-files
```

Some hooks auto-fix - stage changes and retry.

### 10.4. Module Issues

```bash
# Verify module integrity
go mod verify

# Download and check all modules
go mod download
```

## 11. Release Checklist

Before releasing:

1. ✅ All tests pass: `go test ./... -race`
2. ✅ Coverage acceptable: `.github/scripts/coverage.sh`
3. ✅ Pre-commit passes: `pre-commit run --all-files`
4. ✅ Documentation updated
5. ✅ Dependencies audit: `go mod verify`
6. ✅ Build successful: `go build ./cmd/bash-compiler`
7. ✅ Commit message follows guidelines (see `commit-msg-template.md`)

## 12. Additional Resources

- [Go Modules Documentation](https://golang.org/ref/mod)
- [Go Testing](https://golang.org/pkg/testing/)
- [Pre-commit Documentation](https://pre-commit.com/)
- [bash-compiler Compilation Command](https://bash-compiler.devlab.top/docs/compilecommand/)
- [Technical Architecture](https://bash-compiler.devlab.top/docs/technicalarchitecture/)
