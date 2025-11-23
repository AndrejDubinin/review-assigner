package add

import (
	"context"
	"errors"
	"fmt"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
	"github.com/AndrejDubinin/review-assigner/internal/repository/db_team_repo"
)

type (
	repository interface {
		AddTeam(ctx context.Context, team domain.Team) error
	}

	Handler struct {
		repo repository
	}
)

var ErrTeamExists = errors.New("team already exists")

func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) AddTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	err := h.repo.AddTeam(ctx, team)
	if err != nil {
		if errors.Is(err, db_team_repo.ErrTeamExists) {
			return domain.Team{}, ErrTeamExists
		}
		return domain.Team{}, fmt.Errorf("repo.AddItem failed: %w", err)
	}

	return team, nil
}
