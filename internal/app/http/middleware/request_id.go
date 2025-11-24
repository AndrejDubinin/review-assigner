package middleware

import (
	"net/http"

	"github.com/rs/xid"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = xid.New().String()
		}
		ctx := domain.SetRequestID(r.Context(), requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
