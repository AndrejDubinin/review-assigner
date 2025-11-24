package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func LoggingMiddleware(logger logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := NewResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			logger.Info("HTTP Request",
				zap.String("request_id", domain.GetRequestID(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration_ms", time.Since(start)),
				zap.String("content_type", r.Header.Get("Content-Type")),
			)
		})
	}
}
