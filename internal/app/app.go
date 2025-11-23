// Package app provides configuration parsing and the main application server for the review-assigner service.
//
// It defines Options for server settings, config for parsed timeouts and paths,
// and App struct that initializes HTTP mux/server and handles ListenAndServe with handler registration.
package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	validatorV10 "github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	appHttp "github.com/AndrejDubinin/review-assigner/internal/app/http"
	"github.com/AndrejDubinin/review-assigner/internal/domain"
	"github.com/AndrejDubinin/review-assigner/internal/repository/db_team_repo"
	teamService "github.com/AndrejDubinin/review-assigner/internal/services/team/add"
)

type (
	mux interface {
		Handle(pattern string, handler http.Handler)
	}
	server interface {
		ListenAndServe() error
		Close() error
	}
	logger interface {
		Info(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}
	validator interface {
		Struct(s any) error
	}
	teamStorage interface {
		AddTeam(ctx context.Context, team domain.Team) error
	}

	App struct {
		config    config
		mux       mux
		server    server
		logger    logger
		validator validator
		storage   teamStorage
	}
)

func NewApp(config config, logger logger) (*App, error) {
	mux := http.NewServeMux()

	ctx := context.Background()
	dbConfig, err := pgxpool.ParseConfig(config.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	dbConfig.MaxConns = config.db.maxConns
	dbConfig.MinConns = config.db.minConns
	dbConfig.MaxConnLifetime = config.db.maxConnLife
	dbConfig.MaxConnIdleTime = config.db.connMaxIdle

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	validator := validatorV10.New(validatorV10.WithRequiredStructEnabled())

	return &App{
		config: config,
		mux:    mux,
		server: &http.Server{
			Addr:         net.JoinHostPort(config.web.host, config.web.port),
			Handler:      mux,
			ReadTimeout:  config.web.readTimeout,
			WriteTimeout: config.web.writeTimeout,
			IdleTimeout:  config.web.idleTimeout,
		},
		logger:    logger,
		validator: validator,
		storage:   db_team_repo.NewRepo(pool),
	}, nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.index, appHttp.NewIndexHandler(a.logger))
	a.mux.Handle(a.config.path.teamAdd, appHttp.NewTeamAddTeamHandler(teamService.New(a.storage), a.config.path.teamAdd, a.logger, a.validator))

	a.logger.Info("Starting server", zap.String("address", net.JoinHostPort(a.config.web.host, a.config.web.port)))

	return a.server.ListenAndServe()
}
