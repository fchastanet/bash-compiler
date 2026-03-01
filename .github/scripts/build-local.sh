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

GO_BIN="$(go env GOBIN)"
if [[ -f ${GO_BIN}/bash-compiler ]]; then
  echo >&2 "you can run ${GO_BIN}/bash-compiler"
else
  echo >&2 "${GO_BIN}/bash-compiler has not been generated"
  exit 1
fi
