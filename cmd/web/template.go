package main

import (
	"html/template"
	"log"
	"path/filepath"

	"todo.khoirulakmal.dev/internal/models"
)

type templateData struct {
	Year  int
	Lists *[]models.List
	List  *models.List
	Form  any
	Flash string
	Page  string
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
			"./ui/html/partials/data.html",
			page,
		}
		ts, err := template.New(name).ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		templateCache[name] = ts

	}

	return templateCache, nil
}
