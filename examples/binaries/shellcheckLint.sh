#!/usr/bin/env bash

###############################################################################
# GENERATED FROM ${REPOSITORY_URL}/tree/master/${SRC_FILE_PATH}
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
#!/usr/bin/env bash
# @description Log namespace provides 2 kind of functions
# - Log::display* allows to display given message with
#   given display level
# - Log::log* allows to log given message with
#   given log level
# Log::display* functions automatically log the message too
# @see Env::requireLoad to load the display and log level from .env file

# @description log level off
export __LEVEL_OFF=0
# @description log level error
export __LEVEL_ERROR=1
# @description log level warning
export __LEVEL_WARNING=2
# @description log level info
export __LEVEL_INFO=3
# @description log level success
export __LEVEL_SUCCESS=3
# @description log level debug
export __LEVEL_DEBUG=4

# @description verbose level off
export __VERBOSE_LEVEL_OFF=0
# @description verbose level info
export __VERBOSE_LEVEL_INFO=1
# @description verbose level info
export __VERBOSE_LEVEL_DEBUG=2
# @description verbose level info
export __VERBOSE_LEVEL_TRACE=3
#!/usr/bin/env bash

# @description concatenate each element of an array with a separator
# but wrapping text when line length is more than provided argument
# The algorithm will try not to cut the array element if it can.
# - if an arg can be placed on current line it will be,
#   otherwise current line is printed and arg is added to the new
#   current line
# - Empty arg is interpreted as a new line.
# - Add \r to arg in order to force break line and avoid following
#   arg to be concatenated with current arg.
#
# @arg $1 glue:String
# @arg $2 maxLineLength:int
# @arg $3 indentNextLine:int
# @arg $@ array:String[]
Array::wrap2() {
  local glue="${1-}"
  local -i glueLength="${#glue}"
  shift || true
  local -i maxLineLength=$1
  shift || true
  local -i indentNextLine=$1
  shift || true
  local indentStr=""
  if ((indentNextLine > 0)); then
    indentStr="$(head -c "${indentNextLine}" </dev/zero | tr '\0' " ")"
  fi
  if (($# == 0)); then
    return 0
  fi

  printCurrentLine() {
    if ((isNewline == 0)) || ((previousLineEmpty == 1)); then
      echo
    fi
    ((isNewline = 1))
    echo -en "${indentStr}"
    ((currentLineLength = indentNextLine)) || true
  }
  appendToCurrentLine() {
    local text="$1"
    local -i length=$2
    ((currentLineLength += length)) || true
    ((isNewline = 0)) || true
    if [[ "${text: -1}" = $'\r' ]]; then
      text="${text:0:-1}"
      echo -en "${text%%+([[:blank:]])}"
      printCurrentLine
    else
      echo -en "${text%%+([[:blank:]])}"
    fi
  }

  (
    local currentLine
    local -i currentLineLength=0 isNewline=1 argLength=0
    local -a additionalLines
    local -i previousLineEmpty=0
    local arg=""

    while (($# > 0)); do
      arg="$1"
      shift || true

      # replace tab by 2 spaces
      arg="${arg//$'\t'/  }"
      # remove trailing spaces
      arg="${arg%[[:blank:]]}"
      if [[ "${arg}" = $'\n' || -z "${arg}" ]]; then
        printCurrentLine
        ((previousLineEmpty = 1))
        continue
      else
        if ((previousLineEmpty == 1)); then
          printCurrentLine
        fi
        ((previousLineEmpty = 0)) || true
      fi
      # convert eol to args
      mapfile -t additionalLines <<<"${arg}"
      if ((${#additionalLines[@]} > 1)); then
        set -- "${additionalLines[@]}" "$@"
        continue
      fi

      ((argLength = ${#arg})) || true

      # empty arg
      if ((argLength == 0)); then
        if ((isNewline == 0)); then
          # isNewline = 0 means currentLine is not empty
          printCurrentLine
        fi
        continue
      fi

      if ((isNewline == 0)); then
        glueLength="${#glue}"
      else
        glueLength="0"
      fi
      if ((currentLineLength + argLength + glueLength > maxLineLength)); then
        if ((argLength + glueLength > maxLineLength)); then
          # arg is too long to even fit on one line
          # we have to split the arg on current and next line
          local -i remainingLineLength
          ((remainingLineLength = maxLineLength - currentLineLength - glueLength))
          appendToCurrentLine "${glue:0:${glueLength}}${arg:0:${remainingLineLength}}" "$((glueLength + remainingLineLength))"
          printCurrentLine
          arg="${arg:${remainingLineLength}}"
          # remove leading spaces
          arg="${arg##[[:blank:]]}"

          set -- "${arg}" "$@"
        else
          # the arg can fit on next line
          printCurrentLine
          appendToCurrentLine "${arg}" "${argLength}"
        fi
      else
        appendToCurrentLine "${glue:0:${glueLength}}${arg}" "$((glueLength + argLength))"
      fi
    done
    if [[ "${currentLine}" != "" ]] && [[ ! "${currentLine}" =~ ^[\ \t]+$ ]]; then
      printCurrentLine
    fi
  ) | sed -E -e 's/[[:blank:]]+$//'
}
#!/usr/bin/env bash

# @description check if tty (interactive mode) is active
# @noargs
# @exitcode 1 if tty not active
# @env NON_INTERACTIVE if 1 consider as not interactive even if environment is interactive
# @env INTERACTIVE if 1 consider as interactive even if environment is not interactive
Assert::tty() {
  if [[ "${NON_INTERACTIVE:-0}" = "1" ]]; then
    return 1
  fi
  if [[ "${INTERACTIVE:-0}" = "1" ]]; then
    return 0
  fi
  tty -s
}
#!/usr/bin/env bash

# @description ensure env files are loaded
# @arg $@ list of default files to load at the end
# @exitcode 1 if one of env files fails to load
# @stderr diagnostics information is displayed
# shellcheck disable=SC2120
Env::requireLoad() {
  local -a defaultFiles=("$@")
  # get list of possible config files
  local -a configFiles=()
  if [[ -n "${BASH_FRAMEWORK_ENV_FILES[0]+1}" ]]; then
    # BASH_FRAMEWORK_ENV_FILES is an array
    configFiles+=("${BASH_FRAMEWORK_ENV_FILES[@]}")
  fi
  local localFrameworkConfigFile
  localFrameworkConfigFile="$(pwd)/.framework-config"
  if [[ -f "${localFrameworkConfigFile}" ]]; then
    configFiles+=("${localFrameworkConfigFile}")
  fi
  if [[ -f "${FRAMEWORK_ROOT_DIR}/.framework-config" ]]; then
    configFiles+=("${FRAMEWORK_ROOT_DIR}/.framework-config")
  fi
  configFiles+=("${optionEnvFiles[@]}")
  configFiles+=("${defaultFiles[@]}")

  for file in "${configFiles[@]}"; do
    # shellcheck source=/.framework-config
    CURRENT_LOADED_ENV_FILE="${file}" source "${file}" || {
      Log::displayError "while loading config file: ${file}"
      return 1
    }
  done
}
#!/usr/bin/env bash

# @description create a temp file using default TMPDIR variable
# initialized in _includes/_commonHeader.sh
# @env TMPDIR String (default value /tmp)
# @arg $1 templateName:String template name to use(optional)
Framework::createTempFile() {
  mktemp -p "${TMPDIR:-/tmp}" -t "${1:-}.XXXXXXXXXXXX"
}
#!/usr/bin/env bash

declare -g FIRST_LOG_DATE LOG_LAST_LOG_DATE LOG_LAST_LOG_DATE_INIT LOG_LAST_DURATION_STR
FIRST_LOG_DATE="${EPOCHREALTIME/[^0-9]/}"
LOG_LAST_LOG_DATE="${FIRST_LOG_DATE}"
LOG_LAST_LOG_DATE_INIT=1
LOG_LAST_DURATION_STR=""

# @description compute duration since last call to this function
# the result is set in following env variables.
# in ss.sss (seconds followed by milliseconds precision 3 decimals)
# @noargs
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @set LOG_LAST_LOG_DATE_INIT int (default 1) set to 0 at first call, allows to detect reference log
# @set LOG_LAST_DURATION_STR String the last duration displayed
# @set LOG_LAST_LOG_DATE String the last log date that will be used to compute next diff
Log::computeDuration() {
  if ((${DISPLAY_DURATION:-0} == 1)); then
    local -i duration=0
    local -i delta=0
    local -i currentLogDate
    currentLogDate="${EPOCHREALTIME/[^0-9]/}"
    if ((LOG_LAST_LOG_DATE_INIT == 1)); then
      LOG_LAST_LOG_DATE_INIT=0
      LOG_LAST_DURATION_STR="Ref"
    else
      duration=$(((currentLogDate - FIRST_LOG_DATE) / 1000000))
      delta=$(((currentLogDate - LOG_LAST_LOG_DATE) / 1000000))
      LOG_LAST_DURATION_STR="${duration}s/+${delta}s"
    fi
    LOG_LAST_LOG_DATE="${currentLogDate}"
    # shellcheck disable=SC2034
    local microSeconds="${EPOCHREALTIME#*.}"
    LOG_LAST_DURATION_STR="$(printf '%(%T)T.%03.0f\n' "${EPOCHSECONDS}" "${microSeconds:0:3}")(${LOG_LAST_DURATION_STR}) - "
  else
    # shellcheck disable=SC2034
    LOG_LAST_DURATION_STR=""
  fi
}
#!/usr/bin/env bash

# @description Display message using debug color (gray)
# @arg $1 message:String the message to display
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @env LOG_CONTEXT String allows to contextualize the log
Log::displayDebug() {
  if ((BASH_FRAMEWORK_DISPLAY_LEVEL >= __LEVEL_DEBUG)); then
    Log::computeDuration
    echo -e "${__DEBUG_COLOR}DEBUG   - ${LOG_CONTEXT:-}${LOG_LAST_DURATION_STR:-}${1}${__RESET_COLOR}" >&2
  fi
  Log::logDebug "$1"
}
#!/usr/bin/env bash

# @description Display message using error color (red)
# @arg $1 message:String the message to display
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @env LOG_CONTEXT String allows to contextualize the log
Log::displayError() {
  if ((BASH_FRAMEWORK_DISPLAY_LEVEL >= __LEVEL_ERROR)); then
    Log::computeDuration
    echo -e "${__ERROR_COLOR}ERROR   - ${LOG_CONTEXT:-}${LOG_LAST_DURATION_STR:-}${1}${__RESET_COLOR}" >&2
  fi
  Log::logError "$1"
}
#!/usr/bin/env bash

# @description Display message using info color (bg light blue/fg white)
# @arg $1 message:String the message to display
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @env LOG_CONTEXT String allows to contextualize the log
Log::displayInfo() {
  local type="${2:-INFO}"
  if ((BASH_FRAMEWORK_DISPLAY_LEVEL >= __LEVEL_INFO)); then
    Log::computeDuration
    echo -e "${__INFO_COLOR}${type}    - ${LOG_CONTEXT:-}${LOG_LAST_DURATION_STR:-}${1}${__RESET_COLOR}" >&2
  fi
  Log::logInfo "$1" "${type}"
}
#!/usr/bin/env bash

# @description Display message using warning color (yellow)
# @arg $1 message:String the message to display
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @env LOG_CONTEXT String allows to contextualize the log
Log::displayWarning() {
  if ((BASH_FRAMEWORK_DISPLAY_LEVEL >= __LEVEL_WARNING)); then
    Log::computeDuration
    echo -e "${__WARNING_COLOR}WARN    - ${LOG_CONTEXT:-}${LOG_LAST_DURATION_STR:-}${1}${__RESET_COLOR}" >&2
  fi
  Log::logWarning "$1"
}
#!/usr/bin/env bash

# @description Display message using error color (red) and exit immediately with error status 1
# @arg $1 message:String the message to display
# @env DISPLAY_DURATION int (default 0) if 1 display elapsed time information between 2 info logs
# @env LOG_CONTEXT String allows to contextualize the log
Log::fatal() {
  Log::computeDuration
  echo -e "${__ERROR_COLOR}FATAL   - ${LOG_CONTEXT:-}${LOG_LAST_DURATION_STR:-}${1}${__RESET_COLOR}" >&2
  Log::logFatal "$1"
  exit 1
}
#!/usr/bin/env bash

# @description log message to file
# @arg $1 message:String the message to display
Log::logDebug() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_DEBUG)); then
    Log::logMessage "${2:-DEBUG}" "$1"
  fi
}
#!/usr/bin/env bash

# @description log message to file
# @arg $1 message:String the message to display
Log::logError() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_ERROR)); then
    Log::logMessage "${2:-ERROR}" "$1"
  fi
}
#!/usr/bin/env bash

# @description log message to file
# @arg $1 message:String the message to display
Log::logFatal() {
  Log::logMessage "${2:-FATAL}" "$1"
}
#!/usr/bin/env bash

# @description log message to file
# @arg $1 message:String the message to display
Log::logInfo() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_INFO)); then
    Log::logMessage "${2:-INFO}" "$1"
  fi
}
#!/usr/bin/env bash

# @description Internal: common log message
# @example text
#   [date]|[levelMsg]|message
#
# @example text
#   2020-01-19 19:20:21|ERROR  |log error
#   2020-01-19 19:20:21|SKIPPED|log skipped
#
# @arg $1 levelMsg:String message's level description (eg: STATUS, ERROR, ...)
# @arg $2 msg:String the message to display
# @env BASH_FRAMEWORK_LOG_FILE String log file to use, do nothing if empty
# @env BASH_FRAMEWORK_LOG_LEVEL int log level log only if > OFF or fatal messages
# @stderr diagnostics information is displayed
# @require Env::requireLoad
# @require Log::requireLoad
Log::logMessage() {
  local levelMsg="$1"
  local msg="$2"
  local date

  if [[ -n "${BASH_FRAMEWORK_LOG_FILE}" ]] && ((BASH_FRAMEWORK_LOG_LEVEL > __LEVEL_OFF)); then
    date="$(date '+%Y-%m-%d %H:%M:%S')"
    touch "${BASH_FRAMEWORK_LOG_FILE}"
    printf "%s|%7s|%s\n" "${date}" "${levelMsg}" "${msg}" >>"${BASH_FRAMEWORK_LOG_FILE}"
  fi
}
#!/usr/bin/env bash

# @description log message to file
# @arg $1 message:String the message to display
Log::logWarning() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_WARNING)); then
    Log::logMessage "${2:-WARNING}" "$1"
  fi
}
#!/usr/bin/env bash

