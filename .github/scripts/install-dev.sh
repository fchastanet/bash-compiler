#!/usr/bin/env bash
set -euxo pipefail

go install mvdan.cc/gofumpt@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/mgechev/revive@latest
go install github.com/boumenot/gocover-cobertura@latest
