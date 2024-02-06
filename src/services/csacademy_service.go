package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"io"

	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
)

type CsacademyService struct{}

func (s *CsacademyService) Url() string {
	return "https://csacademy.com/contests/?"
}

func (s *CsacademyService) FetchUpcomingContests() <-chan models.ContestDto {
	c := make(chan models.ContestDto)

	go func() {
		defer close(c)

		client := &http.Client{}

		req, err := http.NewRequest("GET", s.Url(), nil)
		if err != nil {
			loggers.LogError(err)
			return
		}

		// Mimicing XHR request gets us a JSON response
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")

		resp, err := client.Do(req)
		if err != nil {
			loggers.LogError(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(io.Reader(resp.Body))
		if err != nil {
			loggers.LogError(err)
			return
		}

		type jsonContest struct {
			Name      string   `json:"longName"`
			Slug      *string  `json:"name"`
			StartTime *float64 `json:"startTime"` // some are sent as floats
			EndTime   *float64 `json:"endTime"`
		}

		type contestsResponse struct {
			Contests []jsonContest `json:"Contest"`
		}

		type jsonResponse struct {
			Data contestsResponse `json:"state"`
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
				contestUrl := fmt.Sprintf("https://csacademy.com/contest/%s", *contest.Slug)
				parsedUrl, err = url.ParseRequestURI(contestUrl)
				if err != nil {
					parsedUrl = nil
				}
			}

			var startTime *time.Time
			var endTime *time.Time

			if contest.StartTime != nil {
				t := time.Unix(int64(*contest.StartTime), 0).UTC()
				startTime = &t
			}

			if contest.EndTime != nil {
				t := time.Unix(int64(*contest.EndTime), 0).UTC()
				endTime = &t
			}

			// skip ended contests
			if startTime == nil || (endTime != nil && endTime.Before(time.Now())) {
				continue
			}

			c <- models.ContestDto{
				Name:      contest.Name,
				Url:       parsedUrl,
				StartTime: startTime,
				EndTime:   endTime,
				Judge:     models.CsAcademy,
			}
		}

	}()

	return c
}

func (s *CsacademyService) parseUrl(rawUrl string) *url.URL {
	parsedUrl, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return nil
	}

	return parsedUrl
}

func (s *CsacademyService) parseStartTime(startTime string) *time.Time {
	t, err := time.Parse("2006-01-02 15:04:05-0700", startTime)

	if err != nil {
		loggers.LogError(err)
		return nil
	}

	return &t
}

func (s *CsacademyService) getEndTime(startTime *time.Time, rawDuration string) *time.Time {
	if startTime == nil {
		return nil
	}

	parts := strings.Split(rawDuration, ":")

	if len(parts) != 2 {
		return nil
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		loggers.LogError(err)
		return nil
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		loggers.LogError(err)
		return nil
	}

	durationSecs := h*3600 + m*60
	t := time.Unix(startTime.Unix()+int64(durationSecs), 0).UTC()

	return &t
}
