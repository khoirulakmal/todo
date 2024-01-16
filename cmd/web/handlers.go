package main

import (
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	result, err := app.todos.GetRows()
	if err != nil {
		app.serverError(w, err)
		return
	}
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	data := app.generateTemplateData()
	data.List = &result
	app.render(w, "base", data)
}

// Add a snippetCreate handler function.
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}
	content := "Taking care of a cat"
	status := "ongoing"
	id, err := app.todos.Insert(content, status)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Print(err)
		return
	}
	w.Header().Add("id", strconv.Itoa(id))
	app.getList(w, r)

}

func (app *application) getList(w http.ResponseWriter, r *http.Request) {
	list, err := app.todos.GetRows()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.generateTemplateData()
	data.List = &list
	app.render(w, "data", data)
}
