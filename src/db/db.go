package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"snow.sahej.io/loggers"
)

type DatabaseInteractor struct {
	db *sql.DB
}

var instance *DatabaseInteractor
var lock = &sync.Mutex{}

func GetInstance() *DatabaseInteractor {
	if instance != nil {
		return instance
	}

	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		db, err := sql.Open("sqlite3", "sql_database.sqlite")
		if err != nil {
			loggers.LogError(err)
			return nil
		}

		instance = &DatabaseInteractor{
			db: db,
		}

	}

	return instance
}

func (d *DatabaseInteractor) Exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)

	return err
}