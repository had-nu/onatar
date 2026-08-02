package httpapi

import (
	"net/http"
)

// AuthEndpoints registers the authentication endpoints.
func (s *Server) AuthEndpoints() http.Handler {
	mux := http.NewServeMux()

	// Public endpoints
	mux.HandleFunc("GET /auth/github", s.handleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", s.handleGitHubCallback)

	// Protected endpoints
	mux.Handle("GET /auth/me", s.RequireAuth(http.HandlerFunc(s.handleGetMe)))
	mux.Handle("POST /auth/logout", s.RequireAuth(http.HandlerFunc(s.handleLogout)))

	return mux
}

// handleGitHubLogin initiates the GitHub OAuth flow.
func (s *Server) handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	state, authURL := s.auth.BeginOAuthFlow()

	// Store state in a short-lived cookie for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   s.auth.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, authURL, http.StatusFound)
}

// handleGitHubCallback handles the GitHub OAuth callback.
func (s *Server) handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state parameter
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		s.logger.Warn("oauth state cookie missing", "error", err)
		writeError(w, http.StatusBadRequest, "INVALID_STATE", "invalid or missing state parameter", nil)
		return
	}

	queryState := r.URL.Query().Get("state")
	if stateCookie.Value != queryState {
		s.logger.Warn("oauth state mismatch", "cookie", stateCookie.Value, "query", queryState)
		writeError(w, http.StatusBadRequest, "INVALID_STATE", "state parameter mismatch", nil)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "MISSING_CODE", "authorization code is required", nil)
		return
	}

	// Exchange code for token and fetch user
	user, err := s.auth.ExchangeCode(r.Context(), code)
	if err != nil {
		s.logger.Error("oauth exchange failed", "error", err)
		writeError(w, http.StatusInternalServerError, "OAUTH_EXCHANGE_FAILED", "failed to exchange authorization code", nil)
		return
	}

	// Upsert user in database
	if err := s.auth.UpsertUser(r.Context(), user); err != nil {
		s.logger.Error("failed to upsert user", "error", err)
		writeError(w, http.StatusInternalServerError, "USER_UPSERT_FAILED", "failed to save user", nil)
		return
	}

	// Create session
	session, err := s.auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		s.logger.Error("failed to create session", "error", err)
		writeError(w, http.StatusInternalServerError, "SESSION_CREATE_FAILED", "failed to create session", nil)
		return
	}

	// Set session cookie
	s.setSessionCookie(w, session.ID)

	// Clear OAuth state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.auth.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to home page (or original destination)
	http.Redirect(w, r, "/", http.StatusFound)
}

// handleGetMe returns the current authenticated user.
func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// handleLogout logs out the current user.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(s.auth.Cfg.SessionCookieName)
	if err == nil {
		_ = s.auth.DeleteSession(r.Context(), cookie.Value)
	}

	// Clear session cookie
	s.clearSessionCookie(w)

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}