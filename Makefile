# ======================================================
# Global variables
# ======================================================
ENV_DEV_DIR                          := $(CURDIR)/build/dev
ENV_LOCAL_DIR                        := $(CURDIR)/build/local
MIGRATE_DIR                          := $(CURDIR)/migrations
DIR_LIST                             := cmd internal

GOOSE_VERSION                        := v3.25.0
GOLANGCI_LINT_VERSION                := v2.1.5
GCI_VERSION                          := v0.13.6
GOFUMPT_VERSION                      := v0.8.0

BIN_DIR                              := $(CURDIR)/bin
GOOSE                                := $(BIN_DIR)/goose
GOLANGCI_LINT                        := $(BIN_DIR)/golangci-lint
GCI                                  := $(BIN_DIR)/gci
GOFUMPT                              := $(BIN_DIR)/gofumpt

SERVICE_NAME                         := review-assigner

-include .env

ENV ?= dev
DB_CONN                              := postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_HOST_PORT}/${POSTGRES_DB}

.DEFAULT_GOAL := help


# ======================================================
# Changing files dependent on environment
# ======================================================
ifeq ($(ENV),local)
  COMPOSE_FILE := $(ENV_LOCAL_DIR)/docker-compose.yml
  ENV_FILE := $(ENV_LOCAL_DIR)/.env.example
  ENV_NAME := review-assigner-local
else ifeq ($(ENV),dev)
  COMPOSE_FILE := $(ENV_DEV_DIR)/docker-compose.yml
  ENV_FILE := $(ENV_DEV_DIR)/.env.example
  ENV_NAME := review-assigner-dev
endif


# ======================================================
# Changing environment variables
# ======================================================
.PHONY: env-local
env-local:
	@echo "Switching to LOCAL environment..."
	@cp $(ENV_LOCAL_DIR)/.env.example .env
	@cp $(ENV_LOCAL_DIR)/docker-compose.yml docker-compose.yml
	@echo "✓ Switched to LOCAL"

.PHONY: env-dev
env-dev:
	@echo "Switching to DEVELOPMENT environment..."
	@cp $(ENV_DEV_DIR)/.env.example .env
	@cp $(ENV_DEV_DIR)/docker-compose.yml docker-compose.yml
	@echo "✓ Switched to DEVELOPMENT"


# ======================================================
# Status command
# ======================================================
.PHONY: status
status:
	@echo "Current environment: $(ENV)"
	@echo "Compose files: $(COMPOSE_FILE)"
	@echo "Env file: $(ENV_FILE)"
	@echo "======================================================"
	@docker compose -p $(ENV_NAME) ps
	@echo "======================================================"
	@echo "Switch env:"
	@echo "   make env-local"
	@echo "   make env-dev"


# ======================================================
# Docker composer
# ======================================================
.PHONY: compose-up
compose-up:
	@echo "Using environment: $(ENV_NAME)"
	@docker compose -p $(ENV_NAME) up -d

.PHONY: compose-down
compose-down:
	@docker compose -p $(ENV_NAME) stop

.PHONY: compose-rm
compose-rm:
	@docker compose -p $(ENV_NAME) rm -fvs

.PHONY: compose-rs
compose-rs: ## remove previously and start new local env
	make compose-rm
	make compose-up

.PHONY: compose-logs
compose-logs:
	@docker compose -p $(ENV_NAME) logs -f $(SERVICE_NAME)


# ======================================================
# Building containers
# ======================================================
.PHONY: all
all: review-assigner

.PHONY: review-assigner
review-assigner:
	@docker compose -p $(ENV_NAME) build $(SERVICE_NAME)


# ======================================================
# Run local
# ======================================================
.PHONY: run
run:
	@go run ./cmd/review-assigner


# ======================================================
# Help command
# ======================================================
.PHONY: help
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  env-local                     Switch environment to local"
	@echo "  env-dev                       Switch environment to development"
	@echo "  status                        Current environment information"
	@echo ""
	@echo "  compose-up                    Start environment"
	@echo "  compose-down                  Terminate environment"
	@echo "  compose-rm                    Remove environment"
	@echo "  compose-rs                    Remove previously and start new environment"
	@echo "  compose-logs                  View logs of review-assigner service"
	@echo ""
	@echo "  run                           Run local review-assigner service"
	@echo ""
	@echo "  all                           Build all app containers"
	@echo "  review-assigner               Build review-assigner container"
	@echo ""
