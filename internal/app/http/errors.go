package http

import (
	"errors"
	"net/http"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
	"go.uber.org/zap"
)

var (
	ErrInvalidJSONSyntax = errors.New("invalid JSON syntax")
	ErrInvalidJSON       = errors.New("invalid JSON")
)

func handleError(w http.ResponseWriter, err error, msg string, logger logger) {
	var statusCode int
	var errCode domain.ErrorCode

	switch {
	case errors.Is(err, ErrInvalidJSONSyntax) || errors.Is(err, ErrInvalidJSON):
		statusCode = http.StatusBadRequest
		errCode = domain.ErrCodeInvalidRequest

	case errors.Is(err, domain.ErrTeamExists):
		statusCode = http.StatusBadRequest
		errCode = domain.ErrCodeTeamExists

	case errors.Is(err, domain.ErrUsersInTeam):
		statusCode = http.StatusBadRequest
		errCode = domain.ErrCodeUserExists

	default:
		logger.Error("internal error", zap.Error(err))
		statusCode = http.StatusInternalServerError
		errCode = domain.ErrCodeInternalError
		msg = "internal server error"
	}

	respErr := GetErrorResponse(w, statusCode, errCode, msg)
	if respErr != nil {
		logger.Error("Failed to send error response", zap.Error(err))
	}
}
