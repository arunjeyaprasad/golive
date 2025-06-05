package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// RecoverMiddleware wraps an http.Handler and recovers from panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"error", err,
					"stacktrace", string(debug.Stack()),
					"path", r.URL.Path,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error": "internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
