package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"todo.khoirulakmal.dev/internal/validator"
)

type todoForm struct {
	Content string `form:"list"`
	Status  string `form:"status"`
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
	app.render(w, "main", data)

}

// Add a snippetCreate handler function.
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
		return
	}

	var decoded todoForm
	err := app.decodeForm(r, &decoded)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	statusList := []string{"ongoing", "pending"}

	decoded.CheckField(validator.MinChars(decoded.Content, 5), "Content", "Content must be more than 5 characters")
	decoded.CheckField(validator.MaxChars(decoded.Content, 20), "Content", "Content must be less than 20 characters")
	decoded.CheckField(validator.PermittedString(decoded.Status, statusList...), "Status", "Status must be select between pending and ongoing")

	if !decoded.Valid() {
		data := app.generateTemplateData(r)
		data.Form = decoded
		data.Page = "status"
		w.Header().Add("HX-Retarget", "#formStatus")
		w.Header().Add("HX-Reswap", "innerHTML")
		app.render(w, "main", data)
		return
	}

	id, err := app.todos.Insert(decoded.Content, decoded.Status)
	if err != nil {
		app.serverError(w, err)
		app.errorLog.Print(err)
		return
	}
	app.session.Put(r.Context(), "dataID", id)
	data := app.generateTemplateData(r)
	data.Flash = "List created success!"
	data.Page = "status"
	data.Form = todoForm{}
	w.Header().Add("HX-Trigger", "newList")
	app.render(w, "main", data)

}

func (app *application) getList(w http.ResponseWriter, r *http.Request) {
	// We can then use the ByName() method to get the value of the "id" named
	// parameter from the slice and validate it as normal.
	id := app.session.GetInt(r.Context(), "dataID")
	app.infoLog.Print(id)
	list, err := app.todos.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	updateList := app.session.PopBool(r.Context(), "updateList")
	if updateList {
		idData := strings.Join([]string{"#data-", strconv.Itoa(int(id))}, "")
		w.Header().Add("HX-Retarget", idData)
		w.Header().Add("HX-Reswap", "outerHTML")
		data := app.generateTemplateData(r)
		data.List = list
		data.Page = "list"
		app.render(w, "main", data)
		return
	}
	data := app.generateTemplateData(r)
	data.List = list
	data.Page = "list"
	w.Header().Add("HX-Retarget", "#data")
	app.infoLog.Print(data.List)
	app.render(w, "main", data)
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
	success, err := app.todos.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if success {
		idData := strings.Join([]string{"#data-", strconv.Itoa(int(id))}, "")
		w.Header().Add("HX-Retarget", idData)
		w.Header().Add("HX-Reswap", "delete")
		data := app.generateTemplateData(r)
		data.Flash = "List succesfully deleted!"
		data.Page = "status"
		data.Form = todoForm{}
		app.render(w, "main", data)
		return
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
	success, err := app.todos.Done(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if success {
		app.infoLog.Printf("Row ID is %v", id)
		app.session.Put(r.Context(), "dataID", int(id))
		app.session.Put(r.Context(), "updateList", true)
		data := app.generateTemplateData(r)
		data.Flash = "List updated success!"
		data.Page = "status"
		data.Form = todoForm{}
		w.Header().Add("HX-Trigger", "updateList")
		app.render(w, "main", data)
	}
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	data := app.generateTemplateData(r)
	app.render(w, "login", data)
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	data := app.generateTemplateData(r)
	app.render(w, "signup", data)
}
