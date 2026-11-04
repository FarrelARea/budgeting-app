.PHONY: dev build test migrate sqlc

dev:
	go run .

build:
	go build -o bin/api .

test:
	go test ./...

migrate:
	psql "$(DATABASE_URL)" -f migrations/000001_init.up.sql

sqlc:
	sqlc generate
