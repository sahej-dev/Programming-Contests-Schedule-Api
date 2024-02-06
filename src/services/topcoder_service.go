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

type TopcoderService struct{}

func (s *TopcoderService) Url() string {
	return "https://api.topcoder.com/v5/challenges/?status=Active&isLightweight=true&perPage=100&tracks%5B%5D=Dev&tracks%5B%5D=Des&tracks%5B%5D=DS&tracks%5B%5D=QA&types%5B%5D=CH&types%5B%5D=F2F&types%5B%5D=TSK"
}

func (s *TopcoderService) FetchUpcomingContests() <-chan models.ContestDto {
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
			Id        string     `json:"id"`
			Name      string     `json:"name"`
			StartDate *time.Time `json:"startDate"`
			EndDate   *time.Time `json:"endDate"`
		}

		type jsonRespose struct {
			Status   string        `json:"status"`
			Contests []jsonContest `json:"future_contests"`
		}

		var data *[]jsonContest
		err = json.Unmarshal([]byte(string(body)), &data)
		if err != nil {
			loggers.LogError(err)
			return
		}

		for _, contest := range *data {
			contestUrl := fmt.Sprintf("https://www.topcoder.com/challenges/%s", contest.Id)
			parsedUrl, err := url.ParseRequestURI(contestUrl)
			if err != nil {
				parsedUrl = nil
			}

			// skip ended contests
			if contest.EndDate != nil && contest.EndDate.Before(time.Now()) {
				continue
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: contest.StartDate,
				EndTime:   contest.EndDate,
				Judge:     models.TopCoder,
			}
		}
	}()

	return c
}
