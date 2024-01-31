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
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.getLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.postLogin))
	router.Handler(http.MethodGet, "/user/register", dynamic.ThenFunc(app.getRegister))
	router.Handler(http.MethodPost, "/user/register", dynamic.ThenFunc(app.postRegister))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.logout))
	router.Handler(http.MethodPost, "/todo/create", dynamic.ThenFunc(app.todoCreate))
	router.Handler(http.MethodGet, "/todo/created", dynamic.ThenFunc(app.getList))
	router.Handler(http.MethodPut, "/todo/delete/:id", dynamic.ThenFunc(app.deleteList))
	router.Handler(http.MethodPut, "/todo/status/:id", dynamic.ThenFunc(app.updateStatus))
	standard := alice.New(app.recoverPanic, app.requestLog, secureHeader)
	return standard.Then(router)
}
