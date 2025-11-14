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


# ======================================================
# Changing files dependent on environment
# ======================================================
ifeq ($(ENV),local)
	COMPOSE_FILE := $(ENV_LOCAL_DIR)/docker-compose.yml
  ENV_FILE := $(ENV_LOCAL_DIR)/.env.example
  ENV_NAME := review-assigner-local
else ifeq ($(ENV),dev)
	COMPOSE_FILE := $(ENV_DEV_DIR)/docker-compose.yml
  ENV_FILE := $(ENV_DEV_DIR)/.env
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
