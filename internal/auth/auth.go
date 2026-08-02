package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/hadnu/onatar/internal/config"
)

// User represents an authenticated user.
type User struct {
	ID        string `json:"id"`
	GitHubID  int64  `json:"github_id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

// Session represents an authenticated session.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthService handles authentication operations.
type AuthService struct {
	db        *sql.DB
	logger    *slog.Logger
	Cfg       *config.AuthConfig
	oauth2Cfg *oauth2.Config
	jwtSecret []byte
}

// NewAuthService creates a new authentication service.
func NewAuthService(db *sql.DB, logger *slog.Logger, cfg *config.AuthConfig) *AuthService {
	return &AuthService{
		db:        db,
		logger:    logger.With("component", "auth"),
		Cfg:       cfg,
		jwtSecret: []byte(cfg.JWTSecret),
		oauth2Cfg: &oauth2.Config{
			ClientID:     cfg.GitHubClientID,
			ClientSecret: cfg.GitHubClientSecret,
			RedirectURL:  cfg.GitHubRedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetOAuth2Config returns the OAuth2 configuration for GitHub.
func (s *AuthService) GetOAuth2Config() *oauth2.Config {
	return s.oauth2Cfg
}

// BeginOAuthFlow generates a random state and returns the authorization URL.
func (s *AuthService) BeginOAuthFlow() (state, authURL string) {
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		s.logger.Error("failed to generate oauth state", "error", err)
		state = uuid.New().String()
	} else {
		state = hex.EncodeToString(stateBytes)
	}
	authURL = s.oauth2Cfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
	return state, authURL
}

// decodeJSON decodes JSON from an io.Reader.
func decodeJSON(reader interface{}, v interface{}) error {
	if r, ok := reader.(interface{ Read([]byte) (int, error) }); ok {
		return json.NewDecoder(r).Decode(v)
	}
	return fmt.Errorf("body does not implement io.Reader")
}

// ExchangeCode exchanges the authorization code for a token and fetches the GitHub user.
func (s *AuthService) ExchangeCode(ctx context.Context, code string) (*User, error) {
	token, err := s.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth2 exchange: %w", err)
	}

	client := s.oauth2Cfg.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("fetch github user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user api returned %d", resp.StatusCode)
	}

	var ghUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}
	if err := decodeJSON(resp.Body, &ghUser); err != nil {
		return nil, fmt.Errorf("decode github user: %w", err)
	}

	// If email is not public, fetch from emails endpoint
	if ghUser.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailsResp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if err := decodeJSON(emailsResp.Body, &emails); err == nil {
				for _, e := range emails {
					if e.Primary && e.Verified {
						ghUser.Email = e.Email
						break
					}
				}
			}
		}
	}

	user := &User{
		ID:        fmt.Sprintf("gh:%d", ghUser.ID),
		GitHubID:  ghUser.ID,
		Login:     ghUser.Login,
		Name:      ghUser.Name,
		AvatarURL: ghUser.AvatarURL,
		Email:     ghUser.Email,
	}

	return user, nil
}

// UpsertUser inserts or updates a user in the database.
func (s *AuthService) UpsertUser(ctx context.Context, user *User) error {
	return UpsertUser(ctx, s.db, user)
}

// CreateSession creates a new session for the user.
func (s *AuthService) CreateSession(ctx context.Context, userID string) (*Session, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(s.Cfg.SessionTTL)

	if err := CreateSession(ctx, s.db, sessionID, userID, expiresAt); err != nil {
		return nil, err
	}

	return &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

// GetSession retrieves a session by ID.
func (s *AuthService) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	return GetSession(ctx, s.db, sessionID)
}

// DeleteSession deletes a session by ID.
func (s *AuthService) DeleteSession(ctx context.Context, sessionID string) error {
	return DeleteSession(ctx, s.db, sessionID)
}

// GetUserByID retrieves a user by ID.
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return GetUserByID(ctx, s.db, userID)
}

// GenerateJWT generates a JWT token for the user.
func (s *AuthService) GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.Cfg.SessionTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT validates a JWT token and returns the user ID.
func (s *AuthService) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", fmt.Errorf("jwt parse: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}
	return "", errors.New("invalid token claims")
}

// CleanupExpiredSessions removes expired sessions from the database.
func (s *AuthService) CleanupExpiredSessions(ctx context.Context) error {
	return CleanupExpiredSessions(ctx, s.db)
}

// UpsertUser inserts or updates a user in the database.
func UpsertUser(ctx context.Context, db *sql.DB, user *User) error {
	query := `
		INSERT INTO users (id, github_id, login, name, avatar_url, email)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			login = VALUES(login),
			name = VALUES(name),
			avatar_url = VALUES(avatar_url),
			email = VALUES(email),
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := db.ExecContext(ctx, query, user.ID, user.GitHubID, user.Login, user.Name, user.AvatarURL, user.Email)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

// CreateSession creates a new session for the user.
func CreateSession(ctx context.Context, db *sql.DB, sessionID, userID string, expiresAt time.Time) error {
	query := `
		INSERT INTO sessions (id, user_id, expires_at)
		VALUES (?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query, sessionID, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

// GetSession retrieves a session by ID.
func GetSession(ctx context.Context, db *sql.DB, sessionID string) (*Session, error) {
	query := `SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?`
	row := db.QueryRowContext(ctx, query, sessionID)

	var session Session
	err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		_ = DeleteSession(ctx, db, sessionID)
		return nil, nil
	}

	return &session, nil
}

// DeleteSession deletes a session by ID.
func DeleteSession(ctx context.Context, db *sql.DB, sessionID string) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, sessionID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID.
func GetUserByID(ctx context.Context, db *sql.DB, userID string) (*User, error) {
	query := `SELECT id, github_id, login, name, avatar_url, email FROM users WHERE id = ?`
	row := db.QueryRowContext(ctx, query, userID)

	var user User
	err := row.Scan(&user.ID, &user.GitHubID, &user.Login, &user.Name, &user.AvatarURL, &user.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// CleanupExpiredSessions removes expired sessions from the database.
func CleanupExpiredSessions(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("cleanup sessions: %w", err)
	}
	return nil
}