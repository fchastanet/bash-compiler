---
title: Requirements
description: The @require annotation - declaring runtime preconditions and letting the compiler order them through the dependency tree
weight: 40
type: docs
categories: [documentation]
tags: [require, requirements, dependencies, compile]
date: '2026-02-14T21:20:53+02:00'
lastmod: '2026-09-17'
version: '1.0'
---

{{% pageinfo %}}

The `@require` annotation declares a precondition that must hold before a function runs: a command being available, a
specific OS, or some loading step having happened. The compiler collects every requirement, orders them, and injects the
checks into the generated script.

{{% /pageinfo %}}

## 1. How the compiler resolves requirements

The compiler during successive passes:

- use existing compiler passes (injectImportedFunctions)
  - will parse `# @require` directives of each newly injected functions
    - error if require name does not begin with require
    - error if require name does not comply naming convention
    - error if `require*` file not found
  - will ignore the disabled requirements
  - a tree of require dependencies will be computed
  - we inject gradually the framework functions linked to the requires functions
- At the end of compiler processing
  - inject the requirements calls in the order specified by dependency tree (see below).

## 2. Requires dependencies tree

The following rules apply:

- Some requirements can depends on each others, the compiler will compute which dependency should be loaded before the
  other. _Eg:_ Log::requireLoad requirement depends on Framework::requireRootDir, so Framework::requireRootDir is loaded
  before. But Log requirement depends also on Env::requireLoad requirement.
- Requirement can be set at namespace level by adding the directive in `_.sh` file or at function level.
- A requirement can be loaded only once.
- A requirement that is used by several functions will be more prioritized and will be loaded before a less prioritized
  requirement.
- `# FUNCTIONS` placeholder should be defined before `# REQUIREMENTS` placeholder
- `# REQUIREMENTS` placeholder should be defined before `# ENTRYPOINT` placeholder

## 3. Require example

The annotation @require added to a function like in this example:

```bash
# @require Env::requireLoad
# @require Log::requireLoad
Log::logMessage() {
  # rest of the function content
  :
}
```

will do the following actions:

- compiler checks that the required functions exist, if not an error is triggered.
- compiler adds code to the required function that will set an environment variable to 1 when the function is called
  (eg: REQUIRE_FUNCTION_ENV_REQUIRE_LOAD_LOADED=1).
- compiler adds code to the function that has these requirements in to check if these environment variables are set and
  exit 1 if not.
- compiler checks if the function is called at least once but it is the developer's responsibility to call the require
  function at the right place.

Code is generated using go templates. The go templates are configured in the yaml file at compiler config level.

```yaml
compilerConfig:
  annotationsConfig:
    requireTemplate: require
    checkRequirementsTemplate: checkRequirements
# rest of the config file content
```

`examples/templates/annotations/require.gtpl` => generates this code:

```bash
Env::RequireLoad() {
  REQUIRE_FUNCTION_ENV_REQUIRE_LOAD_LOADED=1
  # rest of the function content
}
```

`examples/templates/annotations/checkRequirements.gtpl` => generates this code:

```bash
# @require Env::requireLoad
# @require Log::requireLoad
Log::logMessage() {
  if [[ "${REQUIRE_FUNCTION_ENV_REQUIRE_LOAD_LOADED:-0}" != 1 ]]; then
    echo >&2 "Requirement Env::requireLoad has not been loaded"
    exit 1
  fi

  if [[ "${REQUIRE_FUNCTION_LOG_REQUIRE_LOAD_LOADED:-0}" != 1 ]]; then
    echo >&2 "Requirement Log::requireLoad has not been loaded"
    exit 1
  fi
  # rest of the function content
}
```

The aims of a require are the following:

- be to be able to test for a requirement just before executing a function that is marked with @require
- when compiling be able to know if a function with a specific requirement has been used (eg: ubuntu>20)
- There are several kind of requirements:
  - checking that a command is available
    - this requirement needs to be called at the proper level if the binary actually installs this command.
    - @require Aws::requireAwsCommand
    - @require Docker::requireDockerCommand
    - @require Git::requireGitCommand
    - @require Linux::requireCurlCommand
    - @require Linux::requireJqCommand
    - @require Linux::requireRealpathCommand
    - @require Linux::requirePathchkCommand
    - @require Linux::requireSudoCommand
    - @require Linux::requireTarCommand
    - @require Ssh::requireSshKeygenCommand
    - @require Ssh::requireSshKeyscanCommand
  - checking a feature is available
    - @require Git::requireShallowClone actually based on git version
  - checking a specific environment/state is available on execution
    - @require Linux::requireUbuntu
    - @require Linux::Wsl::requireWsl
    - @require Linux::requireExecutedAsUser
    - ubuntu>20
  - ensuring some specific loading are made
    - @require Env::requireLoad
    - @require Log::requireLoad
    - @require UI::requireTheme

