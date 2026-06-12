package middleware

import (
	"fmt"
	"net/http"

	"github.com/fnuritdinov/firstService/internal/rate_limiter"
	"github.com/fnuritdinov/firstService/pkg/utils"
)

func RateLimit(rate *rate_limiter.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := utils.StrToInt(r.Header.Get("UserID"))
		if err != nil {
			fmt.Errorf("error from utils.StrToInt %w", err)
		}

		ok := rate.Allow(r.Context(), userID)
		if !ok {
			http.Error(w, "too many request", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
