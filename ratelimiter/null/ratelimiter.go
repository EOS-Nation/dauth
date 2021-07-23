package null

import (
	"github.com/eosnationftw/dauth/ratelimiter"
)

func init() {
	// null://
	ratelimiter.Register("null", func(configURL string) (ratelimiter.RateLimiter, error) {
		return NewRequestRateLimiter(), nil
	})
}

type RequestRateLimiter struct{}

func NewRequestRateLimiter() *RequestRateLimiter {
	return &RequestRateLimiter{}
}

func (r *RequestRateLimiter) Gate(id string, method string) (allow bool) {
	return true
}
