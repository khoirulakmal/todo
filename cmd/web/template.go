package main

import (
	"html/template"
	"path/filepath"

	"todo.khoirulakmal.dev/internal/models"
)

type templateData struct {
	Year int
	List *[]models.List
}

func parseTemplate() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)
	pages, err := filepath.Glob("ui/html/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseGlob("ui/html/*.html")
		if err != nil {
			return nil, err
		}
		templateCache[name] = ts
	}

	return templateCache, nil
}
