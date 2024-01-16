package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	fs := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.todoCreate)
	mux.HandleFunc("/get", app.getList)
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
