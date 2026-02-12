# Copilot Instructions for bash-compiler

## Project Overview

**bash-compiler** is a Go-based tool that compiles bash scripts by detecting and
inlining framework functions that follow the pattern `Namespace::functionName`.
It processes YAML configuration files that describe bash applications, interprets
templates, and imports necessary bash functions recursively.

**Key Features:**

- Compiles single-file bash executables from multiple sources
- Template-based code generation using Go's `text/template`
- Recursive function dependency resolution
- File/directory embedding with base64 encoding
- YAML-based configuration with CUE language transformations
- Extensive linting and validation

## Repository Structure

```text
bash-compiler/
├── cmd/bash-compiler/        # Main entry point and CLI
│   ├── args.go              # Kong-based CLI argument parsing
│   ├── main.go              # Application entry point
│   └── defaultTemplates/    # Built-in Go templates (*.gtpl)
├── internal/
│   ├── compiler/            # Core compiler logic
│   │   ├── compiler.go      # Main compilation orchestration
│   │   ├── annotationEmbed.go    # @embed directive processing
│   │   └── annotationRequire.go  # Function dependency tracking
│   ├── model/               # Data structures and YAML models
│   ├── render/              # Template rendering engine
│   │   ├── templates.go     # Template loading and parsing
│   │   ├── funcmap.go       # Custom template functions
│   │   └── context.go       # Template execution context
│   ├── services/            # High-level service orchestration
│   └── utils/               # Utility packages
├── pkg/                     # Public API (if any)
├── examples/                # Example configurations
│   └── configReference/     # Reference YAML configs
├── build/                   # Build scripts
├── doc/                     # Diagrams and documentation
└── .github/workflows/       # CI/CD pipelines
```

## Development Workflow

### Building

```bash
# Local build (recommended for development)
./build/build-local.sh

# Docker build
./build/build-docker.sh

# Binary location after build
~/go/bin/bash-compiler
```

### Testing

```bash
# Run all tests
./build/test.sh

# Run specific package tests
go test ./internal/compiler/...

# Run tests with coverage
./build/coverage.sh

# Tests are located in *_test.go files alongside source
```

**Important:** Tests use:

- `github.com/stretchr/testify` for assertions
- `gotest.tools/v3` for advanced test utilities
- Table-driven tests are common pattern

### Linting and Validation

The project uses **MegaLinter** with extensive linters:

```bash
# Pre-commit hooks (must be installed first)
pre-commit install --hook-type pre-commit --hook-type pre-push
pre-commit run --all-files

# Specific linters
go fmt ./...                    # Go formatting
shellcheck script.sh            # Bash script linting
golangci-lint run              # Go linting (currently disabled in CI for Go 1.23)
```

**Key Linter Configs:**

- `.golangci.yml` - golangci-lint configuration (strict settings)
- `.shellcheckrc` - ShellCheck rules for bash validation
- `.mega-linter.yml` - MegaLinter orchestration
- `.pre-commit-config.yaml` - Pre-commit hook definitions

### Running the Tool

```bash
# Basic usage
~/go/bin/bash-compiler examples/configReference/shellcheckLint.yaml

# With options
~/go/bin/bash-compiler <yaml-file> \
  --intermediate-files-dir examples/generated \
  --debug

# Get help
~/go/bin/bash-compiler --help
```

## Go Language Conventions

### Code Style

- **Indentation:** Tabs for Go files (enforced by `.editorconfig`)
- **Line Length:** No hard limit, but aim for readability
- **Error Handling:** Explicit error returns; use custom error types
- **Naming:** Follow standard Go conventions
  - PascalCase for exported identifiers
  - camelCase for unexported identifiers
  - Avoid stuttering (e.g., `compiler.CompilerConfig` → `compiler.Config`)

### Common Patterns in This Project

1. **Custom Error Types:**

   ```go
   type rootDirError struct{ error }
   func (*rootDirError) Error() string { return "message" }
   ```

