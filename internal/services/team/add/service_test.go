package add

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/domain"
)

func TestHandler_AddTeam(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo   func(mc *minimock.Controller) repository
		logger logger
	}
	type args struct {
		//nolint:all
		ctx  context.Context
		team domain.Team
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Team
		wantErr error
	}{
		{
			name: "success: team with multiple members",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "backend",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
								{UserID: "u2", Username: "Bob", IsActive: true},
							},
						},
					).Return(nil)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: domain.SetRequestID(context.Background(), "req-123"),
				team: domain.Team{
					TeamName: "backend",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
						{UserID: "u2", Username: "Bob", IsActive: true},
					},
				},
			},
			want: domain.Team{
				TeamName: "backend",
				Members: []domain.TeamMember{
					{UserID: "u1", Username: "Alice", IsActive: true},
					{UserID: "u2", Username: "Bob", IsActive: true},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: team with single member",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "frontend",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
							},
						},
					).Return(nil)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "frontend",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
					},
				},
			},
			want: domain.Team{
				TeamName: "frontend",
				Members: []domain.TeamMember{
					{UserID: "u1", Username: "Alice", IsActive: true},
				},
			},
			wantErr: nil,
		},
		{
			name: "success: team with inactive members",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "mixed-team",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
								{UserID: "u2", Username: "Bob", IsActive: false},
							},
						},
					).Return(nil)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "mixed-team",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
						{UserID: "u2", Username: "Bob", IsActive: false},
					},
				},
			},
			want: domain.Team{
				TeamName: "mixed-team",
				Members: []domain.TeamMember{
					{UserID: "u1", Username: "Alice", IsActive: true},
					{UserID: "u2", Username: "Bob", IsActive: false},
				},
			},
			wantErr: nil,
		},
		{
			name: "error: team already exists",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "backend",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
							},
						},
					).Return(domain.ErrTeamExists)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "backend",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
					},
				},
			},
			want:    domain.Team{},
			wantErr: domain.ErrTeamExists,
		},
		{
			name: "error: team with empty members",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "empty-team",
					Members:  []domain.TeamMember{},
				},
			},
			want:    domain.Team{},
			wantErr: domain.ErrEmptyTeam,
		},
		{
			name: "error: users already in team",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "payments",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
							},
						},
					).Return(domain.ErrUsersInTeam)
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "payments",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
					},
				},
			},
			want:    domain.Team{},
			wantErr: domain.ErrUsersInTeam,
		},
		{
			name: "error: repository generic error",
			fields: fields{
				repo: func(mc *minimock.Controller) repository {
					repo := NewRepositoryMock(mc)
					repo.AddTeamMock.Expect(
						minimock.AnyContext,
						domain.TeamDTO{
							TeamName: "backend",
							Members: []domain.UserDTO{
								{UserID: "u1", Username: "Alice", IsActive: true},
							},
						},
					).Return(errors.New("database connection failed"))
					return repo
				},
				logger: zap.NewNop(),
			},
			args: args{
				ctx: context.Background(),
				team: domain.Team{
					TeamName: "backend",
					Members: []domain.TeamMember{
						{UserID: "u1", Username: "Alice", IsActive: true},
					},
				},
			},
			want:    domain.Team{},
			wantErr: errors.New("repo.AddItem: database connection failed"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			h := &Handler{
				repo:   tt.fields.repo(mc),
				logger: tt.fields.logger,
			}

			got, err := h.AddTeam(tt.args.ctx, tt.args.team)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
				assert.Equal(t, tt.want, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
