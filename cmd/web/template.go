package main

import (
	"html/template"
	"log"
	"path/filepath"
	"time"

	"todo.khoirulakmal.dev/internal/models"
)

type templateData struct {
	Year  int
	Lists *[]models.List
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
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		log.Print(name)

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/dynamic/status.html",
			"./ui/html/dynamic/data.html",
			"./ui/html/dynamic/list.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		templateCache[name] = ts

	}

	return templateCache, nil
}
