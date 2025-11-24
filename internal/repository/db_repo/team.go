package db_repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

type (
	DBTX interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}

	Repo struct {
		conn *pgxpool.Pool
	}
)

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

func (r *Repo) AddTeam(ctx context.Context, team domain.TeamDTO) error {
	err := r.InTx(ctx, func(tx pgx.Tx) error {
		var err error
		teamID, err := r.addTeam(ctx, tx, team.TeamName)
		if err != nil {
			return fmt.Errorf("r.addTeam: %w", err)
		}

		err = r.addUsers(ctx, tx, teamID, team.Members)
		if err != nil {
			return fmt.Errorf("r.addUsers: %w", err)
		}

		return nil
	})

	return err
}

func (r *Repo) addTeam(ctx context.Context, tx pgx.Tx, teamName string) (int64, error) {
	const query = `
	INSERT INTO teams (name, created_at, updated_at)
	VALUES ($1, $2, $3) RETURNING id;`

	now := time.Now()
	var id int64

	var db DBTX = r.conn
	if tx != nil {
		db = tx
	}

	err := db.QueryRow(ctx, query, teamName, now, now).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, domain.ErrTeamExists
		}
		return 0, err
	}

	return id, nil
}

func (r *Repo) addUsers(ctx context.Context, tx pgx.Tx, teamID int64, users []domain.UserDTO) error {
	if len(users) == 0 {
		return nil
	}

	const colsNum = 6
	now := time.Now()
	var sb strings.Builder
	args := make([]any, 0, len(users)*colsNum)

	sb.WriteString("INSERT INTO users (id, username, team_id, is_active, created_at, updated_at) VALUES ")

	for i, user := range users {
		if i > 0 {
			sb.WriteString(", ")
		}
		paramOffset := i*colsNum + 1
		sb.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", paramOffset, paramOffset+1,
			paramOffset+2, paramOffset+3, paramOffset+4, paramOffset+5))

		args = append(args, user.UserID, user.Username, teamID, user.IsActive, now, now)
	}

	var db DBTX = r.conn
	if tx != nil {
		db = tx
	}

	_, err := db.Exec(ctx, sb.String(), args...)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUsersInTeam
		}
		return err
	}
	return nil
}

func (r *Repo) GetTeam(ctx context.Context, teamName string) (domain.Team, error) {
	const query = `
	SELECT t.id, u.id, u.username, u.is_active from teams t
	LEFT JOIN users u on t.id = u.team_id
	WHERE name = $1;`

	var team domain.Team

	rows, err := r.conn.Query(ctx, query, teamName)
	if err != nil {
		return domain.Team{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var member domain.TeamMember
		var teamID string

		if err := rows.Scan(&teamID, &member.UserID, &member.Username, &member.IsActive); err != nil {
			return domain.Team{}, err
		}

		team.Members = append(team.Members, member)
	}

	if err := rows.Err(); err != nil {
		return domain.Team{}, err
	}

	team.TeamName = teamName
	if len(team.Members) == 0 {
		return domain.Team{}, domain.ErrTeamNotFound
	}

	return team, nil
}
