//go:build !embed_gui

package gui

import "embed"

var DistFS embed.FS

var Embedded = false
