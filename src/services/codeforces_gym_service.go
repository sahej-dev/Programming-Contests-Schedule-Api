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

type CodeforcesGymService struct{}

func (s *CodeforcesGymService) Url() string {
	return "https://codeforces.com/api/contest.list?gym=true"
}

func (s *CodeforcesGymService) FetchUpcomingContests() <-chan models.ContestDto {
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
			Id               int64  `json:"id"`
			Name             string `json:"name"`
			StartTimeSecs    *int64 `json:"startTimeSeconds"`
			DurationSecs     int64  `json:"durationSeconds"`
			RelativeTimeSecs *int64 `json:"relativeTimeSeconds"`
		}

		type jsonResponse struct {
			Status   string        `json:"status"`
			Contests []jsonContest `json:"result"`
		}

		var data *jsonResponse
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range data.Contests {
			// skip ended contests
			if contest.RelativeTimeSecs != nil && *contest.RelativeTimeSecs >= contest.DurationSecs {
				continue
			}

			// skip contests with no startTime
			if contest.StartTimeSecs == nil {
				continue
			}

			contestUrl := fmt.Sprintf("https://codeforces.com/gymRegistration/%d", contest.Id)
			parsedUrl, err := url.ParseRequestURI(contestUrl)
			if err != nil {
				parsedUrl = nil
			}
			var startTime *time.Time
			var endTime *time.Time
			if contest.StartTimeSecs != nil {
				st := time.Unix(*contest.StartTimeSecs, 0).UTC()
				startTime = &st
				et := time.Unix(*contest.StartTimeSecs+contest.DurationSecs, 0).UTC()
				endTime = &et
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: startTime,
				EndTime:   endTime,
				Judge:     models.CodeforcesGym,
			}
		}
	}()

	return c
}
