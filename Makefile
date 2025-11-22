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
	@go run ./cmd/review-assigner \
		-host=$(SERVER_HOST) \
		-port=$(SERVER_PORT)


# ======================================================
# Migration
# ======================================================
.PHONY: install-goose
install-goose:
	@echo "Installing goose if missing..."
	@[ -f $(GOOSE) ] || GOBIN=$(BIN_DIR) \
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

.PHONY: goose-up
goose-up:
	@$(BIN_DIR)/goose -dir $(MIGRATE_DIR) postgres $(DB_CONN) up

.PHONY: goose-down
goose-down:
	@$(BIN_DIR)/goose -dir $(MIGRATE_DIR) postgres $(DB_CONN) down

.PHONY: goose-create
goose-create:
	@$(BIN_DIR)/goose -dir $(MIGRATE_DIR) -s create $(n) sql


# ======================================================
# Formatter installation and formatting
# ======================================================
.PHONY: install-formatters
install-formatters:
	@echo "Installing formatters if missing..."
	@[ -f $(GOFUMPT) ] || GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	@[ -f $(GCI) ] || GOBIN=$(BIN_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION)

.PHONY: format
format:
	@echo "Formatting code with gofumpt..."
	@for module in $(DIR_LIST); do \
		if [ -d $$module ]; then \
			find $$module -type f -name '*.go' ! -path '*/mocks/*' -exec $(GOFUMPT) -extra -w {} +; \
		fi \
	done
	@echo "Sorting imports with gci..."
	@for module in $(DIR_LIST); do \
		if [ -d $$module ]; then \
			find $$module -type f -name '*.go' ! -path '*/mocks/*' -exec $(GCI) write -s standard \
			-s default -s "prefix(github.com/AndrejDubinin/review-assigner)" {} +; \
		fi \
	done


# ======================================================
# Linting
# ======================================================
.PHONY: dev-gotooling
dev-gotooling: install-golangci-lint install-formatters

.PHONY: install-golangci-lint
install-golangci-lint:
	@echo "Installing golangci-lint if missing..."
	@[ -f $(GOLANGCI_LINT) ] || { mkdir -p $(BIN_DIR); GOBIN=$(BIN_DIR) go install \
	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); }

.PHONY: lint
lint:
	@echo "Running golangci-lint on all modules..."
	@set -e; ERR=0; \
	for mod in $(DIR_LIST); do \
		if [ -d $$mod ]; then \
			$(GOLANGCI_LINT) run $$mod/... --config=.golangci.yml || ERR=1; \
		fi \
	done; exit $$ERR


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
	@echo "  install-goose                 Install goose"
	@echo "  goose-up                      Migration up"
	@echo "  goose-down                    Migration down"
	@echo "  goose-create n=<filename>     Create migration file"
	@echo ""
	@echo "  install-formatters            Install formatters"
	@echo "  install-golangci-lint         Install linter"
	@echo "  format                        Format code with gofumpt and sort imports with gci"
	@echo "  lint                          Lint code with golangci-lint"
	@echo ""
