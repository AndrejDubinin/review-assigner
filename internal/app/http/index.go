// Package http provides the root HTTP handler for health checks.
//
// It defines an IndexHandler that responds to requests at the root path ("/")
// with a success message indicating the 'review-assigner' service is online.
// The handler logs errors if the response cannot be sent.
package http

import (
	"net/http"

	"go.uber.org/zap"
)

type (
	logger interface {
		Error(msg string, fields ...zap.Field)
	}

	IndexHandler struct {
		logger logger
	}
)

func NewIndexHandler(logger logger) *IndexHandler {
	return &IndexHandler{logger: logger}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	err := GetSuccessResponseWithBody(w, []byte("Service 'review-assigner' is online"))
	if err != nil {
		h.logger.Error("Failed to send success response", zap.Error(err))
	}
}
