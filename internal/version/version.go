// Package version holds the binary version string.
//
// At link time, override with:
//
//	go build -ldflags "-X github.com/notwillk/nara/internal/version.Version=<value>" .
//
// Environment variables VERSION or VERSION are honored by `just build` when set.
// If not overridden, Version is "local-build".
package version

var Version = "local-build"
