#!/usr/bin/env bash
###############################################################################
# GENERATED FROM ${REPOSITORY_URL}/tree/master/${SRC_FILE_PATH}
# DO NOT EDIT IT
# @generated
###############################################################################
# shellcheck disable=SC2288,SC2034
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
SCRIPT_NAME=${0##*/}
REAL_SCRIPT_FILE="$(readlink -e "$(realpath "${BASH_SOURCE[0]}")")"
if [[ -n "${EMBED_CURRENT_DIR}" ]]; then
  CURRENT_DIR="${EMBED_CURRENT_DIR}"
else
  CURRENT_DIR="${REAL_SCRIPT_FILE%/*}"
fi
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
# @description log message to file
# @arg $1 message:String the message to display
Log::logDebug() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_DEBUG)); then
    Log::logMessage "${2:-DEBUG}" "$1"
  fi
}
# @description log message to file
# @arg $1 message:String the message to display
Log::logError() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_ERROR)); then
    Log::logMessage "${2:-ERROR}" "$1"
  fi
}
# @description log message to file
# @arg $1 message:String the message to display
Log::logInfo() {
  if ((BASH_FRAMEWORK_LOG_LEVEL >= __LEVEL_INFO)); then
    Log::logMessage "${2:-INFO}" "$1"
  fi
}
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
# FUNCTIONS
# ------------------------------------------
# Command shellcheckLintCommand
# ------------------------------------------
# @description parse command options and arguments for shellcheckLintCommand
shellcheckLintCommandParse() {
  Log::displayDebug "Command ${SCRIPT_NAME} - parse arguments: ${BASH_FRAMEWORK_ARGV[*]}"
  Log::displayDebug "Command ${SCRIPT_NAME} - parse filtered arguments: ${BASH_FRAMEWORK_ARGV_FILTERED[*]}"
  optionHelp="0"
  local -i options_parse_optionParsedCountOptionHelp
  ((options_parse_optionParsedCountOptionHelp = 0)) || true
  optionFormat="tty"
  local -i options_parse_optionParsedCountOptionFormat
  ((options_parse_optionParsedCountOptionFormat = 0)) || true
  optionStaged="0"
  local -i options_parse_optionParsedCountOptionStaged
  ((options_parse_optionParsedCountOptionStaged = 0)) || true
  optionXargs="0"
  local -i options_parse_optionParsedCountOptionXargs
  ((options_parse_optionParsedCountOptionXargs = 0)) || true
  argShellcheckFiles=()
  # shellcheck disable=SC2034
  local -i options_parse_parsedArgIndex=0
  while (($# > 0)); do
    local options_parse_arg="$1"
    local argOptDefaultBehavior=0
    case "${options_parse_arg}" in
      # Option 1/4
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
      # Option 2/4
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
      # Option 3/4
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
      # Option 4/4
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
  done
}
shellcheckLintCommandLongDescription="$(cat <<'EOF'
shellcheck wrapper that will:
- install new shellcheck version(${MIN_SHELLCHECK_VERSION}) automatically
- by default, lint all git files of this project which are beginning with a bash shebang
  except if the option --staged is passed
${__HELP_TITLE}Special configuration .shellcheckrc:${__HELP_NORMAL}
use the following line in your .shellcheckrc file to exclude
some files from being checked (use grep -E syntax)
exclude=^bin/bash-tpl$
${__HELP_TITLE_COLOR}SHELLCHECK HELP${__RESET_COLOR}
@@@SHELLCHECK_HELP@@@
EOF
)"
# @description display command options and arguments help for shellcheckLintCommand
shellcheckLintCommandHelp() {
  Array::wrap2 ' ' 80 0 "${__HELP_TITLE_COLOR}DESCRIPTION:${__RESET_COLOR}" \
    "Lint bash files using shellcheck."
  echo
  # ------------------------------------------
  # usage section
  # ------------------------------------------
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" "shellcheckLint [OPTIONS] "
  # ------------------------------------------
  # usage/options section
  # ------------------------------------------
  optionsAltList=(
    "--help|-h"
    "[--format|-f <format>]"
    "--staged"
    "--xargs"
  )
  Array::wrap2 " " 80 2 "${__HELP_TITLE_COLOR}USAGE:${__RESET_COLOR}" \
    "shellcheckLint" "${optionsAltList[@]}"
  # ------------------------------------------
  # options section
  # ------------------------------------------
  echo
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"
  echo -e "${__HELP_OPTION_COLOR}--help${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-h${__HELP_NORMAL} {single}"
  echo
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"
  echo -e "${__HELP_OPTION_COLOR}--format${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-f format${__HELP_NORMAL} {single}"
  echo
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"
  echo -e "${__HELP_OPTION_COLOR}--staged${__HELP_NORMAL} {single}"
  echo
  echo -e "${__HELP_TITLE_COLOR}OPTIONS:${__RESET_COLOR}"
  echo -e "${__HELP_OPTION_COLOR}--xargs${__HELP_NORMAL} {single}"
  # ------------------------------------------
  # longDescription section
  # ------------------------------------------
  Array::wrap2 ' ' 76 0 "${shellcheckLintCommandLongDescription}"
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
