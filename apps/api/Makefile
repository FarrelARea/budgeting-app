.PHONY: dev build test migrate seed sqlc

dev:
	go run .

build:
	go build -o bin/api .

test:
	go test ./...

migrate:
	psql "$(DATABASE_URL)" -f migrations/000001_init.up.sql
	psql "$(DATABASE_URL)" -f migrations/000002_users.up.sql

seed:
	go run ./cmd/seed

sqlc:
	sqlc generate
