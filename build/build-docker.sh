#!/usr/bin/env bash
set -eo pipefail

declare image="scrasnups/bash-compiler"
mkdir -pv logs bin

docker buildx build -t "${image}" .
