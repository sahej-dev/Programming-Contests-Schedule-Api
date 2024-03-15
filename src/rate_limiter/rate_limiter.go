package ratelimiter

import (
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ipMap map[string]*rate.Limiter
	mutex sync.Mutex
}

var instance *RateLimiter
var lock = &sync.Mutex{}

func GetInstance() *RateLimiter {
	if instance != nil {
		return instance
	}

	lock.Lock()
	defer lock.Unlock()

	if instance == nil {

		instance = &RateLimiter{
			ipMap: make(map[string]*rate.Limiter),
		}

	}

	return instance
}

func (r *RateLimiter) GetLimiterForIp(ip string) *rate.Limiter {
	if limiter, ok := r.ipMap[ip]; ok {
		return limiter
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 20 requests allowed per minute
	newLimiter := rate.NewLimiter(1/3.0, 20)
	r.ipMap[ip] = newLimiter

	return newLimiter
}
