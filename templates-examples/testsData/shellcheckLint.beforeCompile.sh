#!/usr/bin/env bash

###############################################################################
# GENERATED FROM https://github.com/fchastanetbash-compiler/tree/master/templates-examples/testsData/shellcheckLint.yaml
# DO NOT EDIT IT
# @generated
###############################################################################
# shellcheck disable=SC2288,SC2034

#!/usr/bin/env bash

# ensure that no user aliases could interfere with
# commands used in this script
unalias -a || true
shopt -u expand_aliases

# shellcheck disable=SC2034
((failures = 0)) || true

# Bash will remember & return the highest exit code in a chain of pipes.
# This way you can catch the error inside pipes, e.g. mysqldump | gzip
set -o pipefail
set -o errexit

# Command Substitution can inherit errexit option since bash v4.4
shopt -s inherit_errexit || true

# if set, and job control is not active, the shell runs the last command
# of a pipeline not executed in the background in the current shell
# environment.
shopt -s lastpipe

# a log is generated when a command fails
set -o errtrace

# use nullglob so that (file*.php) will return an empty array if no file
# matches the wildcard
shopt -s nullglob

# ensure regexp are interpreted without accentuated characters
export LC_ALL=POSIX

export TERM=xterm-256color

# avoid interactive install
export DEBIAN_FRONTEND=noninteractive
export DEBCONF_NONINTERACTIVE_SEEN=true

# store command arguments for later usage
# shellcheck disable=SC2034
declare -a BASH_FRAMEWORK_ARGV=("$@")
# shellcheck disable=SC2034
declare -a ORIGINAL_BASH_FRAMEWORK_ARGV=("$@")

# @see https://unix.stackexchange.com/a/386856
# shellcheck disable=SC2317
interruptManagement() {
  # restore SIGINT handler
  trap - INT
  # ensure that Ctrl-C is trapped by this script and not by sub process
  # report to the parent that we have indeed been interrupted
  kill -s INT "$$"
}
trap interruptManagement INT

################################################
# Temp dir management
################################################

KEEP_TEMP_FILES="${KEEP_TEMP_FILES:-0}"
export KEEP_TEMP_FILES

# PERSISTENT_TMPDIR is not deleted by traps
PERSISTENT_TMPDIR="${TMPDIR:-/tmp}/bash-framework"
export PERSISTENT_TMPDIR
if [[ ! -d "${PERSISTENT_TMPDIR}" ]]; then
  mkdir -p "${PERSISTENT_TMPDIR}"
fi

# shellcheck disable=SC2034
TMPDIR="$(mktemp -d -p "${PERSISTENT_TMPDIR:-/tmp}" -t bash-framework-$$-XXXXXX)"
export TMPDIR

# temp dir cleaning
# shellcheck disable=SC2317
cleanOnExit() {
  local rc=$?
  if [[ "${KEEP_TEMP_FILES:-0}" = "1" ]]; then
    Log::displayInfo "KEEP_TEMP_FILES=1 temp files kept here '${TMPDIR}'"
  elif [[ -n "${TMPDIR+xxx}" ]]; then
    Log::displayDebug "KEEP_TEMP_FILES=0 removing temp files '${TMPDIR}'"
    rm -Rf "${TMPDIR:-/tmp/fake}" >/dev/null 2>&1
  fi
  exit "${rc}"
}
trap cleanOnExit EXIT HUP QUIT ABRT TERM
#!/usr/bin/env bash

