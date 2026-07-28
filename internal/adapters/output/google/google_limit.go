package google

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(10), 5),
	}
}

func NewRateLimiterWithConfig(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

func (r *RateLimiter) Wait(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.limiter == nil {
		return nil
	}

	return r.limiter.Wait(ctx)
}

func (r *RateLimiter) Allow() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.limiter == nil {
		return true
	}

	return r.limiter.Allow()
}

func (r *RateLimiter) Reserve() *rate.Reservation {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.limiter == nil {
		return nil
	}

	return r.limiter.Reserve()
}

func (r *RateLimiter) SetLimit(limit rate.Limit) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.limiter != nil {
		r.limiter.SetLimit(limit)
	}
}

func (r *RateLimiter) SetBurst(burst int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.limiter != nil {
		r.limiter.SetBurst(burst)
	}
}

func (r *RateLimiter) GetLimit() rate.Limit {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.limiter == nil {
		return rate.Limit(0)
	}

	return r.limiter.Limit()
}

func (r *RateLimiter) GetBurst() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.limiter == nil {
		return 0
	}

	return r.limiter.Burst()
}

func (r *RateLimiter) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

}

func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.limiter != nil {
		limit := r.limiter.Limit()
		burst := r.limiter.Burst()
		r.limiter = rate.NewLimiter(limit, burst)
	}
}

func (r *RateLimiter) IsHealthy() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.limiter != nil
}
