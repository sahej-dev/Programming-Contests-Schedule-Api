package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
)

type HackerEarthService struct{}

func (s *HackerEarthService) Url() string {
	return "https://www.hackerearth.com/chrome-extension/events/"
}

func (s *HackerEarthService) FetchUpcomingContests() <-chan models.ContestDto {
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
			Name      string  `json:"title"`
			Url       *string `json:"url"`
			StartTime *string `json:"start_utc_tz"`
			EndTime   *string `json:"end_utc_tz"`
		}

		type jsonResponse struct {
			Contests []jsonContest `json:"response"`
		}

		var data *jsonResponse
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range data.Contests {

			parsedUrl, err := url.ParseRequestURI(*contest.Url)
			if err != nil {
				parsedUrl = nil
			}

			var startTime *time.Time = s.parseDateTime(contest.StartTime)
			var endTime *time.Time = s.parseDateTime(contest.EndTime)

			// skip ended contests
			if endTime != nil && endTime.Before(time.Now()) {
				continue
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: startTime,
				EndTime:   endTime,
				Judge:     models.HackerEarth,
			}
		}
	}()

	return c
}

func (s *HackerEarthService) parseDateTime(datetime *string) *time.Time {
	if datetime == nil {
		return nil
	}

	t, err := time.Parse(time.RFC3339, strings.Join(strings.Split(*datetime, " "), "T"))
	if err != nil {
		return nil
	}

	return &t
}
