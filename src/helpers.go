package main

import (
	"fmt"
	"net/http"

	"snow.sahej.io/loggers"
)

func (app *Application) WriteJson(w http.ResponseWriter, content []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(content)
}

func (app *Application) serverError(w http.ResponseWriter, err error) {
	loggers.LogError(err)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *Application) notFountError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *Application) handleErrorByClosingConnection(w http.ResponseWriter) {
	if err := recover(); err != nil {
		w.Header().Set("Connection", "close")
		app.serverError(w, fmt.Errorf("%s", err))
	}
}