SCRIPT_NAME=${0##*/}
REAL_SCRIPT_FILE="$(readlink -e "$(realpath "${BASH_SOURCE[0]}")")"
if [[ -n "${EMBED_CURRENT_DIR}" ]]; then
  CURRENT_DIR="${EMBED_CURRENT_DIR}"
else
  CURRENT_DIR="${REAL_SCRIPT_FILE%/*}"
fi
# FUNCTIONS
#!/usr/bin/env bash
declare -a BASH_FRAMEWORK_ARGV_FILTERED=()

copyrightCallback() {
  if [[ -z "${copyrightBeginYear}" ]]; then
    copyrightBeginYear="$(date +%Y)"
  fi
  echo "Copyright (c) ${copyrightBeginYear}-now François Chastanet"
}

# shellcheck disable=SC2317 # if function is overridden
updateArgListInfoVerboseCallback() {
  BASH_FRAMEWORK_ARGV_FILTERED+=(--verbose)
}
# shellcheck disable=SC2317 # if function is overridden
updateArgListDebugVerboseCallback() {
  BASH_FRAMEWORK_ARGV_FILTERED+=(-vv)
}
# shellcheck disable=SC2317 # if function is overridden
updateArgListTraceVerboseCallback() {
  BASH_FRAMEWORK_ARGV_FILTERED+=(-vvv)
}
# shellcheck disable=SC2317 # if function is overridden
updateArgListEnvFileCallback() { :; }
# shellcheck disable=SC2317 # if function is overridden
updateArgListLogLevelCallback() { :; }
# shellcheck disable=SC2317 # if function is overridden
updateArgListDisplayLevelCallback() { :; }
# shellcheck disable=SC2317 # if function is overridden
updateArgListNoColorCallback() {
  BASH_FRAMEWORK_ARGV_FILTERED+=(--no-color)
}
# shellcheck disable=SC2317 # if function is overridden
updateArgListThemeCallback() { :; }
# shellcheck disable=SC2317 # if function is overridden
updateArgListQuietCallback() { :; }

# shellcheck disable=SC2317 # if function is overridden
optionHelpCallback() {
  <% ${commandFunctionName} %> help
  exit 0
}

# shellcheck disable=SC2317 # if function is overridden
optionVersionCallback() {
  echo "${SCRIPT_NAME} version <% ${versionNumber} %>"
  exit 0
}

# shellcheck disable=SC2317 # if function is overridden
optionEnvFileCallback() {
  local envFile="$2"
  Log::displayWarning "Command ${SCRIPT_NAME} - Option --env-file is deprecated and will be removed in the future"
  if [[ ! -f "${envFile}" || ! -r "${envFile}" ]]; then
    Log::displayError "Command ${SCRIPT_NAME} - Option --env-file - File '${envFile}' doesn't exist"
    exit 1
  fi
}

# shellcheck disable=SC2317 # if function is overridden
optionInfoVerboseCallback() {
  BASH_FRAMEWORK_ARGS_VERBOSE_OPTION='--verbose'
  BASH_FRAMEWORK_ARGS_VERBOSE=${__VERBOSE_LEVEL_INFO}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_INFO}" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionDebugVerboseCallback() {
  BASH_FRAMEWORK_ARGS_VERBOSE_OPTION='-vv'
  BASH_FRAMEWORK_ARGS_VERBOSE=${__VERBOSE_LEVEL_DEBUG}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_DEBUG}" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionTraceVerboseCallback() {
  BASH_FRAMEWORK_ARGS_VERBOSE_OPTION='-vvv'
  BASH_FRAMEWORK_ARGS_VERBOSE=${__VERBOSE_LEVEL_TRACE}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_DEBUG}" >> "${overrideEnvFile}"
}

getLevel() {
  local levelName="$1"
  case "${levelName^^}" in
    OFF)
      echo "${__LEVEL_OFF}"
      ;;
    ERR | ERROR)
      echo "${__LEVEL_ERROR}"
      ;;
    WARN | WARNING)
      echo "${__LEVEL_WARNING}"
      ;;
    INFO)
      echo "${__LEVEL_INFO}"
      ;;
    DEBUG | TRACE)
      echo "${__LEVEL_DEBUG}"
      ;;
    *)
      Log::displayError "Command ${SCRIPT_NAME} - Invalid level ${level}"
      return 1
  esac
}

getVerboseLevel() {
  local levelName="$1"
  case "${levelName^^}" in
    OFF)
      echo "${__VERBOSE_LEVEL_OFF}"
      ;;
    ERR | ERROR | WARN | WARNING | INFO)
      echo "${__VERBOSE_LEVEL_INFO}"
      ;;
    DEBUG)
      echo "${__VERBOSE_LEVEL_DEBUG}"
      ;;
    TRACE)
      echo "${__VERBOSE_LEVEL_TRACE}"
      ;;
    *)
      Log::displayError "Command ${SCRIPT_NAME} - Invalid level ${level}"
      return 1
  esac
}

