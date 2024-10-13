#!/usr/bin/env bash
set -o pipefail -o errexit

mkdir -pv "${HOME}/go/bin" || true
go env -w "GOBIN=${HOME}/go/bin"

echo >&2 "Check dependencies ..."
go mod download

echo >&2 "Building ..."
go build -ldflags="-w -s" ./...

echo >&2 "Installing ..."
go install ./...

if [[ -f ${HOME}/go/bin/bash-compiler ]]; then
  echo >&2 "you can run ${HOME}/go/bin/bash-compiler"
else
  echo >&2 "${HOME}/go/bin/bash-compiler has not been generated"
fi
