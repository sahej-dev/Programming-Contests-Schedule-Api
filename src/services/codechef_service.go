package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"io"

	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
)

type CodechefService struct{}

func (s *CodechefService) Url() string {
	return "https://www.codechef.com/api/list/contests/all"
}

func (s *CodechefService) FetchUpcomingContests() <-chan models.ContestDto {
	c := make(chan models.ContestDto)

	go func() {
		defer close(c)

		response, err := http.Get(s.Url())
		if err != nil {
			loggers.LogError(err)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(io.Reader(response.Body))
		if err != nil {
			loggers.LogError(err)
			return
		}

		type jsonContest struct {
			Code      string    `json:"contest_code"`
			Name      string    `json:"contest_name"`
			StartDate time.Time `json:"contest_start_date_iso"`
			EndDate   time.Time `json:"contest_end_date_iso"`
		}

		type jsonRespose struct {
			Status   string        `json:"status"`
			Contests []jsonContest `json:"future_contests"`
		}

		var data *jsonRespose
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range data.Contests {
			contestUrl := fmt.Sprintf("https://www.codechef.com/%s", contest.Code)
			parsedUrl, err := url.ParseRequestURI(contestUrl)
			if err != nil {
				parsedUrl = nil
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: &contest.StartDate,
				EndTime:   &contest.EndDate,
				Judge:     models.Codechef,
			}
		}
	}()

	return c
}