# @description activate or not Log::display* and Log::log* functions
# based on BASH_FRAMEWORK_DISPLAY_LEVEL and BASH_FRAMEWORK_LOG_LEVEL
# environment variables loaded by Env::requireLoad
# try to create log file and rotate it if necessary
# @noargs
# @set BASH_FRAMEWORK_LOG_LEVEL int to OFF level if BASH_FRAMEWORK_LOG_FILE is empty or not writable
# @env BASH_FRAMEWORK_DISPLAY_LEVEL int
# @env BASH_FRAMEWORK_LOG_LEVEL int
# @env BASH_FRAMEWORK_LOG_FILE String
# @env BASH_FRAMEWORK_LOG_FILE_MAX_ROTATION int do log rotation if > 0
# @exitcode 0 always successful
# @stderr diagnostics information about log file is displayed
# @require Env::requireLoad
# @require UI::requireTheme
Log::requireLoad() {
  if [[ -z "${BASH_FRAMEWORK_LOG_FILE:-}" ]]; then
    BASH_FRAMEWORK_LOG_LEVEL=${__LEVEL_OFF}
    export BASH_FRAMEWORK_LOG_LEVEL
  fi

  if ((BASH_FRAMEWORK_LOG_LEVEL > __LEVEL_OFF)); then
    if [[ ! -f "${BASH_FRAMEWORK_LOG_FILE}" ]]; then
      if [[ ! -d "${BASH_FRAMEWORK_LOG_FILE%/*}" ]]; then
        if ! mkdir -p "${BASH_FRAMEWORK_LOG_FILE%/*}" 2>/dev/null; then
          BASH_FRAMEWORK_LOG_LEVEL=${__LEVEL_OFF}
          echo -e "${__ERROR_COLOR}ERROR   - directory ${BASH_FRAMEWORK_LOG_FILE%/*} is not writable${__RESET_COLOR}" >&2
        fi
      elif ! touch --no-create "${BASH_FRAMEWORK_LOG_FILE}" 2>/dev/null; then
        BASH_FRAMEWORK_LOG_LEVEL=${__LEVEL_OFF}
        echo -e "${__ERROR_COLOR}ERROR   - File ${BASH_FRAMEWORK_LOG_FILE} is not writable${__RESET_COLOR}" >&2
      fi
    elif [[ ! -w "${BASH_FRAMEWORK_LOG_FILE}" ]]; then
      BASH_FRAMEWORK_LOG_LEVEL=${__LEVEL_OFF}
      echo -e "${__ERROR_COLOR}ERROR   - File ${BASH_FRAMEWORK_LOG_FILE} is not writable${__RESET_COLOR}" >&2
    fi
  fi

  if ((BASH_FRAMEWORK_LOG_LEVEL > __LEVEL_OFF)); then
    # will always be created even if not in info level
    Log::logMessage "INFO" "Logging to file ${BASH_FRAMEWORK_LOG_FILE} - Log level ${BASH_FRAMEWORK_LOG_LEVEL}"
    if ((BASH_FRAMEWORK_LOG_FILE_MAX_ROTATION > 0)); then
      Log::rotate "${BASH_FRAMEWORK_LOG_FILE}" "${BASH_FRAMEWORK_LOG_FILE_MAX_ROTATION}"
    fi
  fi
}
#!/usr/bin/env bash

