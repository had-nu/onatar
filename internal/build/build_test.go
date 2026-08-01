package build

import (
	"testing"

	"github.com/hadnu/onatar/internal/content"
)

func rulesFixture() *content.Content {
	return &content.Content{
		Classes: []content.Class{
			{ID: "sorcerer", Name: "Sorcerer", HitDie: "d6", Spellcaster: true, SubclassLevel: 3,
				Data: map[string]any{"primaryAbility": "CHA"}},
			{ID: "fighter", Name: "Fighter", HitDie: "d10", Spellcaster: false, SubclassLevel: 3,
				Data: map[string]any{"primaryAbility": "STR"}},
		},
		Subclasses: []content.Subclass{
			{ID: "aberrant-sorcery", ClassID: "sorcerer", Name: "Aberrant Sorcery", LevelRequired: 3},
		},
		Species: []content.Species{{ID: "tiefling", Name: "Tiefling"}},
		Backgrounds: []content.Background{{ID: "sage", Name: "Sage"}},
		Spells: []content.Spell{
			{ID: "magic-missile", Name: "Magic Missile", Level: 1, School: "evocation",
				Data: map[string]any{"classes": []string{"sorcerer"}}},
			{ID: "fireball", Name: "Fireball", Level: 3, School: "evocation",
				Data: map[string]any{"classes": []string{"sorcerer"}}},
		},
		Feats: []content.Feat{
			{ID: "war-caster", Name: "War Caster", Prerequisites: map[string]any{"spellcasting": true}},
			{ID: "tough", Name: "Tough"},
		},
		Features: []content.Feature{
			{ID: "font-of-magic", ClassID: "sorcerer", Name: "Font of Magic", Level: 2,
				Data: map[string]any{"description": "Gain sorcery points."}},
			{ID: "metamagic", ClassID: "sorcerer", Name: "Metamagic", Level: 3,
				Data: map[string]any{"description": "Twist spells.", "choose": 2, "options": []string{"Careful Spell", "Quickened Spell"}}},
			{ID: "aberrant-power", SubclassID: "aberrant-sorcery", Name: "Aberrant Power", Level: 3,
				Data: map[string]any{"description": "Psionic power."}},
		},
	}
}

func baseReq() Request {
	return Request{
		Name:         "Onatar",
		Classes:      []ClassInput{{ID: "sorcerer", Level: 6, SubclassID: "aberrant-sorcery"}},
		SpeciesID:    "tiefling",
		BackgroundID: "sage",
		AbilityScores: map[string]int{"STR": 8, "DEX": 14, "CON": 16, "INT": 10, "WIS": 12, "CHA": 18},
		AbilityMethod: "point-buy",
		Skills:        []string{"arcana"},
		Spells:        []string{"magic-missile", "fireball"},
		Feats:         []string{"war-caster"},
	}
}

func TestComputeLevelProficiencyHPAC(t *testing.T) {
	resp, err := Compute(rulesFixture(), baseReq())
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	s := resp.Sheet
	if s.Level != 6 {
		t.Fatalf("level = %d, want 6", s.Level)
	}
	if s.ProficiencyBonus != 3 {
		t.Fatalf("PB = %d, want 3", s.ProficiencyBonus)
	}
	// d6: level 1 = 6 + CON3; levels 2-6 = 5*(4+3) => 9 + 35 = 44.
	if s.HP.Max != 44 || s.HP.Current != 44 {
		t.Fatalf("hp = %+v, want 44/44", s.HP)
	}
	// DEX 14 (+2) -> AC 12.
	if s.AC != 12 {
		t.Fatalf("ac = %d, want 12", s.AC)
	}
	// Sorcerer (full caster) level 6 slots.
	if got, want := s.SpellSlots, []int{4, 3, 3, 0, 0, 0, 0, 0, 0}; !eqSlice(got, want) {
		t.Fatalf("spellSlots = %v, want %v", got, want)
	}
	if m := s.Abilities["STR"]; m.Score != 8 || m.Mod != -1 {
		t.Fatalf("STR = %+v, want score 8 mod -1", m)
	}
	if m := s.Abilities["DEX"]; m.Mod != 2 {
		t.Fatalf("DEX mod = %d, want 2", m.Mod)
	}
	if m := s.Abilities["WIS"]; m.Mod != 1 {
		t.Fatalf("WIS mod = %d, want 1 (score 12)", m.Mod)
	}
	if m := s.Abilities["CHA"]; m.Mod != 4 {
		t.Fatalf("CHA mod = %d, want 4 (score 18)", m.Mod)
	}
}

func TestComputeFeaturesAndPendingChoices(t *testing.T) {
	resp, err := Compute(rulesFixture(), baseReq())
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	feats := resp.Sheet.Features
	names := map[string]bool{}
	for _, f := range feats {
		names[f.Name] = true
	}
	for _, want := range []string{"Font of Magic", "Metamagic", "Aberrant Power"} {
		if !names[want] {
			t.Fatalf("missing feature %q in %v", want, names)
		}
	}
	choices := resp.Sheet.PendingChoices
	if len(choices) != 1 {
		t.Fatalf("pendingChoices = %v, want 1", choices)
	}
	if choices[0].Type != "metamagic" {
		t.Fatalf("choice type = %q, want metamagic", choices[0].Type)
	}
	want := "Choose 2 Metamagic options: Careful Spell, Quickened Spell"
	if choices[0].Description != want {
		t.Fatalf("choice desc = %q, want %q", choices[0].Description, want)
	}
}