2. **Kong CLI Parsing:**

   ```go
   type cli struct {
       YamlFiles []string `arg:"" optional:"" type:"path" help:"..."`
       Debug     bool     `short:"d" help:"..."`
   }
   ```

3. **Validation Methods:**

   ```go
   func (c *Config) Validate() error {
       if c.Field == "" {
           return fmt.Errorf("field required")
       }
       return nil
   }
   ```

4. **Context Structs for Template Rendering:**

   ```go
   type Context struct {
       Template *template.Template
       Name     string
       RootData any
       Data     any
   }
   ```

### Important Go Modules

- **Kong:** CLI argument parsing (`github.com/alecthomas/kong`)
- **Sprig:** Template functions (`github.com/Masterminds/sprig/v3`)
- **go-yaml:** YAML parsing (`github.com/goccy/go-yaml`)
- **KCL:** Configuration language (`kcl-lang.io/kcl-go`)
- **slog:** Structured logging (standard library)

## Template System

### Overview

The project uses Go's `text/template` (not `html/template`) to generate bash
code. Templates have `.gtpl` extension.

### Template Locations

- `cmd/bash-compiler/defaultTemplates/` - Built-in templates
- Users can provide custom templates via YAML config

### Custom Template Functions

Located in `internal/render/funcmap.go`:

- **String functions:**
  - `stringLength` - Get string length
  - `format` - Printf-style formatting (e.g., `format "${%sVar}" .name`)
- **Template inclusion:**
  - `include` - Include another template with data override
  - `includeFile` - Include template by filename
  - `includeFileAsTemplate` - Include and interpret as template
  - `dynamicFile` - Resolve first matching filepath from list
- **Plus all Sprig functions:** <https://masterminds.github.io/sprig/>

### Template Context

Templates receive a `Context` struct with:

- `.Template` - Original template reference
- `.Name` - Template name
- `.RootData` - Initial data passed to rendering
- `.Data` - Current data (may be subset of RootData)

**Important:** Sprig functions don't receive execution context, so we simulate
context by passing it through the render function.

## YAML Configuration

### Structure

Configuration files in `examples/configReference/` show the structure:

```yaml
extends:
  - defaultCommand.yaml
  - frameworkConfig.yaml

vars:
  SRC_FILE_PATH: src/_binaries/commandDefinitions/...

compilerConfig:
  targetFile: "${FRAMEWORK_ROOT_DIR}/bin/shellcheckLint"
  templateFile: binFile.gtpl

binData:
  commands:
    default:
      functionName: shellcheckLintCommand
      version: "1.0"
      ...
```

### Configuration Transformation

The project uses **CUE** language to transform/validate YAML:

```bash
cue export \
  -l input: examples/generated/file-merged.yaml \
  internal/model/binFile.cue --out yaml \
  -e output > transformed.yaml
```

## Annotations and Directives

### @embed Directive

Allows embedding files/directories into compiled bash scripts:

**Syntax:**

```bash
# @embed "srcFile" AS "targetFile"
# @embed "srcDir" AS "targetDir"
```

**Behavior:**

- Files/dirs are base64-encoded and embedded in output
- Automatically generates extraction functions:
  - `Compiler::Embed::extractFile_<asName>`
  - `Compiler::Embed::extractDir_<asName>`
- Variables created: `embed_file_<asName>`, `embed_dir_<asName>`

### FUNCTIONS Directive

```bash
# FUNCTIONS
```

Marks where framework functions should be injected during compilation.

### Framework Function Pattern

Functions must follow: `Namespace::Namespace::functionName`

- Multiple namespaces separated by `::`
- Namespace: starts with `[A-Z]`, contains `[A-Za-z0-9_-]`
- Function name: camelCase, starts lowercase, contains `[a-zA-Z0-9_-]`
- File location: `srcDirs/Namespace/functionName.sh`

## Testing Approach

### Test File Organization

- Tests in `*_test.go` files next to source
- Test data in `testsData/` directories
- Table-driven tests are preferred

