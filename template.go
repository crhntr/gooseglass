package gooseglass

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed *.gohtml
var templateFiles embed.FS

//go:generate go run github.com/typelate/muxt generate --use-receiver-type=Provider --use-receiver-type-package=github.com/pressly/goose/v3 --output-receiver-interface=Provider --output-routes-func routes --output-template-data-type templateData
var templates = template.Must(template.ParseFS(templateFiles, "*"))

func Pages(mux *http.ServeMux, provider Provider) { routes(mux, provider) }

func (td *templateData[R, T]) TriggerRefreshMigrations() *templateData[R, T] {
	return td.Header("HX-Trigger", `{"refreshMigrations":{"target":"#status-table"}}`)
}
