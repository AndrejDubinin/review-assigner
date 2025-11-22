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
