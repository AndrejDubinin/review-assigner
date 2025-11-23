package add

import (
	"context"
	"fmt"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	repository interface {
		AddTeam(ctx context.Context, team domain.TeamDTO) error
	}

	Handler struct {
		repo repository
	}
)

func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) AddTeam(ctx context.Context, team domain.Team) (domain.Team, error) {
	teamDTO := domain.TeamDTO{
		TeamName: team.TeamName,
		Members:  membersToUsers(team.Members),
	}

	err := h.repo.AddTeam(ctx, teamDTO)
	if err != nil {
		return domain.Team{}, fmt.Errorf("repo.AddItem: %w", err)
	}

	return team, nil
}

func membersToUsers(members []domain.TeamMember) []domain.UserDTO {
	users := make([]domain.UserDTO, len(members))
	for i, member := range members {
		users[i] = domain.UserDTO{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}
	return users
}
