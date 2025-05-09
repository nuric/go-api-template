package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// ---------------------------
// Zerolog based middleware for logging HTTP requests
func ZeroLoggerMetrics(next http.Handler) http.Handler {
	handler := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})(next)
	handler = hlog.NewHandler(log.Logger)(handler)
	return handler
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Interface("error", err).Msg("panic recovered")
				log.Error().Str("stack", string(debug.Stack())).Msg("stack trace")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
