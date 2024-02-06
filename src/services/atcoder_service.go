package services

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"snow.sahej.io/loggers"
	"snow.sahej.io/models"
)

type AtcoderService struct{}

func (s *AtcoderService) Url() string {
	return "https://atcoder.jp/contests"
}

func (s *AtcoderService) FetchUpcomingContests() <-chan models.ContestDto {
	c := make(chan models.ContestDto)

	go func() {
		scraper := colly.NewCollector()

		scraper.OnHTML("#contest-table-upcoming tbody tr", func(e *colly.HTMLElement) {

			name := e.ChildText("td:nth-of-type(2) a")
			rawUrl := fmt.Sprintf("https://atcode.jp%s", e.ChildAttr("td:nth-of-type(2) a", "href"))
			rawStartTime := e.ChildText("td:nth-of-type(1) time")
			rawDuration := e.ChildText("td:nth-of-type(3)")

			parsedUrl := s.parseUrl(rawUrl)
			startTime := s.parseStartTime(rawStartTime)
			endTime := s.getEndTime(startTime, rawDuration)

			c <- models.ContestDto{
				Name:      name,
				Url:       parsedUrl,
				StartTime: startTime,
				EndTime:   endTime,
				Judge:     models.AtCoder,
			}

		})

		scraper.OnError(func(_ *colly.Response, err error) {
			loggers.LogError(err)
		})

		scraper.OnScraped(func(r *colly.Response) {
			close(c)
		})

		err := scraper.Visit(s.Url())
		if err != nil {
			loggers.LogError(err)
		}

	}()

	return c
}

func (s *AtcoderService) parseUrl(rawUrl string) *url.URL {
	parsedUrl, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return nil
	}

	return parsedUrl
}

func (s *AtcoderService) parseStartTime(startTime string) *time.Time {
	t, err := time.Parse("2006-01-02 15:04:05-0700", startTime)

	if err != nil {
		loggers.LogError(err)
		return nil
	}

	return &t
}

func (s *AtcoderService) getEndTime(startTime *time.Time, rawDuration string) *time.Time {
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
