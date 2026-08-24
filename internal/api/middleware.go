package api

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func requestLog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		requestPath := r.URL.Path
		next.ServeHTTP(w, r)
		logger.Info("request", "method", r.Method, "path", requestPath, "duration", time.Since(started))
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if value := recover(); value != nil {
				_ = debug.Stack()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
