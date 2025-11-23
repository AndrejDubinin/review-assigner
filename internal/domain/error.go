package domain

import "errors"

type ErrorCode string

const (
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeTeamExists     ErrorCode = "TEAM_EXISTS"
	ErrCodeUserExists     ErrorCode = "USER_EXISTS"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
)

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrUsersInTeam = errors.New("one or more users are already in a team")
)
