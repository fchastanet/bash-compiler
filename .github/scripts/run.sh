#!/usr/bin/env bash
set -e -o pipefail -o errexit

go run ./cmd/bash-compiler/main.go
