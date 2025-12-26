package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// NewMWLogger возвращает middleware, который логирует информацию о запросе
func NewMWLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger := log.With(
			slog.String("component", "middleware/logger"),
		)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srw := &statusResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			entry := logger.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
			)

			start := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", srw.status),
					slog.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
				)
			}()

			next.ServeHTTP(srw, r)
		})
	}
}