# shellcheck disable=SC2317 # if function is overridden
optionDisplayLevelCallback() {
  local level="$2"
  local logLevel verboseLevel
  logLevel="$(getLevel "${level}")"
  verboseLevel="$(getVerboseLevel "${level}")"
  BASH_FRAMEWORK_ARGS_VERBOSE=${verboseLevel}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${logLevel}" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionLogLevelCallback() {
  local level="$2"
  local logLevel verboseLevel
  logLevel="$(getLevel "${level}")"
  verboseLevel="$(getVerboseLevel "${level}")"
  BASH_FRAMEWORK_ARGS_VERBOSE=${verboseLevel}
  echo "BASH_FRAMEWORK_LOG_LEVEL=${logLevel}" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionLogFileCallback() {
  local logFile="$2"
  echo "BASH_FRAMEWORK_LOG_FILE='${logFile}'" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionQuietCallback() {
  echo "BASH_FRAMEWORK_QUIET_MODE=1" >> "${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionNoColorCallback() {
  UI::theme "noColor"
}

# shellcheck disable=SC2317 # if function is overridden
optionThemeCallback() {
  UI::theme "$2"
}

displayConfig() {
  echo "Config"
  UI::drawLine "-"
  local var
  while read -r var; do
    printf '%-40s = %s\n' "${var}" "$(declare -p "${var}" | sed -E -e 's/^[^=]+=(.*)/\1/')"
  done < <(typeset -p | awk 'match($3, "^(BASH_FRAMEWORK_[^=]+)=", m) { print m[1] }' | sort)
  exit 0
}

optionBashFrameworkConfigCallback() {
  if [[ ! -f "$2" ]]; then
    Log::fatal "Command ${SCRIPT_NAME} - Bash framework config file '$2' does not exists"
  fi
}

defaultFrameworkConfig="$(
  cat <<'EOF'
.INCLUDE "${ORIGINAL_TEMPLATE_DIR}/_includes/.framework-config.default"
EOF
)"

overrideEnvFile="$(Framework::createTempFile "overrideEnvFile")"

commandOptionParseFinished() {
  # load default template framework config
  defaultEnvFile="${PERSISTENT_TMPDIR}/.framework-config"
  echo "${defaultFrameworkConfig}" > "${defaultEnvFile}"
  local -a files=("${defaultEnvFile}")
  if [[ -f "${envFile}" ]]; then
    files+=("${envFile}")
  fi
  # shellcheck disable=SC2154
  if [[ -f "${optionBashFrameworkConfig}" ]]; then
    files+=("${optionBashFrameworkConfig}")
  fi
  files+=("${overrideEnvFile}")
  Env::requireLoad "${files[@]}"
  Log::requireLoad
  # shellcheck disable=SC2154
  if [[ "${optionConfig}" = "1" ]]; then
    displayConfig
  fi
}

#!/usr/bin/env bash
declare MIN_SHELLCHECK_VERSION="0.9.0"
declare versionNumber="1.0"
declare commandFunctionName="shellcheckLintCommand"
declare help="Lint bash files using shellcheck."
declare optionFormatDefault="tty"

declare copyrightBeginYear="2022"
declare optionFormat="${optionFormatDefault}"
declare -a shellcheckArgs=()
declare -a shellcheckFiles=()

longDescriptionFunction() {
  Array::wrap2 ' ' 76 0 "${longDescription[@]}"
}

unknownOption() {
  shellcheckArgs+=("$1")
}
argShellcheckFilesCallback() {
  if [[ -f "$1" ]]; then
    shellcheckFiles=("${@::$#-1}")
  else
    shellcheckArgs+=("$1")
  fi
}
shellcheckLintParseCallback() {
  if [[ "${optionStaged}" = "1" ]] && ((${#argShellcheckFiles[@]} > 0)); then
    Log::displayWarning "${SCRIPT_NAME} - --staged option ignored as files have been provided"
    optionStaged="0"
  fi
  shellcheckArgs=(-f "${optionFormat}")
}

optionHelpCallback() {
  local shellcheckHelpFile
  shellcheckHelpFile="$(Framework::createTempFile "shellcheckHelp")"
  (
  if [[ -x "${FRAMEWORK_VENDOR_BIN_DIR}/shellcheck" ]]; then
    "${FRAMEWORK_VENDOR_BIN_DIR}/shellcheck" --help
  else
    Log::displayError "${FRAMEWORK_VENDOR_BIN_DIR}/shellcheck does not exist" 2>&1
  fi
  ) >"${shellcheckHelpFile}" 2>&1

  shellcheckLintCommandHelp |
    sed -E \
      -e "/@@@SHELLCHECK_HELP@@@/r ${shellcheckHelpFile}" \
      -e "/@@@SHELLCHECK_HELP@@@/d"
  exit 0
}


optionVersionCallback() {
  echo -e "${__HELP_TITLE_COLOR}${SCRIPT_NAME} version: ${__RESET_COLOR} ${versionNumber}"
  echo -e -n "${__HELP_TITLE_COLOR}shellcheck Version: ${__RESET_COLOR}"
  "${FRAMEWORK_VENDOR_BIN_DIR}/shellcheck" --version
  exit 0
}



# ------------------------------------------
# Command shellcheckLintCommand
# ------------------------------------------

# options variables initialization
declare optionFormat="tty"
declare optionStaged="0"
declare optionXargs="0"
declare optionHelp="0"
declare optionConfig="0"
declare optionBashFrameworkConfig=""
declare optionInfoVerbose="0"
declare optionDebugVerbose="0"
declare optionTraceVerbose="0"
declare -a optionEnvFiles=()
declare optionLogLevel=""
declare optionLogFile=""
declare optionDisplayLevel=""
declare optionNoColor="0"
declare optionTheme="default"
declare optionVersion="0"
declare optionQuiet="0"
# arguments variables initialization
declare -a argShellcheckFiles=()
# @description parse command options and arguments for shellcheckLintCommand
shellcheckLintCommandParse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  optionFormat="tty"
  local -i options_parse_optionParsedCountOptionFormat
  ((options_parse_optionParsedCountOptionFormat = 0)) || true
  optionStaged="0"
  local -i options_parse_optionParsedCountOptionStaged
  ((options_parse_optionParsedCountOptionStaged = 0)) || true
  optionXargs="0"
  local -i options_parse_optionParsedCountOptionXargs
  ((options_parse_optionParsedCountOptionXargs = 0)) || true
  optionHelp="0"
  local -i options_parse_optionParsedCountOptionHelp
  ((options_parse_optionParsedCountOptionHelp = 0)) || true
  optionConfig="0"
  local -i options_parse_optionParsedCountOptionConfig
  ((options_parse_optionParsedCountOptionConfig = 0)) || true
  optionBashFrameworkConfig=""
  local -i options_parse_optionParsedCountOptionBashFrameworkConfig
  ((options_parse_optionParsedCountOptionBashFrameworkConfig = 0)) || true
  optionInfoVerbose="0"
  local -i options_parse_optionParsedCountOptionInfoVerbose
  ((options_parse_optionParsedCountOptionInfoVerbose = 0)) || true
  optionDebugVerbose="0"
  local -i options_parse_optionParsedCountOptionDebugVerbose
  ((options_parse_optionParsedCountOptionDebugVerbose = 0)) || true
  optionTraceVerbose="0"
  local -i options_parse_optionParsedCountOptionTraceVerbose
  ((options_parse_optionParsedCountOptionTraceVerbose = 0)) || true
  
  optionLogLevel=""
  local -i options_parse_optionParsedCountOptionLogLevel
  ((options_parse_optionParsedCountOptionLogLevel = 0)) || true
  optionLogFile=""
  local -i options_parse_optionParsedCountOptionLogFile
  ((options_parse_optionParsedCountOptionLogFile = 0)) || true
  optionDisplayLevel=""
  local -i options_parse_optionParsedCountOptionDisplayLevel
  ((options_parse_optionParsedCountOptionDisplayLevel = 0)) || true
  optionNoColor="0"
  local -i options_parse_optionParsedCountOptionNoColor
  ((options_parse_optionParsedCountOptionNoColor = 0)) || true
  optionTheme="default"
  local -i options_parse_optionParsedCountOptionTheme
  ((options_parse_optionParsedCountOptionTheme = 0)) || true
  optionVersion="0"
  local -i options_parse_optionParsedCountOptionVersion
  ((options_parse_optionParsedCountOptionVersion = 0)) || true
  optionQuiet="0"
  local -i options_parse_optionParsedCountOptionQuiet
  ((options_parse_optionParsedCountOptionQuiet = 0)) || true
  
  argShellcheckFiles=()
  
  
  # shellcheck disable=SC2034
  local -i options_parse_parsedArgIndex=0
  while (($# > 0)); do
    local options_parse_arg="$1"
    local argOptDefaultBehavior=0
    case "${options_parse_arg}" in
      # Option 1/17
      # optionFormat alts --format|-f
      # type: String min 0 max 1
      # authorizedValues: checkstyle|diff|gcc|json|json1|quiet|tty
      --format | -f)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        if [[ ! "$1" =~ checkstyle|diff|gcc|json|json1|quiet|tty ]]; then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - value '$1' is not part of authorized values([checkstyle diff gcc json json1 quiet tty])"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionFormat >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionFormat))
        # shellcheck disable=SC2034
        optionFormat="$1"
        
        
        ;;
      # Option 2/17
      # optionStaged alts --staged
      # type: Boolean min 0 max 1
      --staged)
        
        # shellcheck disable=SC2034
        optionStaged="1"
        
        
        if ((options_parse_optionParsedCountOptionStaged >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionStaged))
        
        
        ;;
      # Option 3/17
      # optionXargs alts --xargs
      # type: Boolean min 0 max 1
      --xargs)
        
        # shellcheck disable=SC2034
        optionXargs="1"
        
        
        if ((options_parse_optionParsedCountOptionXargs >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionXargs))
        
        
        ;;
      # Option 4/17
      # optionHelp alts --help|-h
      # type: Boolean min 0 max 1
      --help | -h)
        
        # shellcheck disable=SC2034
        optionHelp="1"
        
        
        if ((options_parse_optionParsedCountOptionHelp >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionHelp))
        
        
        optionHelpCallback "${options_parse_arg}" "${optionHelp}"
        
        
        ;;
      # Option 5/17
      # optionConfig alts --config
      # type: Boolean min 0 max 1
      --config)
        
        # shellcheck disable=SC2034
        optionConfig="1"
        
        
        if ((options_parse_optionParsedCountOptionConfig >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionConfig))
        
        
        ;;
      # Option 6/17
      # optionBashFrameworkConfig alts --bash-framework-config
      # type: String min 0 max 1
      --bash-framework-config)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionBashFrameworkConfig >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionBashFrameworkConfig))
        # shellcheck disable=SC2034
        optionBashFrameworkConfig="$1"
        
        
        optionBashFrameworkConfigCallback "${options_parse_arg}" "${optionBashFrameworkConfig}"
        
        
        ;;
      # Option 7/17
      # optionInfoVerbose alts --verbose|-v
      # type: Boolean min 0 max 1
      --verbose | -v)
        
        # shellcheck disable=SC2034
        optionInfoVerbose="1"
        
        
        if ((options_parse_optionParsedCountOptionInfoVerbose >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionInfoVerbose))
        
        
        optionInfoVerboseCallback "${options_parse_arg}" "${optionInfoVerbose}"
        
        
        updateArgListInfoVerboseCallback "${options_parse_arg}" "${optionInfoVerbose}"
        
        
        ;;
      # Option 8/17
      # optionDebugVerbose alts -vv
      # type: Boolean min 0 max 1
      -vv)
        
        # shellcheck disable=SC2034
        optionDebugVerbose="1"
        
        
        if ((options_parse_optionParsedCountOptionDebugVerbose >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionDebugVerbose))
        
        
        optionDebugVerboseCallback "${options_parse_arg}" "${optionDebugVerbose}"
        
        
        updateArgListDebugVerboseCallback "${options_parse_arg}" "${optionDebugVerbose}"
        
        
        ;;
      # Option 9/17
      # optionTraceVerbose alts -vvv
      # type: Boolean min 0 max 1
      -vvv)
        
        # shellcheck disable=SC2034
        optionTraceVerbose="1"
        
        
        if ((options_parse_optionParsedCountOptionTraceVerbose >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionTraceVerbose))
        
        
        optionTraceVerboseCallback "${options_parse_arg}" "${optionTraceVerbose}"
        
        
        updateArgListTraceVerboseCallback "${options_parse_arg}" "${optionTraceVerbose}"
        
        
        ;;
      # Option 10/17
      # optionEnvFiles alts --env-file
      # type: StringArray min 0 max -1
      --env-file)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        
        
        ((++options_parse_optionParsedCountOptionEnvFiles))
        optionEnvFiles+=("$1")
        
        
        optionEnvFileCallback "${options_parse_arg}" "${optionEnvFiles[@]}"
        
        
        updateArgListEnvFileCallback "${options_parse_arg}" "${optionEnvFiles[@]}"
        
        
        ;;
      # Option 11/17
      # optionLogLevel alts --log-level
      # type: String min 0 max 1
      # authorizedValues: OFF|ERR|ERROR|WARN|WARNING|INFO|DEBUG|TRACE
      --log-level)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        if [[ ! "$1" =~ OFF|ERR|ERROR|WARN|WARNING|INFO|DEBUG|TRACE ]]; then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - value '$1' is not part of authorized values([OFF ERR ERROR WARN WARNING INFO DEBUG TRACE])"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionLogLevel >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionLogLevel))
        # shellcheck disable=SC2034
        optionLogLevel="$1"
        
        
        optionLogLevelCallback "${options_parse_arg}" "${optionLogLevel}"
        
        
        updateArgListLogLevelCallback "${options_parse_arg}" "${optionLogLevel}"
        
        
        ;;
      # Option 12/17
      # optionLogFile alts --log-file
      # type: String min 0 max 1
      --log-file)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionLogFile >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionLogFile))
        # shellcheck disable=SC2034
        optionLogFile="$1"
        
        
        optionLogFileCallback "${options_parse_arg}" "${optionLogFile}"
        
        
        updateArgListLogFileCallback "${options_parse_arg}" "${optionLogFile}"
        
        
        ;;
      # Option 13/17
      # optionDisplayLevel alts --display-level
      # type: String min 0 max 1
      # authorizedValues: OFF|ERR|ERROR|WARN|WARNING|INFO|DEBUG|TRACE
      --display-level)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        if [[ ! "$1" =~ OFF|ERR|ERROR|WARN|WARNING|INFO|DEBUG|TRACE ]]; then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - value '$1' is not part of authorized values([OFF ERR ERROR WARN WARNING INFO DEBUG TRACE])"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionDisplayLevel >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionDisplayLevel))
        # shellcheck disable=SC2034
        optionDisplayLevel="$1"
        
        
        optionDisplayLevelCallback "${options_parse_arg}" "${optionDisplayLevel}"
        
        
        updateArgListDisplayLevelCallback "${options_parse_arg}" "${optionDisplayLevel}"
        
        
        ;;
      # Option 14/17
      # optionNoColor alts --no-color
      # type: Boolean min 0 max 1
      --no-color)
        
        # shellcheck disable=SC2034
        optionNoColor="1"
        
        
        if ((options_parse_optionParsedCountOptionNoColor >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionNoColor))
        
        
        optionNoColorCallback "${options_parse_arg}" "${optionNoColor}"
        
        
        updateArgListNoColorCallback "${options_parse_arg}" "${optionNoColor}"
        
        
        ;;
      # Option 15/17
      # optionTheme alts --theme
      # type: String min 0 max 1
      # authorizedValues: default|default-force|noColor
      --theme)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
          return 1
        fi
        
        if [[ ! "$1" =~ default|default-force|noColor ]]; then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - value '$1' is not part of authorized values([default default-force noColor])"
          return 1
        fi
        
        
        
        if ((options_parse_optionParsedCountOptionTheme >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionTheme))
        # shellcheck disable=SC2034
        optionTheme="$1"
        
        
        optionThemeCallback "${options_parse_arg}" "${optionTheme}"
        
        
        updateArgListThemeCallback "${options_parse_arg}" "${optionTheme}"
        
        
        ;;
      # Option 16/17
      # optionVersion alts --version
      # type: Boolean min 0 max 1
      --version)
        
        # shellcheck disable=SC2034
        optionVersion="1"
        
        
        if ((options_parse_optionParsedCountOptionVersion >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionVersion))
        
        
        optionVersionCallback "${options_parse_arg}" "${optionVersion}"
        
        
        ;;
      # Option 17/17
      # optionQuiet alts --quiet|-q
      # type: Boolean min 0 max 1
      --quiet | -q)
        
        # shellcheck disable=SC2034
        optionQuiet="1"
        
        
        if ((options_parse_optionParsedCountOptionQuiet >= 1 )); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - Maximum number of option occurrences reached(1)"
          return 1
        fi
        
        ((++options_parse_optionParsedCountOptionQuiet))
        
        
        optionQuietCallback "${options_parse_arg}" "${optionQuiet}"
        
        
        updateArgListQuietCallback "${options_parse_arg}" "${optionQuiet}"
        
        
        ;;
      
      -*)
        
        
        unknownOption "" "${options_parse_arg}" || argOptDefaultBehavior=$?
        
        ;;
      *)
        if ((0)); then
          # Technical if - never reached
          :
        
        # Argument 1/1
        # argShellcheckFiles min 0 max -1
        
        
        elif (( options_parse_parsedArgIndex >= 0 )); then
        
          ((++options_parse_argParsedCountArgShellcheckFiles))
          # shellcheck disable=SC2034
          
          # shellcheck disable=SC2034
          argShellcheckFiles+=("${options_parse_arg}")
          
          argShellcheckFilesCallback "${argShellcheckFiles[@]}" -- "${@:2}"
          
        
        # else too much args
        else
          
          if [[ "${argOptDefaultBehavior}" = "0" ]]; then
            # too much args and no unknownArgumentCallbacks configured
            Log::displayError "Command ${SCRIPT_NAME} - Argument - too much arguments provided: $*"
            return 1
          fi
          
        fi
        ;;
    esac
    shift || true
  done || return $?
  commandOptionParseFinished
  shellcheckLintParseCallback
  
}

