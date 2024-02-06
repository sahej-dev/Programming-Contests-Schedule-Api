package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"snow.sahej.io/db"
	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
	"snow.sahej.io/services"
)

type Application struct {
}

func main() {
	PORT := flag.String("port", ":4000", "HTTP listen and serve address")

	flag.Parse()

	app := &Application{}

	srv := &http.Server{
		Addr:     *PORT,
		ErrorLog: loggers.GetInstance().ErrorLog,
		Handler:  app.routes(),
	}

	log.Printf("Starting server on %v", *PORT)

	fetchAndSaveContests()

	err := srv.ListenAndServe()

	loggers.GetInstance().ErrorLog.Fatal(err)
}

func fetchAndSaveContests() error {
	SAVE_BATCH_SIZE := 100

	contestGenerator := services.FanIn(
		&services.CodechefService{},
		&services.CodeforcesService{},
		&services.CodeforcesGymService{},
		&services.TopcoderService{},
		&services.HackerEarthService{},
		&services.LeetcodeService{},
		&services.AtcoderService{},
		&services.CsacademyService{},
	)

	d := db.GetInstance()
	if d == nil {
		loggers.GetInstance().ErrorLog.Fatal("Unable to get db instance")
	}

	err := models.EnsureContestsTableExists(d)
	if err != nil {
		// TODO: restore from previous db snapshot
	}
	err = models.EnsureContestsTableEmpty(d)
	if err != nil {
		// TODO: restore from previous db snapshot
		//
		// Idea: just make a copy of the sqlite file before hand
		// as a backup
	}

	var contestsToSave []models.ContestDto

	saveContests := func(cntsts []models.ContestDto) {
		err := models.SaveContests(cntsts, d)
		if err != nil {
			loggers.LogError(err)
		}
	}

	for contest := range contestGenerator {
		fmt.Printf("%s - %s: %s\n", *contest.Judge.String(), contest.EndTime, contest.Url)
		contestsToSave = append(contestsToSave, contest)

		if len(contestsToSave) >= SAVE_BATCH_SIZE {
			saveContests(contestsToSave)

			contestsToSave = nil
		}
	}

	if contestsToSave != nil {
		saveContests(contestsToSave)
	}

	return nil
}
