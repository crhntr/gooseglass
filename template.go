package gooseglass

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/pressly/goose/v3"
)

//go:embed *.gohtml
var templateFiles embed.FS

//go:generate go run github.com/typelate/muxt generate  --receiver-type=Provider --receiver-type-package=github.com/pressly/goose/v3 --receiver-interface=migrationProvider --routes-func routes --template-data-type templateData
var templates = template.Must(template.ParseFS(templateFiles, "*"))

func Pages(mux *http.ServeMux, provider *goose.Provider) { routes(mux, provider) }
