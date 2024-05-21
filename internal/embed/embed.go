package embed

import (
	"embed"
	"io/fs"
)

const (
	// SwaggerUIAssetsPath is the path for the generated assets.
	SwaggerUIAssetsPath = "swagger-ui"
)

//go:embed all:swagger-ui
var bundles embed.FS

// GetAssets returns the embedded assets from the dist directory.
func GetAssets() fs.FS {
	files, _ := fs.Sub(bundles, SwaggerUIAssetsPath)
	return files
}
