package template

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSwaggerHtmlData(t *testing.T) {
	version := "5.17.10"
	swaggerJson := "https://petstore3.swagger.io/api/v3/openapi.json"
	fs := afero.NewOsFs()

	err := fs.MkdirAll("swagger-ui", 0755)
	assert.NoError(t, err)

	got := NewSwaggerHtmlData(fs, version, swaggerJson)

	err = got.DownloadData()
	assert.NoError(t, err)

	err = got.GenerateSwaggerHtml()
	assert.NoError(t, err)
}
