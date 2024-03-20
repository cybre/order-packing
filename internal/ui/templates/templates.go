package templates

import (
	"embed"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/unrolled/render"
)

//go:embed views/*.html
var embeddedTemplates embed.FS

type Renderer struct {
	r *render.Render
}

func New() *Renderer {
	r := render.New(render.Options{
		Directory: "views",
		FileSystem: &render.EmbedFileSystem{
			FS: embeddedTemplates,
		},
		Extensions: []string{".html"},
		Funcs:      []template.FuncMap{},
	})

	return &Renderer{
		r: r,
	}
}

func (t *Renderer) Render(w io.Writer, name string, pageData interface{}, c echo.Context) error {
	return t.r.HTML(w, http.StatusOK, name, pageData)
}
