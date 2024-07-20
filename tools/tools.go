//go:build tools

package tools

// This is a list of tools to be maintained; some non-main
// package in the same repo needs to be imported so they'll be managed in
// go.mod.

import (
	_ "fyne.io/fyne/v2/cmd/fyne"
	_ "github.com/abice/go-enum"
	_ "github.com/fyne-io/fyne-cross"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/pprof"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/tomwright/dasel/cmd/dasel"
	_ "golang.org/x/tools/cmd/godoc"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"
	_ "mvdan.cc/gofumpt"
)
