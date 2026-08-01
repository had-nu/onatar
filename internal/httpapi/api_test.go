package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hadnu/onatar/internal/content"
)

type stubLoader struct {
	c   *content.Content
	err error
}

func (s stubLoader) LoadContent(context.Context) (*content.Content, error) {
	return s.c, s.err
}

func stubContent() *content.Content {
	return &content.Content{
		Classes: []content.Class{
			{ID: "sorcerer", Name: "Sorcerer", HitDie: "d6", Spellcaster: true, SubclassLevel: 3,
				Data: map[string]any{"primaryAbility": "CHA"}},
			{ID: "fighter", Name: "Fighter", HitDie: "d10", Spellcaster: false, SubclassLevel: 3},
		},
		Subclasses: []content.Subclass{
			{ID: "aberrant-sorcery", ClassID: "sorcerer", Name: "Aberrant Sorcery", LevelRequired: 3},
		},
		Species:     []content.Species{{ID: "tiefling", Name: "Tiefling"}},
		Backgrounds: []content.Background{{ID: "sage", Name: "Sage"}},
		Spells:      []content.Spell{{ID: "magic-missile", Name: "Magic Missile", Level: 1, School: "evocation", Data: map[string]any{"classes": []string{"sorcerer"}}}},
		Feats:       []content.Feat{{ID: "war-caster", Name: "War Caster", Prerequisites: map[string]any{"spellcasting": true}}},
		Features: []content.Feature{
			{ID: "font-of-magic", ClassID: "sorcerer", Name: "Font of Magic", Level: 2},
		},
	}
}

func testServer(t *testing.T, l contentLoader) *Server {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := &Server{
		version:       "test",
		logger:        logger,
		contentLoader: l,
		limiter:       newTokenBucket(defaultRatePerMinute, defaultCapacity),
	}
	t.Cleanup(s.limiter.close)
	return s
}

func doReq(t *testing.T, h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, r)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHealth(t *testing.T) {
	s := testServer(t, stubLoader{})
	rec := doReq(t, s.Router(), http.MethodGet, "/health", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" || body["version"] != "test" {
		t.Fatalf("body = %v", body)
	}
}

func TestContentOK(t *testing.T) {
	s := testServer(t, stubLoader{c: stubContent()})
	rec := doReq(t, s.Router(), http.MethodGet, "/api/v1/content", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"classes", "subclasses", "species", "backgrounds", "spells", "feats", "features"} {
		if _, ok := body[key]; !ok {
			t.Fatalf("missing top-level key %q", key)
		}
	}
}

func TestContentDBError(t *testing.T) {
	s := testServer(t, stubLoader{err: context.DeadlineExceeded})
	rec := doReq(t, s.Router(), http.MethodGet, "/api/v1/content", "")
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	assertErrorCode(t, rec, "INTERNAL_ERROR")
}

func TestBuildOK(t *testing.T) {
	s := testServer(t, stubLoader{c: stubContent()})
	body := `{
		"name":"Onatar",
		"classes":[{"id":"sorcerer","level":6,"subclassId":"aberrant-sorcery"}],
		"speciesId":"tiefling",
		"backgroundId":"sage",
		"abilityScores":{"STR":8,"DEX":14,"CON":16,"INT":10,"WIS":12,"CHA":18},
		"abilityMethod":"point-buy",
		"spells":["magic-missile"],
		"feats":["war-caster"]
	}`
	rec := doReq(t, s.Router(), http.MethodPost, "/api/v1/build", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Sheet struct {
			Level  int `json:"level"`
			HP     struct{ Max int } `json:"hp"`
			AC     int `json:"ac"`
			PB     int `json:"proficiencyBonus"`
			Slots  []int `json:"spellSlots"`
		} `json:"sheet"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Sheet.Level != 6 || resp.Sheet.HP.Max != 44 || resp.Sheet.AC != 12 || resp.Sheet.PB != 3 {
		t.Fatalf("sheet = %+v", resp.Sheet)
	}
}

func TestBuildValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
		code string
	}{
		{"malformed json", `{`, "INVALID_DRAFT"},
		{"unknown class", `{"classes":[{"id":"wizard","level":1}]}`, "BUILD_ERROR"},
		{"spell not for class", validBody(`[{"id":"fighter","level":6}]`, `["magic-missile"]`, `[]`), "INVALID_SPELL_SELECTION"},
		{"feat prereq unmet", validBody(`[{"id":"fighter","level":6}]`, `[]`, `["war-caster"]`), "BUILD_ERROR"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := testServer(t, stubLoader{c: stubContent()})
			rec := doReq(t, s.Router(), http.MethodPost, "/api/v1/build", tc.body)
			if rec.Code == http.StatusOK {
				t.Fatalf("expected error status, body = %s", rec.Body.String())
			}
			assertErrorCode(t, rec, tc.code)
		})
	}
}

func TestBuildRouteLimited(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	s := &Server{
		version:       "test",
		logger:        logger,
		contentLoader: stubLoader{c: stubContent()},
		limiter:       newTokenBucket(60, 2), // 1 token/sec, capacity 2
	}
	t.Cleanup(s.limiter.close)

	h := s.Router()
	body := validBody(`[{"id":"sorcerer","level":1}]`, `[]`, `[]`)
	for i := 0; i < 2; i++ {
		rec := doReq(t, h, http.MethodPost, "/api/v1/build", body)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want 200", i, rec.Code)
		}
	}
	rec := doReq(t, h, http.MethodPost, "/api/v1/build", body)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("3rd request status = %d, want 429", rec.Code)
	}
	assertErrorCode(t, rec, "RATE_LIMITED")
	if rec.Header().Get("Retry-After") == "" {
		t.Fatal("missing Retry-After header")
	}
}

func TestRateLimiterDirect(t *testing.T) {
	rl := newTokenBucket(60, 1)
	t.Cleanup(rl.close)
	if !rl.allow("a") {
		t.Fatal("first token should be allowed")
	}
	if rl.allow("a") {
		t.Fatal("second token should be denied (capacity 1)")
	}
	if !rl.allow("b") {
		t.Fatal("different key should be allowed")
	}
}

func assertErrorCode(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("error body: %v (%s)", err, rec.Body.String())
	}
	if body.Error.Code != want {
		t.Fatalf("error code = %q, want %q", body.Error.Code, want)
	}
}

// validBody builds a well-formed /build request with the given JSON fragments.
func validBody(classes, spells, feats string) string {
	return `{"name":"Onatar","classes":` + classes + `,
		"speciesId":"tiefling","backgroundId":"sage",
		"abilityScores":{"STR":8,"DEX":14,"CON":16,"INT":10,"WIS":12,"CHA":18},
		"abilityMethod":"point-buy",
		"spells":` + spells + `,"feats":` + feats + `}`
}
