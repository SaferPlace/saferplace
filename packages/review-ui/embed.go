package ui

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var dist embed.FS

// StaticFiles contains all the static files from the build directory.
var StaticFiles, _ = fs.Sub(dist, "dist")