### Example Test Pattern

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {name: "valid input", input: "foo", want: "bar"},
        {name: "invalid input", input: "", wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Something(tt.input)
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

### Coverage

- Coverage reports in `logs/coverage.log`
- CI requires reasonable coverage (thresholds: 60% warning, 80% good)
- Use `./build/coverage.sh` to generate local reports

## CI/CD Workflows

### Main Workflow (.github/workflows/main.yml)

**Jobs:**

1. **build-docker-images** - Builds and pushes Docker image
   - Go 1.23.1, Ubuntu 22.04
   - Uses BuildKit caching
   - Pushes to DockerHub

2. **pre-commit** - Runs MegaLinter
   - Validates all code style and quality
   - Can auto-fix and create PR with fixes
   - Requires secrets: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`, `GPG_PRIVATE_KEY`

3. **test** - Unit tests and coverage
   - Runs `go test` with race detector
   - Generates coverage reports
   - Posts coverage summary to PR

**Workflow Triggers:**

- Push to `master` branch
- Pull requests to any branch
- Manual dispatch
- Excludes `docs/**` changes

### Status Checks

Uses `akatov/commit-status-updater` to create granular status checks:

- `build-docker`
- `pre-commit-megalinter`
- `build-ubuntu-latest`
- `unit-tests-ubuntu-latest`

## Commit Message Format

**Template (from `commit-msg-template.md`):**

```markdown
Title

Short 3 lines summary

tag: 1.0.0

### Breaking changes

### Bug fixes

### Compiler changes

### Binaries changes

### Updated Bash framework functions

### New Bash framework functions

### Documentation

### Validation/Tooling
```

**Requirements:**

- Markdown format
- Title summarizes changes
- Contains every relevant change
- Use sections as appropriate

## Git Conventions

- **Main branch:** `master` (not `main`)
- **Line endings:** LF (enforced by `.editorconfig` and pre-commit)
- **Gitignore:** Excludes `bin/`, `logs/`, `node_modules/`, `examples/generated/`

## Common Gotchas and Tips

### 1. Go Version

- Project uses **Go 1.24** (specified in `go.mod`)
- CI uses **Go 1.23.1** (specified in workflow)
- Docker uses **Go 1.24** (specified in Dockerfile)
- **Action:** Ensure compatibility across versions

### 2. Template Parsing

- Templates are parsed once at startup
- `ParseGlob` has a bug with multiple calls - avoid it
- Manual template list computation in `internal/render/render.go`

### 3. golangci-lint Disabled in CI

Currently disabled for Go 1.23 compatibility:

```yaml
DISABLE_LINTERS:
  - GO_GOLANGCI_LINT  # TEMPORARY waiting megalinter to support go 1.23
```

**Action:** Re-enable when MegaLinter supports Go 1.23

### 4. Pre-commit Hook Generation

`.github/preCommitGeneration.sh` generates GitHub-specific pre-commit config.
Always run this if modifying `.pre-commit-config.yaml`.

### 5. Shellcheck Configuration

`.shellcheckrc` has specific settings:

```bash
external-sources=true
enable=require-variable-braces
enable=avoid-nullary-conditions
enable=add-default-case
enable=quote-safe-variables
enable=require-double-brackets
source-path=SCRIPTDIR
```

These are strict - follow them in generated bash code.

### 6. GitHub Workflow Line Length

**Important:** Split lines longer than 120 characters in GitHub workflow files.

### 7. RootDir Option

- `--rootDir` flag is hidden in production (only for `go run` or debug)
- Auto-detected when running installed binary
- Detection logic in `cmd/bash-compiler/args.go:isUsingGoRun()`

### 8. Test Data

- Test data files in `testsData/` are excluded from some linters
- Example: end-of-file-fixer and trailing-whitespace hooks skip them

## Documentation Resources

### In-Repository

- **README.md** - Quick start and overview
- **CompileCommand.md** - Detailed compiler documentation
- **doc/** - PlantUML diagrams (class, activity, dependency)
- **examples/configReference/** - Example YAML configurations

### External References

- Go template docs: <https://pkg.go.dev/text/template>
- Sprig functions: <https://masterminds.github.io/sprig/>
- Kong CLI: <https://github.com/alecthomas/kong>
- KCL language: <https://www.kcl-lang.io/>

## Making Changes

### Typical Change Workflow

1. **Explore** - Use grep/glob to find relevant code
2. **Understand** - Read existing tests to understand behavior
3. **Build** - Run `./build/build-local.sh` to verify compilation
4. **Test** - Create/modify tests; run `./build/test.sh`
5. **Lint** - Run pre-commit hooks or rely on CI
6. **Commit** - Follow commit message format
7. **CI** - Monitor GitHub Actions for failures

### When Modifying Templates

1. Update template in `cmd/bash-compiler/defaultTemplates/`
2. Test with example configs in `examples/configReference/`
3. Generate output with `--intermediate-files-dir` to inspect results
4. Validate generated bash with shellcheck

### When Modifying Compiler

1. Update code in `internal/compiler/`
2. Add tests in `internal/compiler/*_test.go`
3. Test with realistic YAML configs
4. Check generated output for correctness

### When Adding Dependencies

1. Use `go get <package>`
2. Run `go mod tidy`
3. Verify `go.mod` and `go.sum` are updated
4. Check license compatibility (project is MIT licensed)

## Security Considerations

- **Secrets:** Never commit secrets; use GitHub Secrets for CI
- **Input Validation:** Always validate YAML inputs and file paths
- **Shell Injection:** Be careful with template generation - validate inputs
- **Docker Security:** Base image uses distroless for minimal attack surface

## Troubleshooting

### Build Fails

- Check Go version: `go version`
- Clean and rebuild: `./build/clean.sh && ./build/build-local.sh`
- Check module cache: `go clean -modcache`

### Tests Fail

- Run single test: `go test -v -run TestName ./internal/package`
- Check test data paths (absolute paths required)
- Review test output carefully - table-driven tests show which case failed

### Pre-commit Fails

- Run manually: `pre-commit run --all-files`
- Check specific hook: `pre-commit run <hook-id> --all-files`
- Some hooks auto-fix - stage changes and retry

### CI Failures

- Check GitHub Actions logs
- MegaLinter reports are uploaded as artifacts
- Coverage reports show which code lacks tests
- Docker build failures - check Dockerfile and base image

## Related Projects

This is part of a suite:

- **bash-tools-framework** - Framework functions source
- **bash-tools** - Command-line tools
- **bash-dev-env** - Development environment
- **my-documents** - Documentation site

## Quick Reference

### Essential Commands

```bash
# Build
./build/build-local.sh

# Test
./build/test.sh

# Run
~/go/bin/bash-compiler <yaml-file>

# Clean
./build/clean.sh

# Coverage
./build/coverage.sh

# Pre-commit
pre-commit run --all-files
```

### File Extensions

- `.go` - Go source files
- `.gtpl` - Go template files
- `.yaml` / `.yml` - Configuration files
- `.sh` - Bash scripts
- `*_test.go` - Go test files
- `*-binary.yaml` - Binary configuration files

### Key Directories to Watch

- `internal/compiler/` - Core compilation logic
- `internal/render/` - Template engine
- `cmd/bash-compiler/defaultTemplates/` - Built-in templates
- `examples/configReference/` - Example configurations

## Summary

**bash-compiler** is a sophisticated tool for building single-file bash
executables from modular source code. When working on this project:

- Understand the template system thoroughly
- Test changes with example configurations
- Follow strict linting rules (especially shellcheck)
- Use table-driven tests
- Commit messages should be detailed and structured
- CI is comprehensive - let it validate your work

**Remember:** This project values correctness, maintainability, and strict
validation over quick hacks. Take time to understand the architecture before
making changes.
