package httpapi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hadnu/onatar/internal/auth"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const (
	userContextKey     contextKey = "user"
	sessionContextKey  contextKey = "session"
)

// AuthMiddleware creates middleware that validates the session cookie.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(s.auth.Cfg.SessionCookieName)
		if err != nil {
			s.logger.Debug("no session cookie", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		session, err := s.auth.GetSession(r.Context(), cookie.Value)
		if err != nil {
			s.logger.Error("failed to get session", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		if session == nil {
			// Session expired or invalid
			s.clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		user, err := s.auth.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			s.logger.Error("failed to get user", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		if user == nil {
			s.clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		// Add user and session to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		ctx = context.WithValue(ctx, sessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is middleware that requires authentication.
// If no valid session, returns 401 Unauthorized.
func (s *Server) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireDM is middleware that requires the user to be the DM of the campaign.
// It expects the campaign ID to be in the URL path as {id}.
func (s *Server) RequireDM(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
			return
		}

		campaignID := chi.URLParam(r, "id")
		if campaignID == "" {
			writeError(w, http.StatusBadRequest, "INVALID_ID", "campaign ID is required", nil)
			return
		}

		isDM, err := s.store.IsUserDM(r.Context(), campaignID, user.ID)
		if err != nil {
			s.logger.Error("check DM", "error", err)
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check permissions", nil)
			return
		}

		if !isDM {
			writeError(w, http.StatusForbidden, "FORBIDDEN", "only DM can perform this action", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// OptionalAuth adds user to context if authenticated, but doesn't require it.
func (s *Server) OptionalAuth(next http.Handler) http.Handler {
	return s.AuthMiddleware(next)
}

// GetUserFromContext retrieves the user from the request context.
func GetUserFromContext(ctx context.Context) *auth.User {
	if user, ok := ctx.Value(userContextKey).(*auth.User); ok {
		return user
	}
	return nil
}

// GetSessionFromContext retrieves the session from the request context.
func GetSessionFromContext(ctx context.Context) *auth.Session {
	if session, ok := ctx.Value(sessionContextKey).(*auth.Session); ok {
		return session
	}
	return nil
}

// clearSessionCookie clears the session cookie.
func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.Cfg.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.auth.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

// setSessionCookie sets the session cookie.
func (s *Server) setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.auth.Cfg.SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(s.auth.Cfg.SessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   s.auth.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}