package domain

import "errors"

type ErrorCode string

const (
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeTeamExists     ErrorCode = "TEAM_EXISTS"
	ErrCodeUserExists     ErrorCode = "USER_EXISTS"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
)

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrUsersInTeam  = errors.New("one or more users are already in a team")
	ErrEmptyTeam    = errors.New("team is empty")
	ErrTeamNotFound = errors.New("team not found")
)
