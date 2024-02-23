package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"snow.sahej.io/db"
	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
	"snow.sahej.io/services"
	"snow.sahej.io/utils"
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

	defer db.GetInstance().Close()
	log.Printf("Starting server on %v", *PORT)

	tryFetchSaveContests := func(tick time.Time) {
		loggers.GetInstance().InfoLog.Printf("Tick: %s", tick)
		backupId, err := db.Backup()

		if _, ok := err.(db.DbDoesNotExist); err != nil && !ok {
			loggers.LogError(err)
			return
		}

		didDbExist := err != nil

		err = fetchAndSaveContests()

		if err != nil {
			loggers.LogError(err)

			if !didDbExist {
				return
			}

			if err := db.Restore(backupId); err != nil {
				loggers.LogError(err)
				return
			}

			return
		}

		loggers.GetInstance().InfoLog.Printf("Contest Fetch and Save Success")
	}

	tryFetchSaveContests(time.Now())
	ticker := utils.ExecuteEvery(time.Duration(30)*time.Second, tryFetchSaveContests)
	defer ticker.Stop()

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
		return err
	}
	err = models.EnsureContestsTableEmpty(d)
	if err != nil {
		return err
	}

	var contestsToSave []models.ContestDto

	saveContests := func(cntsts []models.ContestDto) error {
		err := models.SaveContests(cntsts, d)

		return err
	}

	for contest := range contestGenerator {
		// fmt.Printf("%s - %s: %s\n", *contest.Judge.String(), contest.EndTime, contest.Url)
		contestsToSave = append(contestsToSave, contest)

		if len(contestsToSave) >= SAVE_BATCH_SIZE {
			err := saveContests(contestsToSave)
			if err != nil {
				return err
			}

			contestsToSave = nil
		}
	}

	if contestsToSave != nil {
		err = saveContests(contestsToSave)
		if err != nil {
			return err
		}
	}

	err = models.PopulateContestsApiTable(d)

	return err
}
