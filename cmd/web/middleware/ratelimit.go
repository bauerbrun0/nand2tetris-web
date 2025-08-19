package middleware

import (
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func (m *Middleware) UserLoginPostRateLimiter(next http.Handler) http.Handler {
	if m.Config.Env == "test" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	store := memory.NewStore()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rate, err := limiter.NewRateFromFormatted("10-M")
		if err != nil {
			panic(err)
		}

		instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))
		mw := stdlib.NewMiddleware(instance, stdlib.WithKeyGetter(func(r *http.Request) string {
			if err := r.ParseForm(); err != nil {
				return ""
			}
			username := r.FormValue("username")
			if len(username) > 100 {
				username = username[:100]
			}
			ip := r.Header.Get("X-Real-IP")

			return username + ip
		}))
		mw.Handler(next).ServeHTTP(w, r)
	})
}

func (m *Middleware) GetIPRateLimiter(limit int64, period time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if m.Config.Env == "test" {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
		store := memory.NewStore()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rate := limiter.Rate{
				Period: period,
				Limit:  limit,
			}
			instance := limiter.New(store, rate, limiter.WithClientIPHeader("X-Real-IP"), limiter.WithTrustForwardHeader(true))
			mw := stdlib.NewMiddleware(instance)
			mw.Handler(next).ServeHTTP(w, r)
		})
	}
}
