module github.com/fchastanet/bash-compiler

go 1.26.8

require (
	github.com/Masterminds/sprig/v3 v3.3.0
	github.com/a8m/envsubst v1.4.3
	github.com/alecthomas/kong v1.16.1
	github.com/bmatcuk/doublestar/v4 v4.10.0
	github.com/goccy/go-yaml v1.19.2
	github.com/stretchr/testify v1.12.1
	go.uber.org/automaxprocs v1.6.0
	golang.org/x/text v0.42.0
	gotest.tools/v3 v3.5.2
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/chai2010/jsonv v1.1.3 // indirect
	github.com/chai2010/protorpc v1.1.4 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/gofrs/flock v0.13.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260911204522-f61a6ca850bd // indirect
	google.golang.org/grpc v1.85.0-dev.0.20260825072537-93e31b48545e // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	kcl-lang.io/lib v0.12.5 // indirect
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.6.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	kcl-lang.io/kcl-go v0.12.5
)

// kcl-lang.io/lib v0.12.5 regressed the source location reported in
// EvaluationError messages: type errors point at an unrelated line of the
// schema instead of the offending attribute declaration (e.g. a bad
// `mainFile` value is reported on the `_optionGroups` comprehension).
// The `expect X, got Y` part stays correct, but the snippet is what
// internal/model/transformModel_test.go relies on to know which attribute
// failed. kcl-go v0.12.5 with lib v0.12.4 behaves correctly.
// Drop this replace once the location reporting is fixed upstream.
replace kcl-lang.io/lib v0.12.5 => kcl-lang.io/lib v0.12.4
