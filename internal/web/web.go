package web

import (
	"embed"
	"io/fs"
)

//go:embed templates/**.html
var templates embed.FS

// Templates contain all the template files from the templates directory.
var Templates, _ = fs.Sub(templates, "templates")
