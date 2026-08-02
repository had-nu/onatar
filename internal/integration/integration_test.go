// Package integration exercises the real API router against a real MariaDB
// (TEST_DB_DSN). Skips when the variable is unset so local `go test ./...` runs
// stay hermetic; CI sets it after migrate + seed (see .github/workflows/ci.yml).
package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hadnu/onatar/internal/config"
	"github.com/hadnu/onatar/internal/httpapi"
	"github.com/hadnu/onatar/internal/store"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set; skipping integration test")
	}
	db, err := store.Open(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		Auth: config.AuthConfig{
			GitHubClientID:     "test",
			GitHubClientSecret: "test",
			GitHubRedirectURL:  "http://localhost:5173/auth/callback",
			JWTSecret:          "test-secret",
			SessionCookieName:  "onatar_session",
			SessionTTL:         24 * 3600 * 1000000000,
			CookieSecure:       false,
		},
	}
	srv := httpapi.New(db, logger, "test", cfg)
	ts := httptest.NewServer(srv.Router())
	t.Cleanup(ts.Close)
	return ts
}

func doJSON(t *testing.T, method, url string, body any) (int, map[string]any) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp.StatusCode, out
}

func TestHealth(t *testing.T) {
	ts := testServer(t)
	code, out := doJSON(t, "GET", ts.URL+"/health", nil)
	if code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", code)
	}
	if out["status"] != "ok" {
		t.Fatalf("health status field = %q, want ok", out["status"])
	}
}

func TestContentEndToEnd(t *testing.T) {
	ts := testServer(t)
	code, out := doJSON(t, "GET", ts.URL+"/api/v1/content", nil)
	if code != http.StatusOK {
		t.Fatalf("content status = %d, want 200", code)
	}
	if out["classes"] == nil {
		t.Fatal("content.classes missing or empty")
	}
}

func TestBuildEndToEnd(t *testing.T) {
	ts := testServer(t)
	payload := map[string]any{
		"name":  "Test Hero",
		"classes": []map[string]any{{"id": "fighter", "level": 1}},
		"speciesId": "human",
		"backgroundId": "sage",
		"abilityScores": map[string]int{"STR": 16, "DEX": 13, "CON": 14, "INT": 10, "WIS": 12, "CHA": 8},
		"abilityMethod": "point-buy",
	}
	code, out := doJSON(t, "POST", ts.URL+"/api/v1/build", payload)
	if code != http.StatusOK {
		t.Fatalf("build status = %d, want 200, body: %v", code, out)
	}
	if out["sheet"] == nil {
		t.Fatal("sheet missing in build response")
	}
}