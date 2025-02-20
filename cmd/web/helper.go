package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, e error) {
	trace := fmt.Sprintf("%s\n%s", e.Error(), debug.Stack())
	app.errorLog.Printf(trace)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	pageHTML := fmt.Sprintf("%s.html", page)
	ts, ok := app.templateCache[pageHTML]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	renderpage := "base"
	if len(data.Page) > 0 {
		renderpage = data.Page
	}
	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, renderpage, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) isUserAuth(r *http.Request) bool {
	return app.session.Exists(r.Context(), "authUser")
}

func (app *application) generateTemplateData(r *http.Request) *templateData {
	return &templateData{
		Year:  time.Time.Year(time.Now()),
		Flash: app.session.PopString(r.Context(), "flash"),
		Auth:  app.isUserAuth(r),
	}
}

func (app *application) decodeForm(r *http.Request, form any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecode.Decode(&form, r.PostForm)
	if err != nil {
		return err
	}
	return nil
}
