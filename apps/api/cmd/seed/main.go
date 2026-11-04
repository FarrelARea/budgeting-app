package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/farrelarea/budgeting-app/apps/api/internal/auth"
)

func main() {
	_ = godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	email := os.Getenv("SEED_USER_EMAIL")
	if email == "" {
		email = "admin@budgeting.local"
	}
	password := os.Getenv("SEED_USER_PASSWORD")
	if password == "" {
		password = "password123"
	}
	name := os.Getenv("SEED_USER_NAME")
	if name == "" {
		name = "Admin"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}

	var id string
	err = pool.QueryRow(context.Background(),
		`INSERT INTO users (email, password_hash, name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, name = EXCLUDED.name
		 RETURNING id`, email, hash, name).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("seeded user: %s (id=%s) password=%s\n", email, id, password)
}
