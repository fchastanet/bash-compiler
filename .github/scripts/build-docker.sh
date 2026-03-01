#!/usr/bin/env bash
set -e -o pipefail -o errexit

declare image="scrasnups/bash-compiler"
mkdir -pv logs bin

docker buildx build -t "${image}" .