## 4. Requires dependencies use cases

_Script file example:_

```bash
# FUNCTIONS placeholder
# REQUIRES placeholder
Linux::Apt::update || Log::displayError "impossible to update"
```

- first compiler injectImportedFunctions pass
  - `Linux::Apt::update` requires
    - `Linux::requireSudoCommand`
    - `Linux::requireUbuntu`
  - `Log::display*` requires `Colors::requireTheme`
- second compiler injectImportedFunctions pass
  - `Log::log*` requires `Log::requireLoad`
- third compiler injectImportedFunctions pass
  - `Log::requireLoad` requires `Env::requireLoad`
- fourth compiler injectImportedFunctions pass
  - `Env::requireLoad` requires
    - `Framework::requireRootDir`
    - `Framework::tmpFileManagement` (see `src/_includes/_commonHeader.sh`)
- fifth compiler injectImportedFunctions pass
  - `Framework::tmpFileManagement` requires
    - `Framework::requireRootDir` which is already in the required list

If we order the requirements following reversed pass order, we end up with:

- `Framework::tmpFileManagement`
- `Framework::requireRootDir`
  - here we have an issue as it should come before `Framework::tmpFileManagement`
  - a solution could be to add the element to require list even if it is already in the list. This way it could even
    give a weight at certain requires.
- `Env::requireLoad`
- `Log::requireLoad`
- `Colors::requireTheme`
- `Linux::requireUbuntu`
- `Linux::requireSudoCommand`

To take into consideration:

- at each pass, we will parse the full list of functions and requires
  - it means the array of requires has to be reset at each pass.

Let's take again our above example, pass by pass (we avoided to include some functions intentionally like
`Retry:default` needed by `Linux::Apt::update` to make example easier to understand).

_Pass #1:_ import functions Linux::Apt::update and Log::displayError

```bash
# @require Linux::requireSudoCommand
# @require Linux::requireUbuntu
Linux::Apt::update() { :; }
# @require Log::requireLoad
Log::displayError() {
  #...
  Log:logMessage #...
}
# FUNCTIONS placeholder
# we don't have any yet as we are still parsing the 3 lines
# code above.
# REQUIRES placeholder
Linux::Apt::update || Log::displayError "impossible to update"
```

Functions imported list so far:

- Linux::Apt::update
- Log::displayError

_Pass #2:_ import functions Log:logMessage and import required functions in reverse order Linux::requireSudoCommand,
Linux::requireUbuntu, Log::requireLoad _Note:_ remember that require functions are only filtered using `# @require`

```bash
# @require Linux::requireSudoCommand
# @require Linux::requireUbuntu
Linux::Apt::update() { :; }
# @require Log::requireLoad
Log::displayError() {
  #...
  Log:logMessage #...
}
Log:logMessage() { :; }
Linux::requireSudoCommand() { :; }
Linux::requireUbuntu() { :; }
# @require Env::requireLoad
Log::requireLoad() { :; }
# FUNCTIONS placeholder

Log::requireLoad
Linux::requireUbuntu
Linux::requireSudoCommand
# REQUIRES placeholder
Linux::Apt::update || Log::displayError "impossible to update"
```

Functions imported list so far:

- Linux::Apt::update
- Log::displayError
- Log:logMessage
- Log::requireLoad
- Linux::requireSudoCommand
- Linux::requireUbuntu

_Pass #3:_ import functions, import required functions will import Env::requireLoad so order of requires will be:

```bash
Env:requireLoad
Log::requireLoad
Linux::requireUbuntu
Linux::requireSudoCommand
```

## 5. Next steps

- {{% mdlink text="Directives" path="/directives" %}} — the other annotations available in your source files
- {{% mdlink text="How the compiler works" path="/howitworks" %}} — the passes that drive this resolution
