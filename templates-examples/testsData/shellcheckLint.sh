#!/bin/bash
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
