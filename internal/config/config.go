package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string
	DBPass   string
	HTTPAddr string
}

func Load() (*Config, error) {
	c := &Config{
		DBHost:   env("DB_HOST", "127.0.0.1"),
		DBPort:   env("DB_PORT", "3306"),
		DBName:   env("DB_NAME", "onatar"),
		DBUser:   env("DB_USER", "onatar"),
		DBPass:   env("DB_PASS", ""),
		HTTPAddr: env("HTTP_ADDR", ":8090"),
	}
	if c.DBPass == "" {
		return nil, fmt.Errorf("DB_PASS is required (see .env.example)")
	}
	return c, nil
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
