MAIN_PACKAGE_PATH := .
BINARY_NAME := green-ecolution-backend
APP_VERSION := 0.0.1
APP_GIT_COMMIT := $(shell git rev-parse HEAD)
APP_GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
APP_GIT_REPOSITORY := https://github.com/green-ecolution/green-ecolution-backend
APP_BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
define GOFLAGS
-ldflags=" \
	-s -w \
  -X main.version=$(APP_VERSION) \
  -X github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info.version=$(APP_VERSION) \
  -X github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info.gitCommit=$(APP_GIT_COMMIT) \
  -X github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info.gitBranch=$(APP_GIT_BRANCH) \
  -X github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info.gitRepository=$(APP_GIT_REPOSITORY) \
  -X github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info.buildTime=$(APP_BUILD_TIME) \
"
endef
MOCKERY_VERSION := v2.43.2
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= postgres
POSTGRES_DB ?= postgres
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5432

.PHONY: all
all: build

.PHONY: build/all
build/all: generate
	@echo "Building for all..."
	GOARCH=amd64 GOOS=darwin go build $(GOFLAGS) -o bin/$(BINARY_NAME)-darwin $(MAIN_PACKAGE_PATH)
	GOARCH=amd64 GOOS=linux go build $(GOFLAGS) -o bin/$(BINARY_NAME)-linux $(MAIN_PACKAGE_PATH)
	GOARCH=amd64 GOOS=windows go build $(GOFLAGS) -o bin/$(BINARY_NAME)-windows $(MAIN_PACKAGE_PATH)

.PHONY: build
build: generate
	@echo "Building..."
	go build $(GOFLAGS) -o bin/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

.PHONY: generate
generate:
	@echo "Generating..."
	sqlc generate
	go generate 

.PHONY: setup
setup:
	@echo "Installing..."
	go install github.com/air-verse/air@latest
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/jmattheis/goverter/cmd/goverter@latest
	go mod download

.PHONY: setup/ci
setup/ci:
	@echo "Installing..."
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/jmattheis/goverter/cmd/goverter@latest
	go mod download

.PHONY: clean
clean:
	@echo "Cleaning..."
	go clean
	rm -rf bin
	rm -rf docs
	rm -rf tmp
	rm -rf internal/server/http/entities/info/generated
	rm -rf internal/server/http/entities/sensor/generated
	rm -rf internal/server/http/entities/tree/generated
	rm -rf internal/server/mqtt/entities/sensor/generated
	rm -rf internal/service/_mock
	rm -rf internal/storage/_mock
	rm -rf internal/storage/postgres/_sqlc
	rm -rf internal/storage/postgres/mapper/generated

.PHONY: run
run: build
	@echo "Running..."
	./bin/$(BINARY_NAME)

.PHONY: run/live
run/live: generate
	@echo "Running live..."
	air

.PHONY: migrate/new
migrate/new:
	@echo "Migrating up..."
	@if [ -z "$(name)" ]; then \
		echo "error: name is required"; \
		echo "usage: make migrate/new name=name_of_migration"; \
		exit 1; \
	fi
	goose -dir internal/storage/postgres/migrations create $(name) sql

.PHONY: migrate/up
migrate/up:
	@echo "Migrating up..."
	goose -dir internal/storage/postgres/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" up

.PHONY: migrate/down
migrate/down:
	@echo "Migrating down..."
	goose -dir internal/storage/postgres/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" down

.PHONY: migrate/reset
migrate/reset:
	@echo "Migrating reset..."
	goose -dir internal/storage/postgres/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" reset

.PHONY: migrate/status
migrate/status:
	@echo "Migrating status..."
	goose -dir internal/storage/postgres/migrations postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" status

.PHONY: seed/up
seed/up: migrate/up
	@echo "Seeding up..."
	goose -dir internal/storage/postgres/seed -no-versioning postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" up

.PHONY: seed/reset
seed/reset: migrate/up
	@echo "Seeding reset..."
	goose -dir internal/storage/postgres/seed -no-versioning postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)" reset

.PHONY: tidy
tidy:
	@echo "Fmt and Tidying..."
	go fmt ./...
	go mod tidy

.PHONY: lint
lint:
	@echo "Go fmt..."
	go fmt ./...
	@echo "Linting..."
	golangci-lint run

.PHONY: test
test:
	@echo "Testing..."
	go test -v -cover ./...

.PHONY: config/enc
config/enc:
	@echo "Encrypting config..."
	sops -e config.yaml > config.enc.yaml

.PHONY: config/dec
config/dec:
	@echo "Decrypting config..."
	sops -d config.enc.yaml > config.yaml

.PHONY: config/edit
config/edit:
	@echo "Editing config..."
	sops edit config.enc.yaml
