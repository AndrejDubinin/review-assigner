package main

import (
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/AndrejDubinin/review-assigner/internal/app"
	"github.com/AndrejDubinin/review-assigner/internal/infra/logger"
)

const serviceName = "REVIEW-ASSIGNER"

func main() {
	logger, err := logger.New(serviceName)
	if err != nil {
		log.Fatal("{FATAL}", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Println("logger.Sync:", err)
		}
	}()

	if err := run(logger); err != nil {
		logger.Errorw("startup", "ERROR", err)

		if err := logger.Sync(); err != nil {
			log.Println("logger.Sync:", err)
		}

		//nolint:all
		os.Exit(1)
	}
}

func run(logger *zap.SugaredLogger) error {
	initOpts()
	config, err := app.NewConfig(opts)
	if err != nil {
		return err
	}

	service, err := app.NewApp(config, logger.Desugar())
	if err != nil {
		return err
	}

	if err := service.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
