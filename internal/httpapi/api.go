// Package httpapi exposes the Onatar MVP API (PRD §3.5): /health,
// /api/v1/content and /api/v1/build.
// Beta: adds /auth/* endpoints for GitHub OAuth.
package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hadnu/onatar/internal/auth"
	"github.com/hadnu/onatar/internal/config"
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
	buildTimeout         = 5 * time.Second
	contentTimeout       = 5 * time.Second
)

// Server holds the API dependencies.
type Server struct {
	version       string
	logger        *slog.Logger
	contentLoader contentLoader
	limiter       *tokenBucket
	auth          *auth.AuthService
	cfg           *config.AuthConfig
	db            *sql.DB
	store         *store.Queries
}

// New builds the API server backed by the given DB.
func New(db *sql.DB, logger *slog.Logger, version string, cfg *config.Config) *Server {
	authService := auth.NewAuthService(db, logger, &cfg.Auth)
	return &Server{
		version:       version,
		logger:        logger,
		contentLoader: dbContentLoader{db: db},
		limiter:       newTokenBucket(defaultRatePerMinute, defaultCapacity),
		auth:          authService,
		cfg:           &cfg.Auth,
		db:            db,
		store:         store.NewQueries(db),
	}
}

// Router assembles all routes.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(s.accessLog)

	r.Get("/health", s.handleHealth)

	// Auth endpoints (public)
	r.Mount("/auth", s.AuthEndpoints())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/content", s.handleContent)
		r.With(middleware.Timeout(buildTimeout), s.rateLimit).Post("/build", s.handleBuild)

		// Protected routes (Beta)
		r.Route("/characters", func(r chi.Router) {
			r.Use(s.RequireAuth)
			r.Get("/", s.handleListCharacters)
			r.Post("/", s.handleCreateCharacter)
			r.Get("/{id}", s.handleGetCharacter)
			r.Put("/{id}", s.handleUpdateCharacter)
			r.Delete("/{id}", s.handleDeleteCharacter)
		})

		r.Route("/campaigns", func(r chi.Router) {
			r.Use(s.RequireAuth)
			r.Get("/", s.handleListCampaigns)
			r.Post("/", s.handleCreateCampaign)
			r.Get("/{id}", s.handleGetCampaign)
			r.Get("/{id}/characters", s.handleListCampaignCharacters)
			r.Post("/{id}/members", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.RequireDM(http.HandlerFunc(s.handleAddCampaignMember)).ServeHTTP(w, r)
			}))
			r.Delete("/{id}/members/{userId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				s.RequireDM(http.HandlerFunc(s.handleRemoveCampaignMember)).ServeHTTP(w, r)
			}))
		})
	})
	return r
}

// handleHealth returns the health check response.
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}

// accessLog logs each request.
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

// handleBuild calculates a character sheet from a draft.
func (s *Server) handleBuild(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), buildTimeout)
	defer cancel()

	var req BuildRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DRAFT", "invalid request body", map[string]any{"error": err.Error()})
		return
	}

	if err := validateBuildRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DRAFT", err.Error(), nil)
		return
	}

	sheet, err := s.buildSheet(ctx, req)
	if err != nil {
		s.logger.Error("build failed", "error", err)
		writeError(w, http.StatusUnprocessableEntity, "BUILD_ERROR", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, BuildResponse{Sheet: *sheet})
}

// decodeJSON decodes JSON from the request body.
func decodeJSON(body interface{}, v interface{}) error {
	if r, ok := body.(interface{ Read([]byte) (int, error) }); ok {
		return json.NewDecoder(r).Decode(v)
	}
	return fmt.Errorf("body does not implement io.Reader")
}