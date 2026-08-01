package content

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root, rel string, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadData(t *testing.T) {
	root := t.TempDir()

	writeFixture(t, root, "classes/sorcerer.yaml", `
id: sorcerer
name: Sorcerer
hitDie: d6
spellcaster: true
subclassLevel: 3
primaryAbility: CHA
suggestedSpecies:
  - kalashtar
  - tiefling
suggestedBackgrounds:
  - sage
  - charlatan
data:
  description: "Raw magic"
`)

	writeFixture(t, root, "subclasses/aberrant.yaml", `
id: aberrant
classId: sorcerer
name: Aberrant Mind
levelRequired: 3
recommendedSpells:
  1: [mind-sliver, arms-of-hadar]
  3: [detect-thoughts]
data:
  description: "Psionic"
`)

	writeFixture(t, root, "species/kalashtar.yaml", `
id: kalashtar
name: Kalashtar
abilityScores:
  WIS: 2
  CHA: 1
data:
  traits:
    - name: Dual Mind
`)

	writeFixture(t, root, "backgrounds/sage.yaml", `
id: sage
name: Sage
data:
  skills: [arcana, history]
`)

	writeFixture(t, root, "spells/magic-missile.yaml", `
id: magic-missile
name: Magic Missile
level: 1
school: evocation
data:
  classes: [sorcerer, wizard]
`)

	writeFixture(t, root, "feats/war-caster.yaml", `
id: war-caster
name: War Caster
prerequisites:
  spellcasting: true
data:
  description: "..."
`)

	writeFixture(t, root, "features/font-of-magic.yaml", `
id: font-of-magic
classId: sorcerer
name: Font of Magic
level: 2
data:
  description: "..."
`)

	c, err := LoadData(root)
	if err != nil {
		t.Fatalf("LoadData: %v", err)
	}

	if len(c.Classes) != 1 || c.Classes[0].ID != "sorcerer" {
		t.Fatalf("classes = %+v", c.Classes)
	}
	if len(c.Classes[0].SuggestedSpecies) != 2 || c.Classes[0].SuggestedSpecies[0] != "kalashtar" {
		t.Fatalf("suggestedSpecies = %v", c.Classes[0].SuggestedSpecies)
	}
	if len(c.Subclasses) != 1 || c.Subclasses[0].RecommendedSpells[1][0] != "mind-sliver" {
		t.Fatalf("subclasses = %+v", c.Subclasses)
	}
	if len(c.Species) != 1 || c.Species[0].AbilityScores["WIS"] != 2 {
		t.Fatalf("species = %+v", c.Species)
	}
	if len(c.Backgrounds) != 1 || len(c.Spells) != 1 || len(c.Feats) != 1 || len(c.Features) != 1 {
		t.Fatalf("counts mismatch: %+v", c)
	}
}

func TestLoadDataMissingDirs(t *testing.T) {
	c, err := LoadData(t.TempDir())
	if err != nil {
		t.Fatalf("LoadData on empty dir: %v", err)
	}
	if len(c.Classes) != 0 || len(c.Species) != 0 {
		t.Fatalf("expected empty content, got %+v", c)
	}
}

func TestClassJSONDataMergesSuggestions(t *testing.T) {
	c := Class{
		ID: "sorcerer", Name: "Sorcerer", HitDie: "d6", Spellcaster: true,
		SuggestedSpecies: []string{"kalashtar"},
		Data:             map[string]any{"description": "x"},
	}
	raw, err := c.JSONData()
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["description"] != "x" {
		t.Fatalf("data lost: %v", m)
	}
	species, ok := m["suggestedSpecies"].([]any)
	if !ok || len(species) != 1 || species[0] != "kalashtar" {
		t.Fatalf("suggestedSpecies not merged: %v", m["suggestedSpecies"])
	}
}

func TestParseDirRejectsInvalidYAML(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "classes/broken.yaml", "id: [unclosed\n  - a")
	if _, err := LoadData(root); err == nil {
		t.Fatal("LoadData: want error on invalid YAML, got nil")
	}
}
