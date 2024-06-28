#!/bin/bash
set -o errexit

REAL_SCRIPT_FILE="$(readlink -e "$(realpath "${BASH_SOURCE[0]}")")"
CURRENT_DIR="${REAL_SCRIPT_FILE%/*}"
ROOT_DIR="${REAL_SCRIPT_FILE%/*/*}"

(
  cd "${ROOT_DIR}" || exit 1

  go install github.com/jfeliu007/goplantuml/cmd/goplantuml@latest
  goplantuml \
    -recursive -aggregate-private-members -show-compositions \
    -show-aliases -show-aggregations -show-connection-labels \
    -show-options-as-note -hide-private-members -ignore "builtin" . \
    >"${CURRENT_DIR}/classDiagram.puml"
  goplantuml \
    -recursive -aggregate-private-members -show-compositions \
    -show-aliases -show-aggregations -show-connection-labels \
    -show-options-as-note -ignore "builtin" . \
    >"${CURRENT_DIR}/classDiagramWithPrivateMethods.puml"
)
