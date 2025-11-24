package add

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	repository interface {
		AddTeam(ctx context.Context, team domain.TeamDTO) error
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
		With(fields ...zap.Field) *zap.Logger
	}

	Handler struct {
		repo   repository
		logger logger
	}
)

func New(repo repository, logger logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

func (h *Handler) AddTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	h.logger = h.logger.With(
		zap.String("service", "team.add"),
		zap.String("requestID", domain.GetRequestID(ctx)),
	)

	if len(team.Members) == 0 {
		h.logger.Error("empty_team", zap.String("team_name", team.TeamName))
		return domain.Team{}, domain.ErrEmptyTeam
	}

	teamDTO := domain.TeamDTO{
		TeamName: team.TeamName,
		Members:  membersToUsers(team.Members),
	}

	err := h.repo.AddTeam(ctx, teamDTO)
	if err != nil {
		h.logger.Error("repo.AddTeam", zap.Error(err), zap.String("team_name", teamDTO.TeamName))
		return domain.Team{}, fmt.Errorf("repo.AddItem: %w", err)
	}

	return team, nil
}

func membersToUsers(members []domain.TeamMember) []domain.UserDTO {
	users := make([]domain.UserDTO, len(members))
	for i, member := range members {
		users[i] = domain.UserDTO(member)
	}
	return users
}
