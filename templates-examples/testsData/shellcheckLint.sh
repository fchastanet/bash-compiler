#!/usr/bin/env bash


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
  local -i options_parse_argParsedCountArgShellcheckFiles
  ((options_parse_argParsedCountArgShellcheckFiles = 0)) || true
  
  
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
      # authorizedValues:
      --format | -f)
        
        shift
        if (($# == 0)); then
          Log::displayError "Command ${SCRIPT_NAME} - Option ${options_parse_arg} - a value needs to be specified"
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
        # argShellcheckFiles min 1 max -1
        
        
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
  
  
  
  
  if ((options_parse_argParsedCountArgShellcheckFiles < 1 )); then
    Log::displayError "Command ${SCRIPT_NAME} - Argument 'shellcheckFiles' should be provided at least 1 time(s)"
    return 1
  fi
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
    "Lint bash files using shellcheck.
  "
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
    "[--format|-f <optionFormat>]"
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
  echo -e "${__HELP_OPTION_COLOR}--format${__HELP_NORMAL}, ${__HELP_OPTION_COLOR}-f optionFormat${__HELP_NORMAL} {single}"
  
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
}

