package services

import (
	"sync"

	"snow.sahej.io/models"
)

func FanIn(services ...BaseService) <-chan models.ContestDto {
	c := make(chan models.ContestDto)
	var wg sync.WaitGroup

	for _, s := range services {
		wg.Add(1)
		go func(s BaseService) {
			defer wg.Done()
			for contest := range s.FetchUpcomingContests() {
				c <- contest
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}
