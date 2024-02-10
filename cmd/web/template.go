package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"todo.khoirulakmal.dev/internal/models"
	"todo.khoirulakmal.dev/ui"
)

type templateData struct {
	Year  int
	Lists []*models.List
	List  *models.List
	Form  any
	Flash string
	Page  string
	Auth  bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 02:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func parseTemplate() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := []string{
			"html/base.html",
			"html/partials/nav.html",
			"html/dynamic/status.html",
			"html/dynamic/data.html",
			"html/dynamic/list.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, files...)
		if err != nil {
			return nil, err
		}
		templateCache[name] = ts

	}

	return templateCache, nil
}
