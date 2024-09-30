package csvimport

import "embed"

// Embed a single file
//go:embed dist/*
var F embed.FS
