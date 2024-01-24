package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"todo.khoirulakmal.dev/internal/validator"
)

type todoForm struct {
	Content string
	Status  string
	validator.Validator
}

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
	data := app.generateTemplateData(r)
	data.Lists = &result
	data.Form = todoForm{
		Status: "ongo",
		Validator: validator.Validator{
			FieldErrors: nil,
		},
	}
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

	form := todoForm{
		Content: r.PostForm.Get("list"),
		Status:  "ongoing",
	}

	form.CheckField(validator.MinChars(form.Content, 5), "Content", "Content must be more than 5 characters")
	form.CheckField(validator.MaxChars(form.Content, 20), "Content", "Content must be less than 20 characters")

	if !form.Valid() {
		list, err := app.todos.GetRows()
		if err != nil {
			app.serverError(w, err)
			return
		}
		data := app.generateTemplateData(r)
		data.Form = form
		data.Lists = &list
		// w.Header().Add("HX-Retarget", "#container")
		// w.Header().Add("HX-Reswap", "afterend")
		app.render(w, "main", data)
		return
	}

	id, err := app.todos.Insert(form.Content, form.Status)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Print(err)
		return
	}
	w.Header().Add("form ID", strconv.Itoa(id))
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
	data := app.generateTemplateData(r)
	data.Lists = &list
	data.Form = todoForm{
		Status: "ongoing",
	}
	app.render(w, "main", data)
}

// func (app *application) getList(w http.ResponseWriter, r *http.Request, id int) {
// 	list, err := app.todos.Get(id)
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}
// 	data := app.generateTemplateData()
// 	data.List = list
// 	app.infoLog.Print(data.List.Status)
// 	app.render(w, "data", data)
// }

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
