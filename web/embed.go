// Package web contains the embedded static files for the web UI.
// These files are compiled into the binary at build time using Go's embed package.
package web

import "embed"

// StaticFiles embeds all files in the static/ directory.
// The //go:embed directive tells the Go compiler to include these files
// in the binary at compile time. This means:
//   - No external files needed at runtime
//   - The binary is self-contained
//   - Deployment is simpler (just one file)
//
// The embed.FS type provides a read-only filesystem interface
// to access these embedded files.
//
//go:embed static/*
var StaticFiles embed.FS
