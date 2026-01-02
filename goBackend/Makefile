.PHONY: all proto build run run-all test clean docker-up docker-down deps help


# ======================
# Variables
# ======================
# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

PROTO_DIR = proto
SERVICES  = auth-service project-service task-service analytics-service media-service
PROTO_SVCS = auth project task analytics media

# Default Database Config (fallback if not in .env)
DB_USER ?= postgres
DB_PASSWORD ?= 123456789
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= gobackend
DB_SSLMODE ?= disable

# ======================
# Default
# ======================
all: proto build

# ======================
# OS Detection
# ======================
ifeq ($(OS),Windows_NT)

# =========================================================
# ===================== WINDOWS (cmd) =====================
# =========================================================

proto:
	@echo Generating protobuf files on Windows...
	@for %%s in ($(PROTO_SVCS)) do ( \
		protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. $(PROTO_DIR)\%%s\*.proto \
	)

build:
	@echo Building all services on Windows...
	@for %%s in ($(SERVICES)) do ( \
		echo Building %%s... && \
		cd services\%%s && \
		go build -o ..\..\bin\%%s.exe cmd\main.go && \
		cd ..\..\ \
	)

deps:
	@echo Downloading dependencies on Windows...
	go work sync
	@for %%s in ($(SERVICES)) do ( \
		cd services\%%s && go mod download && cd ..\..\ \
	)

clean:
	@echo Cleaning on Windows...
	@if exist bin rmdir /s /q bin

else

# =========================================================
# =================== LINUX / MACOS =======================
# =========================================================

SHELL := bash

proto:
	@echo "Generating protobuf files on Linux/macOS..."
	@for svc in $(PROTO_SVCS); do \
		protoc --go_out=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		 $(PROTO_DIR)/$$svc/*.proto; \
	done

build:
	@echo "Building all services on Linux/macOS..."
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		cd services/$$svc && go build -o ../../bin/$$svc ./cmd/main.go; \
		cd ../..; \
	done

deps:
	@echo "Downloading dependencies on Linux/macOS..."
	go work sync
	@for svc in $(SERVICES); do \
		cd services/$$svc && go mod download && cd ../..; \
	done

clean:
	@echo "Cleaning on Linux/macOS..."
	@rm -rf bin/

endif

# ======================
# Common Targets
# ======================

run:
	@echo "Running $(service)..."
	@echo "services/$(service)..."
	@cd services/$(service) && go run ./cmd/main.go
ifeq ($(OS),Windows_NT)
SHELL := C:/Program Files/Git/bin/bash.exe
run-all:
	@for svc in $(SERVICES); do \
		APP_BASE_DIR=$(CURDIR)/services/$$svc; \
		(cd $(CURDIR)/services/$$svc && go run ./cmd/main.go) & \
	done; \
	wait
else
run-all:
	@cd $(CURDIR) && \
	for svc in $(SERVICES); do \
		echo "Running $$svc..."; \
		cd services/$$svc && go run ./cmd/main.go & \
	done; \
	wait
endif
test:
	@echo "Running tests..."
	@go test ./... -v

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-logs:
	docker-compose logs -f

db-migrate:
	@echo "Running migrations..."
	docker-compose exec postgres psql -U postgres -d portfolio -f /docker-entrypoint-initdb.d/001_init.sql

db-create-local:
	@echo "Creating local database if not exists..."
ifeq ($(OS),Windows_NT)
	@psql "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/postgres" \
	-c "CREATE DATABASE $(DB_NAME);" || echo "Database might already exist"
else
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
	-c "CREATE DATABASE $(DB_NAME);" || true
endif

db-migrate-local:
	@echo "Running local migrations with psql..."
ifeq ($(OS),Windows_NT)
	@for %%f in (migrations\*.sql) do ( \
		echo Applying %%f && \
		psql "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" -f %%f \
	)
else
	@for f in migrations/*.sql; do (\
		echo "Applying $$f"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f "$$f"; \
	)
endif

db-reset:
	@echo "Resetting database..."
	docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS portfolio;"
	docker-compose exec postgres psql -U postgres -c "CREATE DATABASE portfolio;"
	$(MAKE) db-migrate

install-proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

help:
	@echo "Available targets:"
	@echo "  all                  - Generate proto + build"
	@echo "  proto                - Generate protobuf files"
	@echo "  build                - Build all services"
	@echo "  run SERVICE=xxx      - Run a service"
	@echo "  test                 - Run tests"
	@echo "  clean                - Clean artifacts"
	@echo "  deps                 - Download dependencies"
	@echo "  docker-up/down       - Docker compose up/down"
	@echo "  docker-build         - Build docker images"
	@echo "  docker-logs          - Tail docker logs"
	@echo "  lint                 - Run linter"

lint:
	@echo "Running linter..."
	golangci-lint run ./...
