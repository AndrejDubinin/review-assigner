package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	addTeamService interface {
		AddTeam(ctx context.Context, team domain.Team) (domain.Team, error)
	}
	validator interface {
		Struct(s any) error
	}

	addTeamRequest struct {
		Team domain.Team `json:"team" validate:"required"`
	}
	addTeamResponse struct {
		Team domain.Team `json:"team"`
	}

	AddTeamHandler struct {
		name           string
		addTeamService addTeamService
		logger         logger
		validator      validator
	}
)

func NewTeamAddTeamHandler(service addTeamService, name string, logger logger, validator validator) *AddTeamHandler {
	return &AddTeamHandler{
		name:           name,
		addTeamService: service,
		logger:         logger,
		validator:      validator,
	}
}

func (h *AddTeamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request *addTeamRequest
		err     error
	)

	if request, err = h.getRequestData(r); err != nil {
		handleError(w, ErrInvalidJSONSyntax, "invalid json syntax", h.logger)
		return
	}

	if err = h.validator.Struct(request); err != nil {
		handleError(w, ErrInvalidJSONSyntax, ConvertValidationErrors(err).String(), h.logger)
		return
	}

	team, err := h.addTeamService.AddTeam(ctx, request.Team)
	if err != nil {
		var msg string
		if errors.Is(err, domain.ErrTeamExists) {
			msg = fmt.Sprintf("%s already exists", request.Team.TeamName)
		} else if errors.Is(err, domain.ErrUsersInTeam) {
			msg = "one or more users are already in a team"
		}
		handleError(w, err, msg, h.logger)
		return
	}

	response := &addTeamResponse{
		Team: team,
	}

	marshaledTeam, err := json.Marshal(response)
	if err != nil {
		handleError(w, err, "failed to marshal team", h.logger)
		return
	}

	GetSuccessResponseWithBody(w, marshaledTeam)
}

func (h *AddTeamHandler) getRequestData(r *http.Request) (request *addTeamRequest, err error) {
	request = &addTeamRequest{}
	if err = json.NewDecoder(r.Body).Decode(request); err != nil {
		return
	}

	return
}
