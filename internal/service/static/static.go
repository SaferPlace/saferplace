package static

import (
	"io/fs"
	"net/http"
)

// Register the static files handler which serves the frontend.
func Register(path string, fsys fs.FS) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return path, http.FileServer(http.FS(fsys))
	}
}
