package config

import (
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("DB_PASS", "secret")

	c, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c.DBHost != "127.0.0.1" || c.DBPort != "3306" || c.DBName != "onatar" || c.DBUser != "onatar" {
		t.Fatalf("unexpected defaults: %+v", c)
	}
	if c.HTTPAddr != ":8090" {
		t.Fatalf("HTTPAddr = %q, want :8090", c.HTTPAddr)
	}
}

func TestLoadRequiresDBPass(t *testing.T) {
	t.Setenv("DB_PASS", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load: want error when DB_PASS empty, got nil")
	}
}

func TestDSN(t *testing.T) {
	c := &Config{
		DBHost: "db.local", DBPort: "3307", DBName: "onatar",
		DBUser: "u", DBPass: "p@ss",
	}
	dsn := c.DSN()
	if !strings.HasPrefix(dsn, "u:p@ss@tcp(db.local:3307)/onatar?") {
		t.Fatalf("unexpected DSN: %s", dsn)
	}
	if !strings.Contains(dsn, "parseTime=true") {
		t.Fatalf("DSN missing parseTime: %s", dsn)
	}
}
