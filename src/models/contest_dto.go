package models

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"snow.sahej.io/db"
	"snow.sahej.io/loggers"
)

type ContestDto struct {
	Name      string     `db:"name"`
	Url       *url.URL   `db:"url"`
	StartTime *time.Time `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
	Judge     Judge      `db:"site"`
}

func (c ContestDto) String() string {
	url := "<nil>"
	if c.Url != nil {
		url = c.Url.String()
	}
	judge := "<nil>"
	if c.Judge.String() != nil {
		judge = *c.Judge.String()
	}
	return fmt.Sprintf("Name: %v\nURL: %v\nStart Time: %v\nEnd Time: %v\nJudge: %v",
		c.Name, url, c.StartTime.Format(time.RFC3339), c.EndTime.Format(time.RFC3339), judge)
}

func (contest *ContestDto) Save(d *db.DatabaseInteractor) error {
	query := `
        INSERT INTO contests (name, url, start_time, end_time, judge)
        VALUES (?, ?, ?, ?, ?)
    `

	var urlString string
	if contest.Url != nil {
		urlString = contest.Url.String()
	}

	err := d.Exec(query, contest.Name, urlString, contest.StartTime, contest.EndTime, contest.Judge)
	return err
}

func SaveContests(contests []ContestDto, d *db.DatabaseInteractor) error {
	if len(contests) == 0 {
		return nil
	}

	query := `
        INSERT INTO contests (name, url, start_time, end_time, judge)
        VALUES
    `
	values := []interface{}{}

	for i, contest := range contests {
		if i > 0 {
			query += ","
		}
		query += " (?, ?, ?, ?, ?)"
		urlString := ""
		if contest.Url != nil {
			urlString = contest.Url.String()
		}
		values = append(values, contest.Name, urlString, contest.StartTime, contest.EndTime, contest.Judge)
	}

	err := d.Exec(query, values...)
	return err
}

func EnsureContestsTableExists(d *db.DatabaseInteractor) error {
	err := d.Exec(`
        CREATE TABLE IF NOT EXISTS contests (
            name TEXT,
            url TEXT,
            start_time TIMESTAMP,
            end_time TIMESTAMP,
            judge TEXT
        )
    `)

	return err
}

func EnsureContestsTableEmpty(d *db.DatabaseInteractor) error {
	err := d.Exec(`
        DELETE FROM contests
        WHERE TRUE
    `)

	return err
}

func PopulateContestsApiTable(d *db.DatabaseInteractor) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        CREATE TABLE IF NOT EXISTS contests_api (
            name TEXT,
            url TEXT,
            start_time TIMESTAMP,
            end_time TIMESTAMP,
            judge TEXT
        )
    `)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
        DELETE FROM contests_api
        WHERE TRUE
    `)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO contests_api SELECT * FROM contests")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()

	return err
}

func GetAllApiContests(d *db.DatabaseInteractor) ([]ContestDto, error) {
	contests := make([]ContestDto, 0)

	rows, err := d.Query(`
		SELECT name, url, start_time, end_time, judge
		FROM contests_api
		ORDER BY start_time
	`)
	if err != nil {
		return contests, err
	}
	defer rows.Close()

	for rows.Next() {
		var contest ContestDto
		var urlString sql.NullString
		err := rows.Scan(&contest.Name, &urlString, &contest.StartTime, &contest.EndTime, &contest.Judge)
		if err != nil {
			loggers.LogError(err)
			continue
		}
		if urlString.Valid {
			url, err := url.Parse(urlString.String)
			if err != nil {
				return contests, err
			}
			contest.Url = url
		}
		contests = append(contests, contest)
	}

	if err := rows.Err(); err != nil {
		return contests, err
	}

	return contests, nil
}