func TestComputeNonCaster(t *testing.T) {
	req := baseReq()
	req.Classes = []ClassInput{{ID: "fighter", Level: 3}}
	req.Spells = nil
	req.Feats = nil
	resp, err := Compute(rulesFixture(), req)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	s := resp.Sheet
	if s.Level != 3 {
		t.Fatalf("level = %d, want 3", s.Level)
	}
	// d10: 10 + CON3 + 2*(6+3) = 13 + 18 = 31.
	if s.HP.Max != 31 {
		t.Fatalf("hp = %d, want 31", s.HP.Max)
	}
	if got, want := s.SpellSlots, []int{0, 0, 0, 0, 0, 0, 0, 0, 0}; !eqSlice(got, want) {
		t.Fatalf("non-caster spellSlots = %v, want zeros", got)
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*Request)
		kind ErrorKind
	}{
		{"no classes", func(r *Request) { r.Classes = nil }, KindInvalidDraft},
		{"unknown class", func(r *Request) { r.Classes = []ClassInput{{ID: "wizard", Level: 1}} }, KindBuildError},
		{"level out of range", func(r *Request) { r.Classes[0].Level = 21 }, KindInvalidDraft},
		{"unknown subclass", func(r *Request) { r.Classes[0].SubclassID = "nope" }, KindBuildError},
		{"subclass wrong class", func(r *Request) { r.Classes[0].SubclassID = "aberrant-sorcery"; r.Classes = []ClassInput{{ID: "fighter", Level: 6, SubclassID: "aberrant-sorcery"}} }, KindBuildError},
		{"subclass level not reached", func(r *Request) { r.Classes[0].Level = 2 }, KindBuildError},
		{"total level exceeds 20", func(r *Request) { r.Classes = []ClassInput{{ID: "sorcerer", Level: 12}, {ID: "fighter", Level: 12}} }, KindInvalidDraft},
		{"unknown species", func(r *Request) { r.SpeciesID = "elf" }, KindBuildError},
		{"unknown background", func(r *Request) { r.BackgroundID = "noble" }, KindBuildError},
		{"ability out of range", func(r *Request) { r.AbilityScores["CHA"] = 21 }, KindInvalidDraft},
		{"missing ability", func(r *Request) { delete(r.AbilityScores, "CON") }, KindInvalidDraft},
		{"spell not for class", func(r *Request) { r.Spells = []string{"magic-missile"}; r.Classes = []ClassInput{{ID: "fighter", Level: 6}} }, KindInvalidSpell},
		{"spell level unavailable", func(r *Request) { r.Spells = []string{"fireball"}; r.Classes = []ClassInput{{ID: "sorcerer", Level: 2}} }, KindInvalidSpell},
		{"unknown spell", func(r *Request) { r.Spells = []string{"wish"} }, KindBuildError},
		{"spell without caster", func(r *Request) { r.Spells = []string{"magic-missile"}; r.Classes = []ClassInput{{ID: "fighter", Level: 3}} }, KindInvalidSpell},
		{"unknown feat", func(r *Request) { r.Feats = []string{"great-weapon-master"} }, KindBuildError},
		{"feat prereq unmet", func(r *Request) { r.Spells = nil; r.Feats = []string{"war-caster"}; r.Classes = []ClassInput{{ID: "fighter", Level: 6}} }, KindBuildError},
		{"too many classes", func(r *Request) {
			r.Classes = []ClassInput{{ID: "sorcerer", Level: 2}, {ID: "fighter", Level: 2}, {ID: "sorcerer", Level: 2},
				{ID: "fighter", Level: 2}, {ID: "sorcerer", Level: 2}, {ID: "fighter", Level: 2},
				{ID: "sorcerer", Level: 2}, {ID: "fighter", Level: 2}, {ID: "sorcerer", Level: 2},
				{ID: "fighter", Level: 2}, {ID: "sorcerer", Level: 2}}
		}, KindInvalidDraft},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := baseReq()
			tc.mut(&req)
			_, err := Compute(rulesFixture(), req)
			if err == nil {
				t.Fatalf("expected error")
			}
			re, ok := err.(*RuleError)
			if !ok {
				t.Fatalf("error type = %T, want *RuleError", err)
			}
			if re.Kind != tc.kind {
				t.Fatalf("kind = %q, want %q (%s)", re.Kind, tc.kind, re.Message)
			}
		})
	}
}

func TestProficiencyBonusTable(t *testing.T) {
	want := map[int]int{1: 2, 4: 2, 5: 3, 8: 3, 9: 4, 12: 4, 13: 5, 16: 5, 17: 6, 20: 6}
	for lvl, pb := range want {
		if got := proficiencyBonus(lvl); got != pb {
			t.Fatalf("PB(%d) = %d, want %d", lvl, got, pb)
		}
	}
}

func eqSlice(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
