BINARY        := go-ecommerce
CMD           := ./cmd/api
BIN           := bin/$(BINARY)
MIGRATIONS    := ./internal/db/migrations

DATABASE_URL  ?= postgres://goecommerce:goecommerce@postgres:5432/goecommerce?sslmode=disable

# Make these visible to goose without repeating them on every target
export GOOSE_DRIVER        := postgres
export GOOSE_DBSTRING      := $(DATABASE_URL)
export GOOSE_MIGRATION_DIR := $(MIGRATIONS)

.DEFAULT_GOAL := build
.PHONY: build run dev test lint tidy audit \
        migrate-up migrate-down migrate-status migrate-reset migrate-create

build:
	go build -o $(BIN) $(CMD)

run:
	go run $(CMD)

dev:
	air -c .air.toml

test:
	go test -race ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
	go mod verify

audit: lint test
	govulncheck ./...

migrate-up:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

migrate-reset:
	goose reset

# usage: make migrate-create name=add_products
migrate-create:
	goose -s create $(name) sql

docker-up:
	docker compose -f docker/compose.yaml up -d

docker-down:
	docker compose -f docker/compose.yaml down