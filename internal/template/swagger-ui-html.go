package template

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"text/template"
)

const (
	SwaggerHtmlTemplate = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="./swagger-ui.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="./swagger-ui-bundle.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '{{.SwaggerJson}}',
        dom_id: '#swagger-ui',
        deepLinking: true,
      });
    };
  </script>
  </body>
</html>
`

	Favicon16x16Png = `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAABNVBMVEVisTRhsTReqzVbpTVXoDdVnTdSlzhRljgvXkAuXUAtWkErV0EzZj40Zj85bz0lTkMkTUMkT0MmTUIkS0IjTEIhSUMkS0IkTEIkTUIlTUIkTkMlTkMcQUQcP0UfQ0QdQ0QfREQgRUMiSUMiSUMjSkInU0EkTEMmUEEiR0IiSEMpVkErWT8kTUElTUIUNkYVNEQVMkcRM0QSNUYQMUIMMUkVK0AAJEkAM00AMzMAAAAAAACF6i2E6SyD6CyC5i2B5Sx/4i6A4S593S583S520jB00DByyjFxyTFwyDFvxjJtxTFtxDFswzJrwDJqvzJpvjNouzNoujNnuDNLjTlKijpKiTpEfztDfzxAeT0+dz05bj44bT44bj82aj81aD8zZT8bPUUbPkUcP0UcPUUeQ0UfREQgRkRgJREvAAAAO3RSTlP09PX19vX39u7u7/Dq6ufh4eDg4+Pf3Nvb2tnY2NvPv7y6rKupqaGZlpSOiYWETDEkHh0fFQwHCgUBAAcHrskAAADYSURBVHjaPc/ZLkNRGIbhz26KjVJpqSKGtjHPc9a7W7OEEhtBjDWUO3XghqQSwVrNTp+j///OXhlrLpdJdg9MLblbxqwPd5RLUDpOjK66YWMwTqRpaM0OhZbo3dskljea9+HyAevxHtoWVAjhfQtr5w3CSfUE8BrgvEDQpxRc3eyfH5wenlQuIO39Sb9x/8uv+bXvmPSjbABPRZznIkGvxkOo7mJtV+FsQsutcFvBuruG9kWZMY+G5pzxlMp/KPKZSUs2cLrzyMWVEyP1OGtlNpvs6p+p5/8DzUo5hMDku9EAAAAASUVORK5CYII=`
	Favicon32x32Png = `iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAMAAABEpIrGAAAAkFBMVEUAAAAQM0QWNUYWNkYXNkYALjoWNUYYOEUXN0YaPEUPMUAUM0QVNUYWNkYWNUYWNUUWNUYVNEYWNkYWNUYWM0eF6i0XNkchR0OB5SwzZj9wyTEvXkA3az5apTZ+4C5DgDt31C9frjU5bz5uxTI/eDxzzjAmT0IsWUEeQkVltzR62S6D6CxIhzpKijpJiDpOkDl4b43lAAAAFXRSTlMAFc304QeZ/vj+ECB3xKlGilPXvS2Ka/h0AAABfklEQVR42oVT2XaCMBAdJRAi7pYJa2QHxbb//3ctSSAUPfa+THLmzj4DBvZpvyauS9b7kw3PWDkWsrD6fFQhQ9dZLfVbC5M88CWCPERr+8fLZodJ5M8QJbjbGL1H2M1fIGfEm+wJN+bGCSc6EXtNS/8FSrq2VX6YDv++XLpJ8SgDWMnwqznGo6alcTbIxB2CHKn8VFikk2mMV2lEnV+CJd9+jJlxXmMr5dW14YCqwgbFpO8FNvJxwwM4TPWPo5QalEsRMAcusXpi58/QUEWPL0AK1ThM5oQCUyXPoPINkdd922VBw4XgTV9zDGWWFrgjIQs4vwvOg6xr+6gbCTqE+DYhlMGX0CF2OknK5gQ2JrkDh/W6TOEbYDeVecKbJtyNXiCfGmW7V93J2hDus1bDfhxWbIZVYDXITA7Lo6E0Ktgg9eB4KWuR44aj7ppBVPazhQH7/M/KgWe9X1qAg8XypT6nxIMJH+T94QCsLvj29IYwZxyO9/F8vCbO9tX5/wDGjEZ7vrgFZwAAAABJRU5ErkJggg==`
)

type SwaggerHtmlData struct {
	version     string
	fs          afero.Fs
	SwaggerJson string
}

func NewSwaggerHtmlData(fs afero.Fs, version, swaggerJson string) SwaggerHtmlData {
	return SwaggerHtmlData{
		version:     version,
		fs:          fs,
		SwaggerJson: swaggerJson,
	}
}

func (d SwaggerHtmlData) DownloadData() error {
	// Download swagger-ui-bundle.js
	jsUrl := fmt.Sprintf("https://cdn.jsdelivr.net/npm/swagger-ui-dist@%s/swagger-ui-bundle.js", d.version)
	jsPath := "swagger-ui/swagger-ui-bundle.js"

	if err := downloadFile(d.fs, jsPath, jsUrl); err != nil {
		return fmt.Errorf("failed to download swagger-ui.css: %w", err)
	}

	// Download swagger-ui.css
	cssUrl := fmt.Sprintf("https://cdn.jsdelivr.net/npm/swagger-ui-dist@%s/swagger-ui.css", d.version)
	cssPath := "swagger-ui/swagger-ui.css"

	if err := downloadFile(d.fs, cssPath, cssUrl); err != nil {
		return fmt.Errorf("failed to download swagger-ui.css: %w", err)
	}

	decodeString, err := base64.StdEncoding.DecodeString(Favicon16x16Png)
	if err != nil {
		return err
	}

	if err = afero.WriteFile(d.fs, "swagger-ui/favicon-16x16.png", decodeString, 0644); err != nil {
		return err
	}

	decodeString, err = base64.StdEncoding.DecodeString(Favicon32x32Png)
	if err != nil {
		return err
	}

	if err = afero.WriteFile(d.fs, "swagger-ui/favicon-32x32.png", decodeString, 0644); err != nil {
		return err
	}
	return nil
}

func (d SwaggerHtmlData) GenerateSwaggerHtml() error {
	tmpl, err := template.New("swagger").Parse(SwaggerHtmlTemplate)
	if err != nil {
		return err
	}

	htmlPath := "swagger-ui/index.html"

	file, err := d.fs.Create(htmlPath)
	if err != nil {
		return err
	}

	return tmpl.Execute(file, d)
}

func downloadFile(fs afero.Fs, path, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := fs.Create(path)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}
