ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Docker Compose
up:
	docker compose up -d --remove-orphans

down:
	docker compose down --remove-orphans

build:
	docker compose up -d --build --remove-orphans

# Database
core-seed:
	docker compose exec -T core_database psql -U $(CORE_DB_USERNAME) -d $(CORE_DB_DATABASE) < database/base/seeds/seed.sql

users-seed:
	docker compose exec -T users_database psql -U $(USERS_DB_USERNAME) -d $(USERS_DB_DATABASE) < database/users/seeds/seed.sql

db-reset:
	docker compose down -v
	docker compose up -d

# Tests
test: unit-test api-test

unit-test:
	go test ./internal/base/core/services/... ./internal/users/core/services/...

unit-test-core:
	go test ./internal/base/core/services/...

unit-test-users:
	go test ./internal/users/core/services/...

api-test:
	go test ./internal/base/http/handlers/... ./internal/users/http/handlers/...

api-test-core:
	go test ./internal/base/http/handlers/...

api-test-users:
	go test ./internal/users/http/handlers/...

# Swagger
swagger-gen: swagger-gen-core swagger-gen-users

swagger-gen-core:
	swag init --parseDependency --parseInternal --dir cmd/base,internal/base --generalInfo main.go --output docs/core

swagger-gen-users:
	swag init --parseDependency --parseInternal --dir cmd/users,internal/users --generalInfo main.go --output docs/users

.PHONY: up down build core-seed users-seed db-reset test unit-test unit-test-core unit-test-users api-test api-test-core api-test-users swagger-gen swagger-gen-core swagger-gen-users