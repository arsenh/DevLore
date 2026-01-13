package server

import (
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
	}
}

func (r *RateLimiter) GetLimiterIp(ip string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exist := r.visitors[ip]

	if !exist {
		limiter = rate.NewLimiter(5.0/60.0, 5)
		r.visitors[ip] = limiter
	}
	return limiter
}
