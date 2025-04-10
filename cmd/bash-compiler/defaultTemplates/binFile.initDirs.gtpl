#!/usr/bin/env bash

SCRIPT_NAME=${0##*/}
REAL_SCRIPT_FILE="$(readlink -e "$(realpath "${BASH_SOURCE[0]}")")"
if [[ -n "${EMBED_CURRENT_DIR}" ]]; then
  CURRENT_DIR="${EMBED_CURRENT_DIR}"
else
  CURRENT_DIR="${REAL_SCRIPT_FILE%/*}"
fi
{{ if .RootData.compilerConfig.RelativeRootDirBasedOnTargetDir -}}
FRAMEWORK_ROOT_DIR="$(cd "${CURRENT_DIR}/{{- .RootData.compilerConfig.RelativeRootDirBasedOnTargetDir -}}" && pwd -P)"
FRAMEWORK_SRC_DIR="${FRAMEWORK_ROOT_DIR}/src"
FRAMEWORK_BIN_DIR="${FRAMEWORK_ROOT_DIR}/bin"
FRAMEWORK_VENDOR_DIR="${FRAMEWORK_ROOT_DIR}/vendor"
FRAMEWORK_VENDOR_BIN_DIR="${FRAMEWORK_ROOT_DIR}/vendor/bin"
{{ end -}}
