package app

import (
	"net"
	"net/http"

	"go.uber.org/zap"

	appHttp "github.com/AndrejDubinin/review-assigner/internal/app/http"
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
		Warn(msg string, fields ...zap.Field)
	}

	App struct {
		config config
		mux    mux
		server server
		logger logger
	}
)

func NewApp(config config, logger logger) (*App, error) {
	mux := http.NewServeMux()

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
		logger: logger,
	}, nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.index, appHttp.NewIndexHandler(a.logger))

	a.logger.Info("Starting server", zap.String("address", net.JoinHostPort(a.config.web.host, a.config.web.port)))

	return a.server.ListenAndServe()
}
