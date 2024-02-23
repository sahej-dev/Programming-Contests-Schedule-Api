// https://github.com/AliOsm/kontests

package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// name of the model is a tribute to the inspiration and the
// spiritual predecessor of this project: github.com/AliOsm/kontests
type Kontest struct {
	Name      string         `json:"name"`
	Url       url.URL        `json:"url"`
	StartTime *time.Time     `json:"start_time"`
	EndTime   *time.Time     `json:"end_time"`
	Duration  *time.Duration `json:"duration"`
	Site      Judge          `json:"site"`

	// "Yes" or "No"
	In24Hours string `json:"in_24_hours"`

	// CODING if the contest is running, BEFORE otherwise
	Status string `json:"status"`
}

func (k *Kontest) MarshalJSON() ([]byte, error) {
	type Alias Kontest

	startTimeStr := "-"
	endTimeStr := "-"
	durationStr := "-"

	if k.StartTime != nil {
		startTimeStr = k.StartTime.Format(time.RFC3339)
	}

	if k.EndTime != nil {
		endTimeStr = k.EndTime.Format(time.RFC3339)
	}

	if k.Duration != nil {
		durationStr = fmt.Sprint(k.Duration.Seconds())
	}

	return json.Marshal(&struct {
		Url       string `json:"url"`
		Judge     string `json:"site"`
		Duration  string `json:"duration"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		*Alias
	}{
		Url:       k.Url.String(),
		Judge:     *k.Site.String(),
		Duration:  durationStr,
		StartTime: startTimeStr,
		EndTime:   endTimeStr,
		Alias:     (*Alias)(k),
	})
}

func AsKontest(dto *ContestDto) *Kontest {
	if dto.Url == nil {
		return nil
	}

	var duration *time.Duration
	if dto.StartTime != nil && dto.EndTime != nil {
		d := dto.EndTime.Sub(*dto.StartTime)
		duration = &d
	}

	var in24Hours string
	if duration != nil && duration.Hours() < 24 {
		in24Hours = "Yes"
	} else {
		in24Hours = "No"
	}

	var status string
	if dto.StartTime != nil && dto.StartTime.Before(time.Now()) && dto.EndTime != nil && dto.EndTime.After(time.Now()) {
		status = "CODING"
	} else {
		status = "BEFORE"
	}

	k := Kontest{
		Name:      dto.Name,
		Url:       *dto.Url,
		StartTime: dto.StartTime,
		EndTime:   dto.EndTime,
		Duration:  duration,
		Site:      dto.Judge,
		In24Hours: in24Hours,
		Status:    status,
	}

	return &k

}
