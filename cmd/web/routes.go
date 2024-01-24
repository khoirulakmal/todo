package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	fs := http.FileServer(http.Dir("./ui/static/"))

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodPost, "/create", app.todoCreate)
	router.HandlerFunc(http.MethodPut, "/delete/:id", app.deleteList)
	router.HandlerFunc(http.MethodPut, "/status/:id", app.updateStatus)
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))
	standard := alice.New(app.recoverPanic, app.requestLog, secureHeader, app.session.LoadAndSave).Then(router)
	return standard
}
