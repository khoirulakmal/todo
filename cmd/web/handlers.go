package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"todo.khoirulakmal.dev/internal/models"
	"todo.khoirulakmal.dev/internal/validator"
)

type todoForm struct {
	Content string `form:"list"`
	Status  string `form:"status"`
	validator.Validator
}

type userRegister struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	validator.Validator
}

type userLogin struct {
	Email    string `form:"email"`
	Password string `form:"password"`
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
	data.Flash = app.session.PopString(r.Context(), "flash")
	data.Form = todoForm{
		Status: "ongo",
		Validator: validator.Validator{
			FieldErrors: nil,
		},
	}
	app.render(w, http.StatusOK, "main", data)

}

// Add a snippetCreate handler function.
func (app *application) todoCreate(w http.ResponseWriter, r *http.Request) {
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
		app.render(w, http.StatusUnprocessableEntity, "main", data)
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
	app.render(w, http.StatusAccepted, "main", data)

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
		app.render(w, http.StatusAccepted, "main", data)
		return
	}
	data := app.generateTemplateData(r)
	data.List = list
	data.Page = "list"
	w.Header().Add("HX-Retarget", "#data")
	app.infoLog.Print(data.List)
	app.render(w, http.StatusAccepted, "main", data)
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
		app.render(w, http.StatusAccepted, "main", data)
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
		app.render(w, http.StatusAccepted, "main", data)
	}
}

func (app *application) getLogin(w http.ResponseWriter, r *http.Request) {
	data := app.generateTemplateData(r)
	data.Flash = app.session.PopString(r.Context(), "flash")
	data.Form = userRegister{}
	app.render(w, http.StatusOK, "login", data)
}

func (app *application) getRegister(w http.ResponseWriter, r *http.Request) {
	data := app.generateTemplateData(r)
	data.Form = userRegister{}
	app.render(w, http.StatusOK, "signup", data)
}

func (app *application) postRegister(w http.ResponseWriter, r *http.Request) {
	var decodeForm userRegister
	err := app.decodeForm(r, &decodeForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	decodeForm.CheckField(validator.NotBlank(decodeForm.Name), "name", "Name must be filled!")
	decodeForm.CheckField(validator.NotBlank(decodeForm.Email), "email", "Email must be filled!")
	decodeForm.CheckField(validator.NotBlank(decodeForm.Password), "password", "Password must be filled!")
	decodeForm.CheckField(validator.Matches(decodeForm.Email, validator.EmailRX), "email", "Email must be in a correct format!")
	decodeForm.CheckField(validator.MinChars(decodeForm.Password, 8), "password", "Password must be more than 8 characters!")
	if !decodeForm.Valid() {
		data := app.generateTemplateData(r)
		data.Form = decodeForm
		app.render(w, http.StatusUnprocessableEntity, "signup", data)
		return
	}
	err = app.users.Insert(decodeForm.Name, decodeForm.Email, decodeForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			decodeForm.AddFieldError("email", "Email address is already in use")
			data := app.generateTemplateData(r)
			data.Form = decodeForm
			app.render(w, http.StatusUnprocessableEntity, "signup", data)
		} else {
			app.serverError(w, err)
		}
		app.errorLog.Print(err)
		return
	}
	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked.
	app.session.Put(r.Context(), "flash", "Your signup was succesful. Please log in")

	// And redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) postLogin(w http.ResponseWriter, r *http.Request) {
	var form userLogin
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "Email must be filled!")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password must be filled!")
	if !form.Valid() {
		data := app.generateTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login", data)
	}
	id, err := app.users.Auth(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFielderror("Email or password is incorrect")
			data := app.generateTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.session.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r.Context(), "authUser", id)
	app.session.Put(r.Context(), "flash", "Login Success!")

	// Redirect the user to the create snippet page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {

}
