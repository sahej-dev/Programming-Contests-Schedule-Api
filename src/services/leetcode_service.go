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

type LeetcodeService struct{}

func (s *LeetcodeService) Url() string {
	return "https://leetcode.com/graphql?query=%7B%20allContests%20%7B%20title%20titleSlug%20startTime%20duration%20__typename%20%7D%20%7D"
}

func (s *LeetcodeService) FetchUpcomingContests() <-chan models.ContestDto {
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
			Slug      *string `json:"titleSlug"`
			StartTime *int64  `json:"startTime"`
			Duration  *int64  `json:"duration"`
		}

		type contestsResponse struct {
			Contests []jsonContest `json:"allContests"`
		}

		type jsonResponse struct {
			Data contestsResponse `json:"data"`
		}

		var data *jsonResponse
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range data.Data.Contests {
			var parsedUrl *url.URL
			if *&contest.Slug != nil {
				contestUrl := fmt.Sprintf("https://leetcode.com/contest/%s/", *contest.Slug)
				parsedUrl, err = url.ParseRequestURI(contestUrl)
				if err != nil {
					parsedUrl = nil
				}
			}

			var startTime *time.Time
			var endTime *time.Time

			if contest.StartTime != nil {
				t := time.Unix(*contest.StartTime, 0).UTC()
				startTime = &t
			}

			if startTime != nil && contest.Duration != nil {
				t := time.Unix(*contest.StartTime+*contest.Duration, 0).UTC()
				endTime = &t
			}

			// skip ended contests
			if endTime != nil && endTime.Before(time.Now()) {
				continue
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: startTime,
				EndTime:   endTime,
				Judge:     models.Leetcode,
			}
		}
	}()

	return c
}
