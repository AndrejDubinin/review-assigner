package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func GetSuccessResponseWithBody(w http.ResponseWriter, body []byte) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(body)
	if err != nil {
		return fmt.Errorf("w.Write: %w", err)
	}
	return nil
}

func GetErrorResponse(w http.ResponseWriter, status int, err error, details string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errEnc := json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
	if errEnc != nil {
		return fmt.Errorf("json.NewEncoder(w).Encode: %w", errEnc)
	}
	return nil
}
