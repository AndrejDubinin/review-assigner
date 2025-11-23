package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	APIError struct {
		Code    domain.ErrorCode `json:"code"`
		Message string           `json:"message"`
	}
	ErrorResponse struct {
		Error APIError `json:"error"`
	}
)

func GetSuccessResponseWithBody(w http.ResponseWriter, body []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	if err != nil {
		return fmt.Errorf("w.Write: %w", err)
	}
	return nil
}

func GetErrorResponse(w http.ResponseWriter, statusCode int, code domain.ErrorCode, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errEnc := json.NewEncoder(w).Encode(ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
	if errEnc != nil {
		return fmt.Errorf("json.NewEncoder(w).Encode: %w", errEnc)
	}
	return nil
}
