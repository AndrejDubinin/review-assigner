package db_team_repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

var ErrTeamExists = errors.New("team already exists")

type Repo struct {
	conn *pgxpool.Pool
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	return &Repo{
		conn: conn,
	}
}

func (r *Repo) InTx(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	err = f(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) AddTeam(ctx context.Context, team domain.Team) error {
	err := r.InTx(ctx, func(tx pgx.Tx) error {
		var err error
		teamID, err := r.addTeam(ctx, team.TeamName)
		if err != nil {
			return err
		}

		err = r.addUsers(ctx, membersToUsers(teamID, team.Members))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *Repo) addTeam(ctx context.Context, teamName string) (int64, error) {
	const query = `
	INSERT INTO teams (name, created_at, updated_at)
	VALUES ($1, $2, $3) RETURNING id;`

	now := time.Now()
	var id int64
	err := r.conn.QueryRow(ctx, query, teamName, now, now).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repo) addUsers(ctx context.Context, users []domain.User) error {
	if len(users) == 0 {
		return nil
	}

	const colsNum = 5
	now := time.Now()
	var sb strings.Builder

	sb.WriteString("INSERT INTO users (username, team_id, is_active, created_at, updated_at) VALUES ")

	args := make([]any, 0, len(users)*colsNum)

	for i, user := range users {
		if i > 0 {
			sb.WriteString(", ")
		}
		paramOffset := i*colsNum + 1
		sb.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", paramOffset, paramOffset+1,
			paramOffset+2, paramOffset+3, paramOffset+4))

		args = append(args,
			user.Username,
			user.TeamID,
			user.IsActive,
			now,
			now,
		)
	}

	_, err := r.conn.Exec(ctx, sb.String(), args...)
	return err
}

func membersToUsers(teamID int64, members []domain.TeamMember) []domain.User {
	users := make([]domain.User, len(members))
	for i, member := range members {
		users[i] = domain.User{
			Username: member.Username,
			TeamID:   teamID,
			IsActive: member.IsActive,
		}
	}
	return users
}
