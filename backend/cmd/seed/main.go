package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/codebytamla/borderform-backend/internal/auth"
	"github.com/codebytamla/borderform-backend/internal/database/db"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	name := os.Getenv("ADMIN_NAME")

	if email == "" || password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD are required")
	}
	if name == "" {
		name = "System Admin"
	}

	if err := auth.ValidatePasswordPolicy(password); err != nil {
		log.Fatalf("invalid password: %v", err)
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("begin transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := db.New(tx)

	// Create the admin user.
	user, err := q.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		FullName:     name,
	})
	if err != nil {
		log.Fatalf("create admin user: %v", err)
	}

	// Assign the ADMIN role.
	role, err := q.GetRoleByName(ctx, "ADMIN")
	if err != nil {
		log.Fatalf("lookup admin role: %v", err)
	}
	if err := q.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: user.ID,
		RoleID: role.ID,
	}); err != nil {
		log.Fatalf("assign admin role: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit transaction: %v", err)
	}

	fmt.Printf("admin user created: %s (%s)\n", email, user.ID.Bytes)
}
