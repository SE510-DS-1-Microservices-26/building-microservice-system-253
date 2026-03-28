ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Docker Compose
up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose up -d --build

# Database

seed:
	docker compose exec -T $(DB_HOST) psql -U $(DB_USERNAME) -d $(DB_DATABASE) < database/seeds/seed.sql

db-reset:
	docker compose down -v
	docker compose up -d

# Tests
test: unit-test api-test

unit-test:
	go test ./internal/core/services/...

api-test:
	go test ./internal/http/handlers/...

# Swagger
swagger-gen:
	swag init --parseDependency --parseInternal --dir cmd,internal --generalInfo api/main.go

.PHONY: up down build seed db-reset test unit-test api-test swagger-gen