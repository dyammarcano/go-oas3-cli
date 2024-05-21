package embed

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var assets embed.FS

// GetAssets returns the embedded assets from the dist directory.
func GetAssets() fs.FS {
	files, _ := fs.Sub(assets, "dist")
	return files
}
