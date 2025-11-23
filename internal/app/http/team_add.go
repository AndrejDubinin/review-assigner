package http

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

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
		// ctx     = r.Context()
		request *addTeamRequest
		err     error
	)

	if request, err = h.getRequestData(r); err != nil {
		respErr := GetErrorResponse(w, http.StatusBadRequest, domain.ErrInvalidRequest, "invalid json")
		if respErr != nil {
			h.logger.Error("Failed to send error response", zap.Error(err))
		}
		return
	}

	if err = h.validator.Struct(request); err != nil {
		GetErrorResponse(w, http.StatusBadRequest, domain.ErrInvalidRequest,
			ConvertValidationErrors(err).String())
		return
	}
	/*
		err = h.addTeamService.AddTeam(

			ctx,
			request.User,
			domain.Item{
				SKU:   uint32(request.SKU),
				Count: request.Count,
			},

		)

			if err != nil {
				if errors.Is(err, add.ErrInvalidSKU) {
					GetErrorResponse(w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusNotFound)
					return
				}
				GetErrorResponse(w, h.name, fmt.Errorf("command handler failed: %w", err), http.StatusInternalServerError)
				return
			}

		GetSuccessResponseWithBody(w)
	*/
}

func (h *AddTeamHandler) getRequestData(r *http.Request) (request *addTeamRequest, err error) {
	request = &addTeamRequest{}
	if err = json.NewDecoder(r.Body).Decode(request); err != nil {
		return
	}

	return
}
