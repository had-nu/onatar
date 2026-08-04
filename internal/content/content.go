package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Content struct {
	Classes     []Class     `json:"classes"`
	Subclasses  []Subclass  `json:"subclasses"`
	Species     []Species   `json:"species"`
	Backgrounds []Background `json:"backgrounds"`
	Spells      []Spell     `json:"spells"`
	Feats       []Feat      `json:"feats"`
	Features    []Feature   `json:"features"`
	Items       []Item      `json:"items"`
}

// Class is the normalized form of data/classes/<id>.yaml (PRD §8.1).
// JSON tags match the GET /content response schema (PRD §3.5).
type Class struct {
	ID                  string         `yaml:"id" json:"id"`
	Name                string         `yaml:"name" json:"name"`
	HitDie              string         `yaml:"hitDie" json:"hitDie"`
	Spellcaster         bool           `yaml:"spellcaster" json:"spellcaster"`
	SubclassLevel       int            `yaml:"subclassLevel" json:"subclassLevel"`
	PrimaryAbility      string         `yaml:"primaryAbility" json:"-"`
	SuggestedSpecies    []string       `yaml:"suggestedSpecies" json:"suggestedSpecies"`
	SuggestedBackground []string       `yaml:"suggestedBackgrounds" json:"suggestedBackgrounds"`
	Data                map[string]any `yaml:"data" json:"data"`
}

type Subclass struct {
	ID                string            `yaml:"id" json:"id"`
	ClassID           string            `yaml:"classId" json:"classId"`
	Name              string            `yaml:"name" json:"name"`
	LevelRequired     int               `yaml:"levelRequired" json:"levelRequired"`
	RecommendedSpells map[int][]string  `yaml:"recommendedSpells" json:"-"`
	Data              map[string]any    `yaml:"data" json:"data"`
}

type Species struct {
	ID            string         `yaml:"id" json:"id"`
	Name          string         `yaml:"name" json:"name"`
	AbilityScores map[string]int `yaml:"abilityScores" json:"-"`
	Data          map[string]any `yaml:"data" json:"data"`
}

type Background struct {
	ID   string         `yaml:"id" json:"id"`
	Name string         `yaml:"name" json:"name"`
	Data map[string]any `yaml:"data" json:"data"`
}

type Spell struct {
	ID     string         `yaml:"id" json:"id"`
	Name   string         `yaml:"name" json:"name"`
	Level  int            `yaml:"level" json:"level"`
	School string         `yaml:"school" json:"school"`
	Data   map[string]any `yaml:"data" json:"data"`
}

type Feat struct {
	ID            string         `yaml:"id" json:"id"`
	Name          string         `yaml:"name" json:"name"`
	Prerequisites map[string]any `yaml:"prerequisites" json:"prerequisites"`
	Data          map[string]any `yaml:"data" json:"data"`
}

type Feature struct {
	ID         string         `yaml:"id" json:"id"`
	ClassID    string         `yaml:"classId" json:"classId"`
	SubclassID string         `yaml:"subclassId" json:"subclassId"`
	Name       string         `yaml:"name" json:"name"`
	Level      int            `yaml:"level" json:"level"`
	Data       map[string]any `yaml:"data" json:"data"`
}

type Item struct {
	ID       string         `yaml:"id" json:"id"`
	Name     string         `yaml:"name" json:"name"`
	Type     string         `yaml:"type" json:"type"`
	Rarity   string         `yaml:"rarity" json:"rarity"`
	Source   string         `yaml:"source" json:"source"`
	Edition  string         `yaml:"edition" json:"edition"`
	Data     map[string]any `yaml:"data" json:"data"`
}

// JSONData serializes the flexible payload for the `data JSON` column,
// merging suggestion fields per PRD §8.1.
func (c Class) JSONData() ([]byte, error) {
	d := map[string]any{}
	for k, v := range c.Data {
		d[k] = v
	}
	if len(c.SuggestedSpecies) > 0 {
		d["suggestedSpecies"] = c.SuggestedSpecies
	}
	if len(c.SuggestedBackground) > 0 {
		d["suggestedBackgrounds"] = c.SuggestedBackground
	}
	if c.PrimaryAbility != "" {
		d["primaryAbility"] = c.PrimaryAbility
	}
	return json.Marshal(d)
}

func (s Subclass) JSONData() ([]byte, error) {
	d := map[string]any{}
	for k, v := range s.Data {
		d[k] = v
	}
	if len(s.RecommendedSpells) > 0 {
		d["recommendedSpells"] = s.RecommendedSpells
	}
	return json.Marshal(d)
}

func (s Species) JSONData() ([]byte, error) {
	d := map[string]any{}
	for k, v := range s.Data {
		d[k] = v
	}
	if len(s.AbilityScores) > 0 {
		d["abilityScores"] = s.AbilityScores
	}
	return json.Marshal(d)
}

func (b Background) JSONData() ([]byte, error) { return json.Marshal(b.Data) }
func (s Spell) JSONData() ([]byte, error)      { return json.Marshal(s.Data) }

func (f Feat) JSONData() ([]byte, error) { return json.Marshal(f.Data) }

func (f Feature) JSONData() ([]byte, error) { return json.Marshal(f.Data) }

func (i Item) JSONData() ([]byte, error) { return json.Marshal(i.Data) }

// LoadData reads the data/ directory tree and parses all content files.
func LoadData(root string) (*Content, error) {
	c := &Content{}
	var err error

	if c.Classes, err = parseDir[Class](filepath.Join(root, "classes")); err != nil {
		return nil, fmt.Errorf("classes: %w", err)
	}
	if c.Subclasses, err = parseDir[Subclass](filepath.Join(root, "subclasses")); err != nil {
		return nil, fmt.Errorf("subclasses: %w", err)
	}
	if c.Species, err = parseDir[Species](filepath.Join(root, "species")); err != nil {
		return nil, fmt.Errorf("species: %w", err)
	}
	if c.Backgrounds, err = parseDir[Background](filepath.Join(root, "backgrounds")); err != nil {
		return nil, fmt.Errorf("backgrounds: %w", err)
	}
	if c.Spells, err = parseDir[Spell](filepath.Join(root, "spells")); err != nil {
		return nil, fmt.Errorf("spells: %w", err)
	}
	if c.Feats, err = parseDir[Feat](filepath.Join(root, "feats")); err != nil {
		return nil, fmt.Errorf("feats: %w", err)
	}
	if c.Features, err = parseDir[Feature](filepath.Join(root, "features")); err != nil {
		return nil, fmt.Errorf("features: %w", err)
	}
	return c, nil
}

func parseDir[T any](dir string) ([]T, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, err
	}
	out := make([]T, 0, len(files))
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		var v T
		// #nosec G304 -- dir is DATA_DIR (admin-controlled, repo-local seed input).
		raw, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(raw, &v); err != nil {
			return nil, fmt.Errorf("%s: %w", f.Name(), err)
		}
		out = append(out, v)
	}
	return out, nil
}
