package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		slog.Info(fmt.Sprintf("Response: %s %s, took %v", r.Method, r.URL.Path, time.Since(start)))
	})
}
