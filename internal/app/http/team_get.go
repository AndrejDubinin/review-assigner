package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

const (
	minTeamNameLength = 3
	maxTeamNameLength = 255
)

var (
	ErrTeamNameRequired = errors.New("team_name query required")
	ErrTeamNameTooShort = fmt.Errorf("team name is too short min length is %d", minTeamNameLength)
	ErrTeamNameTooLong  = fmt.Errorf("team name is too long max length is %d", maxTeamNameLength)
)

type (
	getTeamService interface {
		GetTeam(ctx context.Context, teamName string) (domain.Team, error)
	}

	GetTeamHandler struct {
		name           string
		getTeamService getTeamService
		logger         logger
		validator      validator
	}
)

func NewGetTeamHandler(service getTeamService, name string, logger logger, validator validator) *GetTeamHandler {
	return &GetTeamHandler{
		name:           name,
		getTeamService: service,
		logger:         logger,
		validator:      validator,
	}
}

func (h *GetTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.logger = h.logger.With(
		zap.String("service", "team.get"),
		zap.String("requestID", domain.GetRequestID(ctx)),
	)

	teamName := r.URL.Query().Get("team_name")
	if err := h.validateTeamName(teamName); err != nil {
		handleError(w, ErrInvalidQuery, err.Error(), h.logger)
		return
	}

	team, err := h.getTeamService.GetTeam(ctx, teamName)
	if err != nil {
		msg := err.Error()
		if errors.Is(err, domain.ErrTeamNotFound) {
			msg = "resource not found"
		}
		handleError(w, err, msg, h.logger)
		return
	}

	teamJSON, err := json.Marshal(team)
	if err != nil {
		handleError(w, err, "failed to marshal team", h.logger)
		return
	}

	if err = GetSuccessResponseWithBody(w, teamJSON); err != nil {
		h.logger.Error("GetSuccessResponseWithBody", zap.Error(err))
	}
}

func (h *GetTeamHandler) validateTeamName(teamName string) error {
	if teamName == "" {
		return ErrTeamNameRequired
	}
	if len(teamName) < minTeamNameLength {
		return ErrTeamNameTooShort
	}
	if len(teamName) > maxTeamNameLength {
		return ErrTeamNameTooLong
	}
	return nil
}
