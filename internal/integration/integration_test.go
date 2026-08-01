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
	srv := httpapi.New(db, logger, "test")
	ts := httptest.NewServer(srv.Router())
	t.Cleanup(ts.Close)
	return ts
}

func doJSON(t *testing.T, method, url string, body any) (int, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer res.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res.StatusCode, out
}

func TestHealth(t *testing.T) {
	ts := testServer(t)
	res, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("get health: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d, want 200", res.StatusCode)
	}
	var out map[string]string
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatalf("decode health: %v", err)
	}
	if out["status"] != "ok" {
		t.Fatalf("health status field = %q, want ok", out["status"])
	}
}

func TestContentEndToEnd(t *testing.T) {
	ts := testServer(t)
	code, out := doJSON(t, http.MethodGet, ts.URL+"/api/v1/content", nil)
	if code != http.StatusOK {
		t.Fatalf("content status = %d, want 200 (body: %v)", code, out)
	}
	classes, ok := out["classes"].([]any)
	if !ok || len(classes) == 0 {
		t.Fatalf("content.classes missing or empty: %v", out)
	}
	first, ok := classes[0].(map[string]any)
	if !ok || first["id"] == "" || first["hitDie"] == "" {
		t.Fatalf("content.classes[0] malformed: %v", classes[0])
	}
}

func TestBuildEndToEnd(t *testing.T) {
	ts := testServer(t)
	payload := map[string]any{
		"name":          "Integration",
		"classes":       []map[string]any{{"id": "sorcerer", "level": 6}},
		"speciesId":     "tiefling",
		"backgroundId":  "sage",
		"abilityScores": map[string]int{"STR": 8, "DEX": 14, "CON": 16, "INT": 12, "WIS": 10, "CHA": 18},
		"abilityMethod": "point-buy",
	}
	code, out := doJSON(t, http.MethodPost, ts.URL+"/api/v1/build", payload)
	if code != http.StatusOK {
		t.Fatalf("build status = %d, want 200 (body: %v)", code, out)
	}
	sheet, ok := out["sheet"].(map[string]any)
	if !ok {
		t.Fatalf("build response missing sheet: %v", out)
	}
	hp, ok := sheet["hp"].(map[string]any)
	if !ok || hp["max"].(float64) != 44 {
		t.Fatalf("sorcerer 6 hp.max = %v, want 44", hp["max"])
	}
}

func TestBuildInvalidClassEndToEnd(t *testing.T) {
	ts := testServer(t)
	payload := map[string]any{
		"name":          "Bogus",
		"classes":       []map[string]any{{"id": "wizard", "level": 1}},
		"abilityScores": map[string]int{"STR": 8, "DEX": 8, "CON": 8, "INT": 8, "WIS": 8, "CHA": 8},
	}
	code, out := doJSON(t, http.MethodPost, ts.URL+"/api/v1/build", payload)
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("build status = %d, want 422 (body: %v)", code, out)
	}
	apiErr, ok := out["error"].(map[string]any)
	if !ok || apiErr["code"] != "BUILD_ERROR" {
		t.Fatalf("build error object = %v, want BUILD_ERROR", out)
	}
}
