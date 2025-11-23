package domain

type TeamMember struct {
	UserID   string `json:"user_id" validate:"required,gte=2,lte=255"`
	Username string `json:"username" validate:"required,gte=3,lte=255"`
	IsActive bool   `json:"is_active" validate:"required,boolean"`
}

type Team struct {
	TeamName string       `json:"team_name" validate:"required,gte=3,lte=255"`
	Members  []TeamMember `json:"members" validate:"required,dive"`
}

type TeamDTO struct {
	TeamName string
	Members  []UserDTO
}
