package store

import (
	"context"
	"testing"

	"github.com/hadnu/onatar/internal/content"
)

func TestContentQuery(t *testing.T) {
	dsn := seedTestDSN(t)
	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()
	data := &content.Content{
		Classes: []content.Class{
			{ID: "tst-sorcerer", Name: "Test Sorcerer", HitDie: "d6", Spellcaster: true, SubclassLevel: 3,
				PrimaryAbility: "CHA",
				SuggestedSpecies: []string{"tiefling"},
				Data:            map[string]any{"description": "inherent magic"}},
		},
		Subclasses: []content.Subclass{
			{ID: "tst-aberrant", ClassID: "tst-sorcerer", Name: "Test Aberrant", LevelRequired: 3,
				RecommendedSpells: map[int][]string{1: {"arms-of-hadar"}},
				Data:              map[string]any{"description": "Far Realm"}},
		},
		Species:     []content.Species{{ID: "tst-tiefling", Name: "Test Tiefling", Data: map[string]any{"size": "Medium"}}},
		Backgrounds: []content.Background{{ID: "tst-sage", Name: "Test Sage", Data: map[string]any{"skills": []string{"arcana"}}}},
		Spells:      []content.Spell{{ID: "tst-fire-bolt", Name: "Test Fire Bolt", Level: 0, School: "evocation", Data: map[string]any{"range": "120 feet"}}},
		Feats: []content.Feat{
			{ID: "tst-war-caster", Name: "Test War Caster", Prerequisites: map[string]any{"spellcasting": true}, Data: map[string]any{"description": "x"}},
		},
		Features: []content.Feature{
			{ID: "tst-font-of-magic", ClassID: "tst-sorcerer", Name: "Test Font of Magic", Level: 2, Data: map[string]any{"description": "points"}},
		},
	}
	if err := Seed(ctx, db, data); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	got, err := Content(ctx, db)
	if err != nil {
		t.Fatalf("Content: %v", err)
	}

	if len(got.Classes) == 0 || len(got.Subclasses) == 0 || len(got.Species) == 0 ||
		len(got.Backgrounds) == 0 || len(got.Spells) == 0 || len(got.Feats) == 0 || len(got.Features) == 0 {
		t.Fatalf("incomplete content: %+v", got)
	}

	var c content.Class
	for _, cc := range got.Classes {
		if cc.ID == "tst-sorcerer" {
			c = cc
			break
		}
	}
	if c.ID != "tst-sorcerer" || c.HitDie != "d6" || !c.Spellcaster || c.SubclassLevel != 3 {
		t.Fatalf("class round-trip: %+v", c)
	}
	if len(c.SuggestedSpecies) == 0 || c.SuggestedSpecies[0] != "tiefling" {
		t.Fatalf("suggestedSpecies split: %v", c.SuggestedSpecies)
	}
	if c.Data["description"] != "inherent magic" {
		t.Fatalf("class data: %v", c.Data)
	}
	if c.Data["primaryAbility"] != "CHA" {
		t.Fatalf("primaryAbility merged into data: %v", c.Data["primaryAbility"])
	}

	var s content.Subclass
	for _, ss := range got.Subclasses {
		if ss.ID == "tst-aberrant" {
			s = ss
			break
		}
	}
	if s.ClassID != "tst-sorcerer" || s.LevelRequired != 3 {
		t.Fatalf("subclass round-trip: %+v", s)
	}
	if _, ok := s.Data["recommendedSpells"]; !ok {
		t.Fatalf("recommendedSpells missing from data: %v", s.Data)
	}

	var f content.Feat
	for _, ff := range got.Feats {
		if ff.ID == "tst-war-caster" {
			f = ff
			break
		}
	}
	if f.Prerequisites["spellcasting"] != true {
		t.Fatalf("feat prerequisites: %v", f.Prerequisites)
	}
}
