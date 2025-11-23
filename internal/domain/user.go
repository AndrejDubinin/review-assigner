package domain

type User struct {
	Username string
	TeamID   int64
	IsActive bool
}

type UserDTO struct {
	UserID   string
	Username string
	IsActive bool
}
