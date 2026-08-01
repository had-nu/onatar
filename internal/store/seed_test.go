package store

import (
	"context"
	"os"
	"testing"

	"github.com/hadnu/onatar/internal/content"
)

// seedTestDSN is set in CI (migrate + seed smoke) or via TEST_DB_DSN locally.
func seedTestDSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set; skipping DB integration test")
	}
	return dsn
}

func TestSeedIdempotent(t *testing.T) {
	dsn := seedTestDSN(t)
	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()
	data := &content.Content{
		Classes: []content.Class{{ID: "sorcerer", Name: "Sorcerer", HitDie: "d6", Spellcaster: true}},
		Subclasses: []content.Subclass{
			{ID: "aberrant-sorcery", ClassID: "sorcerer", Name: "Aberrant Sorcery", LevelRequired: 3},
		},
		Species:     []content.Species{{ID: "tiefling", Name: "Tiefling"}},
		Backgrounds: []content.Background{{ID: "sage", Name: "Sage"}},
		Spells:      []content.Spell{{ID: "magic-missile", Name: "Magic Missile", Level: 1, School: "evocation"}},
		Feats:       []content.Feat{{ID: "war-caster", Name: "War Caster"}},
		Features: []content.Feature{
			{ID: "font-of-magic", ClassID: "sorcerer", Name: "Font of Magic", Level: 2},
			{ID: "no-class-feature", Name: "Generic", Level: 1},
		},
	}

	for i := 0; i < 2; i++ {
		if err := Seed(ctx, db, data); err != nil {
			t.Fatalf("Seed run %d: %v", i, err)
		}
	}

	tests := []struct {
		table, id string
	}{
		{"classes", "sorcerer"},
		{"subclasses", "aberrant-sorcery"},
		{"species", "tiefling"},
		{"backgrounds", "sage"},
		{"spells", "magic-missile"},
		{"feats", "war-caster"},
		{"features", "font-of-magic"},
	}
	for _, tc := range tests {
		var n int
		if err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM `+tc.table+` WHERE id=?`, tc.id).Scan(&n); err != nil {
			t.Fatalf("%s: %v", tc.table, err)
		}
		if n != 1 {
			t.Fatalf("%s rows = %d, want 1 (idempotent)", tc.table, n)
		}
	}

	var name string
	if err := db.QueryRowContext(ctx, `SELECT name FROM spells WHERE id='magic-missile'`).Scan(&name); err != nil {
		t.Fatal(err)
	}
	if name != "Magic Missile" {
		t.Fatalf("spell name = %q", name)
	}
}
