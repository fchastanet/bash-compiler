#!/usr/bin/env bash

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