# @description To be called before logging in the log file
# @arg $1 file:string log file name
# @arg $2 maxLogFilesCount:int maximum number of log files
Log::rotate() {
  local file="$1"
  local maxLogFilesCount="${2:-5}"

  if [[ ! -f "${file}" ]]; then
    Log::displayDebug "Log file ${file} doesn't exist yet"
    return 0
  fi
  local i
  for ((i = maxLogFilesCount - 1; i > 0; i--)); do
    Log::displayInfo "Log rotation ${file}.${i} to ${file}.$((i + 1))"
    mv "${file}."{"${i}","$((i + 1))"} &>/dev/null || true
  done
  if cp "${file}" "${file}.1" &>/dev/null; then
    echo >"${file}" # reset log file
    Log::displayInfo "Log rotation ${file} to ${file}.1"
  fi
}
#!/usr/bin/env bash

# @description draw a line with the character passed in parameter repeated depending on terminal width
# @arg $1 character:String character to use as separator (default value #)
UI::drawLine() {
  local character="${1:-#}"
  local -i width=${COLUMNS:-0}
  if ((width == 0)) && [[ -t 1 ]]; then
    width=$(tput cols)
  fi
  if ((width == 0)); then
    width=80
  fi
  printf -- "${character}%.0s" $(seq "${COLUMNS:-$([[ -t 1 ]] && tput cols || echo '80')}")
  echo
}
#!/usr/bin/env bash

