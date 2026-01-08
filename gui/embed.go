//go:build embed_gui

package gui

import "embed"

//go:embed all:dist
var DistFS embed.FS

var Embedded = true
