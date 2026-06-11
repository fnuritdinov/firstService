package rate_limiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	request map[int]RequestUserInfo
}

type RequestUserInfo struct {
	Counter     int
	RequestedAt time.Time
}

func New() *RateLimiter {
	return &RateLimiter{
		request: make(map[int]RequestUserInfo),
	}
}

func (r *RateLimiter) Allow(ctx context.Context, id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	info, ok := r.request[id]
	if !ok {
		r.request[id] = RequestUserInfo{
			Counter:     1,
			RequestedAt: time.Now(),
		}
		return true
	}
	if info.Counter >= 5 {
		return false
	}
	info.Counter++
	r.request[id] = info
	return true
}

func (r *RateLimiter) WorkerClear(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx is done")
			return
		case <-ticker.C:
			r.mu.Lock()
			for id, val := range r.request {
				if time.Since(val.RequestedAt) > time.Minute {
					delete(r.request, id)
				}
			}
			r.mu.Unlock()
		}
	}
}
