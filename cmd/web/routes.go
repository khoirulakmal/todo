package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	fs := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	dynamic := alice.New(app.session.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(app.signIn))
	router.Handler(http.MethodGet, "/register", dynamic.ThenFunc(app.signUp))
	router.Handler(http.MethodPost, "/create", dynamic.ThenFunc(app.todoCreate))
	router.Handler(http.MethodGet, "/created", dynamic.ThenFunc(app.getList))
	router.Handler(http.MethodPut, "/delete/:id", dynamic.ThenFunc(app.deleteList))
	router.Handler(http.MethodPut, "/status/:id", dynamic.ThenFunc(app.updateStatus))
	standard := alice.New(app.recoverPanic, app.requestLog, secureHeader)
	return standard.Then(router)
}
