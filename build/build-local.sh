#!/usr/bin/env bash
set -e -o pipefail -o errexit

mkdir -pv bin || true

echo >&2 "Check dependencies ..."
go mod download

echo >&2 "Building ..."
go build -ldflags="-w -s" -o "bin" ./...
