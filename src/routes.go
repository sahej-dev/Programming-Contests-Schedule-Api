package main

import (
	"encoding/json"
	"net/http"

	"snow.sahej.io/db"
	"snow.sahej.io/models"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)

	return app.rateLimitIfRequired(app.logRequest(mux))
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	contests, err := models.GetAllApiContests(db.GetInstance())
	if err != nil {
		app.serverError(w, err)
		return
	}

	var kontests []models.Kontest

	for _, contest := range contests {
		kontests = append(kontests, *models.AsKontest(&contest))
	}

	jsonData, err := json.Marshal(kontests)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.WriteJson(w, jsonData)
}
