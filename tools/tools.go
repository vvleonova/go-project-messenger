//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/matryer/moq"
)

var _ = true

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/matryer/moq
