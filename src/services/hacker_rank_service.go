package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
)

type HackerRankService struct{}

func (s *HackerRankService) Url() string {
	return "https://www.hackerrank.com/rest/contests/upcoming?limit=100"
}

func (s *HackerRankService) FetchUpcomingContests() <-chan models.ContestDto {
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
			Name      string  `json:"name"`
			Slug      *string `json:"slug"`
			StartTime *string `json:"get_starttimeiso"`
			EndTime   *string `json:"get_endtimeiso"`
		}

		type jsonResponse struct {
			Contests []jsonContest `json:"models"`
		}

		var data *jsonResponse
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range data.Contests {
			var parsedUrl *url.URL
			if *&contest.Slug != nil {
				contestUrl := fmt.Sprintf("https://www.hackerrank.com/contests/%s/challenges", *contest.Slug)
				parsedUrl, err = url.ParseRequestURI(contestUrl)
				if err != nil {
					parsedUrl = nil
				}
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
				Judge:     models.HackerRank,
			}
		}
	}()

	return c
}

func (s *HackerRankService) parseDateTime(datetime *string) *time.Time {
	if datetime == nil {
		return nil
	}

	t, err := time.Parse(time.RFC3339, *datetime)
	if err != nil {
		return nil
	}

	return &t
}
