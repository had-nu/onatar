// Package httpapi exposes the Onatar MVP API (PRD §3.5): /health,
// /api/v1/content and /api/v1/build.
package httpapi

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/store"
)

// contentLoader decouples handlers from the DB so they can be unit-tested.
type contentLoader interface {
	LoadContent(ctx context.Context) (*content.Content, error)
}

type dbContentLoader struct{ db *sql.DB }

func (l dbContentLoader) LoadContent(ctx context.Context) (*content.Content, error) {
	return store.Content(ctx, l.db)
}

const (
	defaultRatePerMinute = 10.0
	defaultCapacity      = 10.0
	// buildTimeout bounds POST /build work (PRD threat model §7).
	buildTimeout = 5 * time.Second
	// contentTimeout bounds the content DB read.
	contentTimeout = 5 * time.Second
)

// Server holds the API dependencies.
type Server struct {
	version       string
	logger        *slog.Logger
	contentLoader contentLoader
	limiter       *tokenBucket
}

// New builds the API server backed by the given DB.
func New(db *sql.DB, logger *slog.Logger, version string) *Server {
	return &Server{
		version:       version,
		logger:        logger,
		contentLoader: dbContentLoader{db: db},
		limiter:       newTokenBucket(defaultRatePerMinute, defaultCapacity),
	}
}

// Router assembles all routes.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(s.accessLog)

	r.Get("/health", s.handleHealth)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/content", s.handleContent)
		r.With(middleware.Timeout(buildTimeout), s.rateLimit).Post("/build", s.handleBuild)
	})
	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}

func (s *Server) accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		s.logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"dur_ms", time.Since(start).Milliseconds(),
		)
	})
}
