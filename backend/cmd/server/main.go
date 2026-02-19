package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/codebytamla/borderform-backend/internal/auth"
	"github.com/codebytamla/borderform-backend/internal/database"
	mw "github.com/codebytamla/borderform-backend/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present (silently ignored in production).
	_ = godotenv.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// -----------------------------------------------------------------
	// Config (env vars with safe defaults)
	// -----------------------------------------------------------------
	dbURL := envOr("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/borderform?sslmode=disable")
	listenAddr := envOr("LISTEN_ADDR", ":8080")
	jwtSecret := envOr("JWT_SECRET", "") // MUST be set in production
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// -----------------------------------------------------------------
	// Database
	// -----------------------------------------------------------------
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}
	log.Println("connected to database")

	store := database.NewPGStore(pool)

	// -----------------------------------------------------------------
	// JWT config
	// -----------------------------------------------------------------
	jwtCfg := auth.JWTConfig{
		SigningKey: []byte(jwtSecret),
		Issuer:     "borderform",
		Audience:   "borderform-api",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
		ClockSkew:  30 * time.Second,
	}

	// -----------------------------------------------------------------
	// Auth service + handlers
	// -----------------------------------------------------------------
	authSvc := auth.NewService(store, jwtCfg, store)
	authHandler := auth.NewHandler(authSvc)

	// -----------------------------------------------------------------
	// Router (see routes.go for all route definitions)
	// -----------------------------------------------------------------
	mux := http.NewServeMux()
	registerRoutes(mux, authHandler, jwtCfg, store)

	// -----------------------------------------------------------------
	// Server
	// -----------------------------------------------------------------

	// Wrap entire mux with CORS and request logging middleware.
	allowedOrigins := envOr("CORS_ORIGINS", "http://localhost:5173")
	handler := mw.RequestLogger(mw.CORS(strings.Split(allowedOrigins, ",")...)(mux))

	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", listenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("stopped")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// chain applies middleware in reverse order so they execute left-to-right.
// Usage: chain(handler, mw1, mw2) => mw1(mw2(handler))
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
