package content

import (
	"encoding/json"
	"testing"
)

func TestSubclassJSONData(t *testing.T) {
	s := Subclass{
		ID: "aberrant-sorcery", ClassID: "sorcerer", Name: "Aberrant Sorcery",
		RecommendedSpells: map[int][]string{1: {"arms-of-hadar"}},
		Data:              map[string]any{"description": "x"},
	}
	raw, err := s.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["recommendedSpells"].(map[string]any); !ok {
		t.Fatalf("recommendedSpells not merged: %v", m)
	}
}

func TestSpeciesJSONData(t *testing.T) {
	s := Species{
		ID: "tiefling", Name: "Tiefling",
		AbilityScores: map[string]int{"CHA": 2},
		Data:          map[string]any{"size": "Medium"},
	}
	raw, err := s.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["size"] != "Medium" {
		t.Fatalf("data lost: %v", m)
	}
	as, ok := m["abilityScores"].(map[string]any)
	if !ok || as["CHA"] != float64(2) {
		t.Fatalf("abilityScores not merged: %v", m["abilityScores"])
	}
}

func TestBackgroundJSONData(t *testing.T) {
	b := Background{ID: "sage", Name: "Sage", Data: map[string]any{"skills": []string{"arcana"}}}
	raw, err := b.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) == "" || len(raw) == 0 {
		t.Fatal("empty JSON")
	}
}

func TestSpellJSONData(t *testing.T) {
	s := Spell{ID: "fire-bolt", Name: "Fire Bolt", Level: 0, School: "evocation",
		Data: map[string]any{"range": "120 feet"}}
	raw, err := s.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["range"] != "120 feet" {
		t.Fatalf("data lost: %v", m)
	}
}

func TestFeatJSONData(t *testing.T) {
	f := Feat{ID: "war-caster", Name: "War Caster",
		Prerequisites: map[string]any{"spellcasting": true},
		Data:          map[string]any{"description": "x"}}
	raw, err := f.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) == 0 {
		t.Fatal("empty JSON")
	}
}

func TestFeatureJSONData(t *testing.T) {
	f := Feature{ID: "metamagic", ClassID: "sorcerer", Name: "Metamagic", Level: 3,
		Data: map[string]any{"options": []string{"Quickened Spell"}}}
	raw, err := f.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["options"].([]any); !ok {
		t.Fatalf("options lost: %v", m)
	}
}

func TestParseDirSkipsNonYAML(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "classes/notes.txt", "not yaml")
	writeFixture(t, root, "classes/sorcerer.yaml", "id: sorcerer\nname: Sorcerer\n")
	c, err := LoadData(root)
	if err != nil {
		t.Fatalf("LoadData: %v", err)
	}
	if len(c.Classes) != 1 {
		t.Fatalf("classes = %d, want 1", len(c.Classes))
	}
}

func TestJSONDataEmptyData(t *testing.T) {
	c := Class{ID: "sorcerer", Name: "Sorcerer"}
	raw, err := c.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Fatalf("expected empty object, got %v", m)
	}
}
