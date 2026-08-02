package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string
	DBPass   string
	HTTPAddr string

	// Auth configuration
	Auth AuthConfig
}

type AuthConfig struct {
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
	JWTSecret          string
	SessionCookieName  string
	SessionTTL         time.Duration
	CookieSecure       bool
}

func Load() (*Config, error) {
	c := &Config{
		DBHost:   env("DB_HOST", "127.0.0.1"),
		DBPort:   env("DB_PORT", "3306"),
		DBName:   env("DB_NAME", "onatar"),
		DBUser:   env("DB_USER", "onatar"),
		DBPass:   env("DB_PASS", ""),
		HTTPAddr: env("HTTP_ADDR", ":8090"),
		Auth: AuthConfig{
			GitHubClientID:     env("GITHUB_CLIENT_ID", ""),
			GitHubClientSecret: env("GITHUB_CLIENT_SECRET", ""),
			GitHubRedirectURL:  env("GITHUB_REDIRECT_URL", "http://localhost:5173/auth/callback"),
			JWTSecret:          env("JWT_SECRET", "dev-secret-change-in-production"),
			SessionCookieName:  env("SESSION_COOKIE_NAME", "onatar_session"),
			SessionTTL:         parseDuration(env("SESSION_TTL_HOURS", "24")),
			CookieSecure:       envBool("COOKIE_SECURE", false),
		},
	}
	if c.DBPass == "" {
		return nil, fmt.Errorf("DB_PASS is required (see .env.example)")
	}
	if c.Auth.GitHubClientID == "" || c.Auth.GitHubClientSecret == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET are required")
	}
	return c, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes"
}

func parseDuration(hoursStr string) time.Duration {
	hours := 24
	if hoursStr != "" {
		var h int
		if _, err := fmt.Sscanf(hoursStr, "%d", &h); err == nil && h > 0 {
			hours = h
		}
	}
	return time.Duration(hours) * time.Hour
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}