# @description display command options and arguments help for shellcheckLintCommand
shellcheckLintCommandHelp() {
  Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
    "Lint bash files using shellcheck."
  echo
  
  # ------------------------------------------
  # usage section
  # ------------------------------------------
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" "shellcheckLint [OPTIONS] "
  echo
  # ------------------------------------------
  # usage/options section
  # ------------------------------------------
  optionsAltList=(
    "[--format|-f <format>]"
    "[--staged]"
    "[--xargs]"
    "[--help|-h]"
    "[--config]"
    "[--bash-framework-config <bash-framework-config>]"
    "[--verbose|-v]"
    "[-vv]"
    "[-vvv]"
    "[--env-file <env-file>]"
    "[--log-level <log-level>]"
    "[--log-file <log-file>]"
    "[--display-level <display-level>]"
    "[--no-color]"
    "[--theme <theme>]"
    "[--version]"
    "[--quiet|-q]"
  )
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" \
    "shellcheckLint" "${optionsAltList[@]}"
  
  
  # ------------------------------------------
  # options section
  # ------------------------------------------
  
  
  echo
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"
  echo -e "  ${__HELP_OPTION_COLOR}--format${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-f format${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    define output format of this command"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--staged${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    lint only staged git files(files added to file list to be committed) and which are beginning with a bash shebang."
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--xargs${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    uses parallelization(using xargs command) only if tty format"
  echo
  
  
  
  echo
  echo -e "${__HELP_TITLE_COLOR}GLOBAL OPTIONS:${__RESET_COLOR}"
  echo -e "  ${__HELP_OPTION_COLOR}--help${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-h${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Displays this command help"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--config${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Displays configuration"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--bash-framework-config bash-framework-config${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Use alternate bash framework configuration."
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--verbose${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-v${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Info level verbose mode (alias of --display-level INFO)"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}-vv${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Debug level verbose mode (alias of --display-level DEBUG)"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}-vvv${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Trace level verbose mode (alias of --display-level TRACE)"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--env-file env-file${__HELP_NORMAL} {list} (optional)"
  Array::wrap2 ' ' 76 4 "    Load the specified env file (deprecated, please use --bash-framework-config option instead)"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--log-level log-level${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Set log level"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--log-file log-file${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Set log file"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--display-level display-level${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Set display level"
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--no-color${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Produce monochrome output. alias of --theme noColor."
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--theme theme${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Choose color theme - default-force means colors will be produced even if command is piped."
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--version${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Print version information and quit."
  echo
  
  
  echo -e "  ${__HELP_OPTION_COLOR}--quiet${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-q${__HELP_NORMAL} {single}"
  Array::wrap2 ' ' 76 4 "    Quiet mode, doesn't display any output."
  echo
  
  
  # ------------------------------------------
  # longDescription section
  # ------------------------------------------
  echo
  declare -a shellcheckLintCommandLongDescription=(
  
  
  
  "shellcheck wrapper that will:"
  
  
  
  "- install new shellcheck version(${MIN_SHELLCHECK_VERSION}) automatically"
  
  
  
  $'\r'
  
  
  
  "- by default, lint all git files of this project which are beginning with a bash shebang"
  
  
  
  "except if the option --staged is passed"
  
  
  
  ""
  
  
  
  ${__HELP_TITLE}Special configuration .shellcheckrc:${__HELP_NORMAL}
  
  
  
  "use the following line in your .shellcheckrc file to exclude"
  
  
  
  "some files from being checked (use grep -E syntax)"
  
  
  
  "exclude=^bin/bash-tpl$"
  
  
  
  ""
  
  
  
  ${__HELP_TITLE_COLOR}SHELLCHECK HELP${__RESET_COLOR}
  
  
  
  ""
  
  
  
  "@@@SHELLCHECK_HELP@@@"
  
  
  
  ""
  
  
  )
  Array::wrap2 ' ' 76 0 "${shellcheckLintCommandLongDescription[@]}"
  echo
  # ------------------------------------------
  # version section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}VERSION: ${__RESET_COLOR}"
  echo '1.0'
  # ------------------------------------------
  # author section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}AUTHOR: ${__RESET_COLOR}"
  echo '[François Chastanet](https://github.com/fchastanet)'
  # ------------------------------------------
  # sourceFile section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}SOURCE FILE: ${__RESET_COLOR}"
  echo 'https://github.com/fchastanetbash-compiler/tree/master/templates-examples/testsData/shellcheckLint.yaml'
  # ------------------------------------------
  # license section
  # ------------------------------------------
  echo
  echo -n -e "${__HELP_TITLE_COLOR}LICENSE: ${__RESET_COLOR}"
  echo 'MIT License'
  # ------------------------------------------
  # copyright section
  # ------------------------------------------
  Array::wrap2 ' ' 76 0 """copyrightCallback"""
}

MAIN_FUNCTION_NAME="main"
main() {
  #!/usr/bin/env bash
  
  SCRIPT_NAME=${0##*/}
  REAL_SCRIPT_FILE="$(readlink -e "$(realpath "${BASH_SOURCE[0]}")")"
  if [[ -n "${EMBED_CURRENT_DIR}" ]]; then
    CURRENT_DIR="${EMBED_CURRENT_DIR}"
  else
    CURRENT_DIR="${REAL_SCRIPT_FILE%/*}"
  fi
  FRAMEWORK_ROOT_DIR="$(cd "${CURRENT_DIR}/.." && pwd -P)"
  FRAMEWORK_SRC_DIR="${FRAMEWORK_ROOT_DIR}/src"
  FRAMEWORK_BIN_DIR="${FRAMEWORK_ROOT_DIR}/bin"
  FRAMEWORK_VENDOR_DIR="${FRAMEWORK_ROOT_DIR}/vendor"
  FRAMEWORK_VENDOR_BIN_DIR="${FRAMEWORK_ROOT_DIR}/vendor/bin"
  Env::requireLoad
  UI::requireTheme
  Log::requireLoad
  shellcheckLintCommandParse "$@"
  
}

# if file is sourced avoid calling main function
# shellcheck disable=SC2178
BASH_SOURCE=".$0" # cannot be changed in bash
# shellcheck disable=SC2128
test ".$0" != ".${BASH_SOURCE}" || main "$@"