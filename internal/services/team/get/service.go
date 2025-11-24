package get

import (
	"context"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	repository interface {
		GetTeam(ctx context.Context, teamName string) (domain.Team, error)
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

func (h *Handler) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	h.logger = h.logger.With(
		zap.String("service", "team.get"),
		zap.String("requestID", domain.GetRequestID(ctx)),
	)

	team, err := h.repo.GetTeam(ctx, teamName)
	if err != nil {
		h.logger.Error("repo.GetTeam", zap.Error(err), zap.String("team_name", teamName))
		return domain.Team{}, err
	}

	return team, nil
}
