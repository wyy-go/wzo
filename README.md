[![GoDoc](https://godoc.org/github.com/wyy-go/go-cli-template?status.svg)](https://godoc.org/github.com/wyy-go/go-cli-template)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wyy-go/go-cli-template?tab=doc)
[![codecov](https://codecov.io/gh/wyy-go/go-cli-template/branch/main/graph/badge.svg)](https://codecov.io/gh/wyy-go/go-cli-template)
[![Tests](https://github.com/wyy-go/go-cli-template/actions/workflows/ci.yaml/badge.svg)](https://github.com/wyy-go/go-cli-template/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wyy-go/go-cli-template)](https://goreportcard.com/report/github.com/wyy-go/go-cli-template)
[![Licence](https://img.shields.io/github/license/wyy-go/go-cli-template)](https://raw.githubusercontent.com/wyy-go/go-cli-template/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/wyy-go/go-cli-template)](https://github.com/wyy-go/go-cli-template/tags)

# go-cli-template
This is template that help you to quick implement some CLI using Go.

This repository is contains following.

- minimal CLI implementation using [spf13/cobra](https://github.com/spf13/cobra)
- CI/CD
  - [golangci-lint](https://golangci-lint.run/usage/linters/)
  - go test
  - goreleaser
  - dependabot for github-actions and Go
  - CodeQL Analysis (Go)

## How to use
1. fork this repository
2. replace `wyy-go` to your user name using `sed`(or others)
3. run `make init`

## Author
wyy-go

## References

- [go-cli-template](https://github.com/skanehira/go-cli-template)

