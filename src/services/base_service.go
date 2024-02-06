package services

import "snow.sahej.io/models"

type BaseService interface {
	Url() string
	FetchUpcomingContests() <-chan models.ContestDto
}
