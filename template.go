package gooseglass

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed *.gohtml
var templateFiles embed.FS

//go:generate go run github.com/typelate/muxt generate  --receiver-type=Provider --receiver-type-package=github.com/pressly/goose/v3 --receiver-interface=Provider --routes-func routes --template-data-type templateData
var templates = template.Must(template.ParseFS(templateFiles, "*"))

func Pages(mux *http.ServeMux, provider Provider) { routes(mux, provider) }

func (td *templateData[T]) TriggerRefreshMigrations() *templateData[T] {
	return td.Header("HX-Trigger", `{"refreshMigrations":{"target":"#status-table"}}`)
}