# @description load colors theme constants
# @warning if tty not opened, noColor theme will be chosen
# @arg $1 theme:String the theme to use (default, noColor)
# @arg $@ args:String[]
# @set __ERROR_COLOR String indicate error status
# @set __INFO_COLOR String indicate info status
# @set __SUCCESS_COLOR String indicate success status
# @set __WARNING_COLOR String indicate warning status
# @set __SKIPPED_COLOR String indicate skipped status
# @set __DEBUG_COLOR String indicate debug status
# @set __HELP_COLOR String indicate help status
# @set __TEST_COLOR String not used
# @set __TEST_ERROR_COLOR String not used
# @set __HELP_TITLE_COLOR String used to display help title in help strings
# @set __HELP_OPTION_COLOR String used to display highlight options in help strings
#
# @set __RESET_COLOR String reset default color
#
# @set __HELP_EXAMPLE String to remove
# @set __HELP_TITLE String to remove
# @set __HELP_NORMAL String to remove
# shellcheck disable=SC2034
UI::theme() {
  local theme="${1-default}"
  if [[ ! "${theme}" =~ -force$ ]] && ! Assert::tty; then
    theme="noColor"
  fi
  case "${theme}" in
    default | default-force)
      theme="default"
      ;;
    noColor) ;;
    *)
      Log::fatal "invalid theme provided"
      ;;
  esac
  if [[ "${theme}" = "default" ]]; then
    BASH_FRAMEWORK_THEME="default"
    # check colors applicable https://misc.flogisoft.com/bash/tip_colors_and_formatting
    __ERROR_COLOR='\e[31m'         # Red
    __INFO_COLOR='\e[44m'          # white on lightBlue
    __SUCCESS_COLOR='\e[32m'       # Green
    __WARNING_COLOR='\e[33m'       # Yellow
    __SKIPPED_COLOR='\e[33m'       # Yellow
    __DEBUG_COLOR='\e[37m'         # Gray
    __HELP_COLOR='\e[7;49;33m'     # Black on Gold
    __TEST_COLOR='\e[100m'         # Light magenta
    __TEST_ERROR_COLOR='\e[41m'    # white on red
    __HELP_TITLE_COLOR="\e[1;37m"  # Bold
    __HELP_OPTION_COLOR="\e[1;34m" # Blue
    # Internal: reset color
    __RESET_COLOR='\e[0m' # Reset Color
    # shellcheck disable=SC2155,SC2034
    __HELP_EXAMPLE="$(echo -e "\e[2;97m")"
    # shellcheck disable=SC2155,SC2034
    __HELP_TITLE="$(echo -e "\e[1;37m")"
    # shellcheck disable=SC2155,SC2034
    __HELP_NORMAL="$(echo -e "\033[0m")"
  else
    BASH_FRAMEWORK_THEME="noColor"
    # check colors applicable https://misc.flogisoft.com/bash/tip_colors_and_formatting
    __ERROR_COLOR=''
    __INFO_COLOR=''
    __SUCCESS_COLOR=''
    __WARNING_COLOR=''
    __SKIPPED_COLOR=''
    __DEBUG_COLOR=''
    __HELP_COLOR=''
    __TEST_COLOR=''
    __TEST_ERROR_COLOR=''
    __HELP_TITLE_COLOR=''
    __HELP_OPTION_COLOR=''
    # Internal: reset color
    __RESET_COLOR=''
    __HELP_EXAMPLE=''
    __HELP_TITLE=''
    __HELP_NORMAL=''
  fi
}
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
  ${commandFunctionName} <% % >help
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
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_INFO}" >>"${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionDebugVerboseCallback() {
  BASH_FRAMEWORK_ARGS_VERBOSE_OPTION='-vv'
  BASH_FRAMEWORK_ARGS_VERBOSE=${__VERBOSE_LEVEL_DEBUG}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_DEBUG}" >>"${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionTraceVerboseCallback() {
  BASH_FRAMEWORK_ARGS_VERBOSE_OPTION='-vvv'
  BASH_FRAMEWORK_ARGS_VERBOSE=${__VERBOSE_LEVEL_TRACE}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${__LEVEL_DEBUG}" >>"${overrideEnvFile}"
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
      ;;
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
      ;;
  esac
}

# shellcheck disable=SC2317 # if function is overridden
optionDisplayLevelCallback() {
  local level="$2"
  local logLevel verboseLevel
  logLevel="$(getLevel "${level}")"
  verboseLevel="$(getVerboseLevel "${level}")"
  BASH_FRAMEWORK_ARGS_VERBOSE=${verboseLevel}
  echo "BASH_FRAMEWORK_DISPLAY_LEVEL=${logLevel}" >>"${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionLogLevelCallback() {
  local level="$2"
  local logLevel verboseLevel
  logLevel="$(getLevel "${level}")"
  verboseLevel="$(getVerboseLevel "${level}")"
  BASH_FRAMEWORK_ARGS_VERBOSE=${verboseLevel}
  echo "BASH_FRAMEWORK_LOG_LEVEL=${logLevel}" >>"${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionLogFileCallback() {
  local logFile="$2"
  echo "BASH_FRAMEWORK_LOG_FILE='${logFile}'" >>"${overrideEnvFile}"
}

# shellcheck disable=SC2317 # if function is overridden
optionQuietCallback() {
  echo "BASH_FRAMEWORK_QUIET_MODE=1" >>"${overrideEnvFile}"
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
  echo "${defaultFrameworkConfig}" >"${defaultEnvFile}"
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

        if ((options_parse_optionParsedCountOptionFormat >= 1)); then
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

        if ((options_parse_optionParsedCountOptionStaged >= 1)); then
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

        if ((options_parse_optionParsedCountOptionXargs >= 1)); then
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

        if ((options_parse_optionParsedCountOptionHelp >= 1)); then
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

        if ((options_parse_optionParsedCountOptionConfig >= 1)); then
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

        if ((options_parse_optionParsedCountOptionBashFrameworkConfig >= 1)); then
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

        if ((options_parse_optionParsedCountOptionInfoVerbose >= 1)); then
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

        if ((options_parse_optionParsedCountOptionDebugVerbose >= 1)); then
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

        if ((options_parse_optionParsedCountOptionTraceVerbose >= 1)); then
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

        if ((options_parse_optionParsedCountOptionLogLevel >= 1)); then
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

        if ((options_parse_optionParsedCountOptionLogFile >= 1)); then
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

        if ((options_parse_optionParsedCountOptionDisplayLevel >= 1)); then
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

        if ((options_parse_optionParsedCountOptionNoColor >= 1)); then
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

        if ((options_parse_optionParsedCountOptionTheme >= 1)); then
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

        if ((options_parse_optionParsedCountOptionVersion >= 1)); then
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

        if ((options_parse_optionParsedCountOptionQuiet >= 1)); then
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

        elif ((options_parse_parsedArgIndex >= 0)); then

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
  done
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
  echo '${REPOSITORY_URL}/tree/master/${SRC_FILE_PATH}'
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
  shellcheckLintCommandParse "$@"

}

# if file is sourced avoid calling main function
# shellcheck disable=SC2178
BASH_SOURCE=".$0" # cannot be changed in bash
# shellcheck disable=SC2128
test ".$0" != ".${BASH_SOURCE}" || main "$@"