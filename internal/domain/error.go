package domain

import "errors"

type ErrorCode string

const (
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
)

var (
	ErrTeamExists  = errors.New("team already exists")
	ErrUsersInTeam = errors.New("one or more users are already in a team")
)
