package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/farrelarea/budgeting-app/apps/api/internal/config"
	"github.com/farrelarea/budgeting-app/apps/api/internal/handler"
)

// Server wires HTTP routing with dependencies.
type Server struct {
	cfg    config.Config
	router *chi.Mux
	pool   *pgxpool.Pool
}

// New builds a server, connecting to Postgres. The caller owns pool lifecycle.
func New(ctx context.Context, cfg config.Config) (*Server, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	s := &Server{cfg: cfg, pool: pool}
	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(30 * time.Second))

	s.router.Get("/healthz", handler.Health(cfg, pool))
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/accounts", handler.ListAccounts(pool))
	})

	return s, nil
}

// Handler returns the http.Handler for the server.
func (s *Server) Handler() http.Handler { return s.router }

// Close releases pooled resources.
func (s *Server) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}
