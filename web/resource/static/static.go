//go:build gui

package static

import "embed"

//go:embed js/*
//go:embed css/*
var Embed embed.FS
