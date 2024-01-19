package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	data.Lists = &result
	app.render(w, "base", data)
}

// Add a snippetCreate handler function.
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	content := r.PostForm.Get("list")
	status := "ongoing"
	id, err := app.todos.Insert(content, status)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Print(err)
		return
	}
	w.Header().Add("id", strconv.Itoa(id))
	app.getLists(w, r)

}

func (app *application) getLists(w http.ResponseWriter, r *http.Request) {
	list, err := app.todos.GetRows()
	if err != nil {
		app.serverError(w, err)
		return
	}
	messages := app.session.GetString(r.Context(), "Message")
	w.Header().Add("Sessions", messages)
	data := app.generateTemplateData()
	data.Lists = &list
	app.render(w, "data", data)
}

func (app *application) getList(w http.ResponseWriter, r *http.Request, id int) {
	list, err := app.todos.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.generateTemplateData()
	data.List = list
	app.infoLog.Print(data.List.Status)
	app.render(w, "data", data)
}

func (app *application) deleteList(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	deleteList, err := app.todos.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}
	if deleteList != -1 {
		app.getLists(w, r)
	}
}

func (app *application) updateStatus(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We'll talk about request context
	// in detail later in the book, but for now it's enough to know that you can
	// use the ParamsFromContext() function to retrieve a slice containing these
	// parameter names and values like so:
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	doneList, err := app.todos.Done(id)
	app.session.Put(r.Context(), "Message", "Yeah its me baby")
	if err != nil {
		app.serverError(w, err)
	}
	if doneList != -1 {
		app.getLists(w, r)
	}
}
