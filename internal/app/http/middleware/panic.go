package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

func PanicMiddleware(logger logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(
						"panic",
						zap.Any("error", err),
						zap.String("request_id", domain.GetRequestID(r.Context())),
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
