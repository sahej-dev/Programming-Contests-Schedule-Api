package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"snow.sahej.io/loggers"
	"snow.sahej.io/utils"
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
		utils.EnsureDirExists(dB_PATH)
		db, err := sql.Open("sqlite3", GetDbPath())
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

func (d *DatabaseInteractor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *DatabaseInteractor) Close() error {
	return d.db.Close()
}

func (d *DatabaseInteractor) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}
