package main

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	return app.logRequest(mux)
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Snow 9000"))
}
