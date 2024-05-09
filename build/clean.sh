#!/usr/bin/env bash
set -eo pipefail

echo "Cleaning ..."
rm -rvf bin logs || true
go mod tidy || true
docker image rm -f scrasnups/bash-compiler || true
