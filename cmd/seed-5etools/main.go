package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/hadnu/onatar/internal/config"
	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/store"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}

	dataDir := os.Getenv("DATA_DIR_5ETOOLS")
	if dataDir == "" {
		dataDir = "/home/hadnu/workspace/homelab/dnd-project/dnd-arts/5etools-src/data"
	}

	slog.Info("loading 5etools data", "dir", dataDir)

	fiveEData, err := load5eToolsData(dataDir)
	if err != nil {
		slog.Error("load 5etools data", "error", err)
		os.Exit(1)
	}

	onatarData, err := convertToOnatar(fiveEData)
	if err != nil {
		slog.Error("convert to onatar", "error", err)
		os.Exit(1)
	}

	// Print summary
	slog.Info("parsed 5etools data",
		"classes", len(onatarData.Classes),
		"subclasses", len(onatarData.Subclasses),
		"species", len(onatarData.Species),
		"backgrounds", len(onatarData.Backgrounds),
		"spells", len(onatarData.Spells),
		"feats", len(onatarData.Feats),
		"features", len(onatarData.Features),
	)

	// Skip DB if no password provided
	if cfg.DBPass == "" {
		slog.Info("skipping database seed (no DB_PASS)")
		return
	}

	db, err := store.Open(cfg.DSN())
	if err != nil {
		slog.Error("connect db", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := store.Seed(ctx, db, onatarData); err != nil {
		slog.Error("seed", "error", err)
		os.Exit(1)
	}

	slog.Info("seed complete",
		"classes", len(onatarData.Classes),
		"subclasses", len(onatarData.Subclasses),
		"species", len(onatarData.Species),
		"backgrounds", len(onatarData.Backgrounds),
		"spells", len(onatarData.Spells),
		"feats", len(onatarData.Feats),
		"features", len(onatarData.Features),
	)
}

type FiveEClass struct {
	Name                string   `json:"name"`
	Source              string   `json:"source"`
	Page                int      `json:"page"`
	SRD                 bool     `json:"srd"`
	HD                  HD       `json:"hd"`
	Proficiency         []string `json:"proficiency"`
	SpellcastingAbility string   `json:"spellcastingAbility"`
	CasterProgression   string   `json:"casterProgression"`
	StartingProficiencies struct {
		Weapons []any `json:"weapons"`
		Skills  []struct {
			Choose struct {
				From  []string `json:"from"`
				Count int      `json:"count"`
			} `json:"choose"`
		} `json:"skills"`
	} `json:"startingProficiencies"`
	StartingEquipment struct {
		Default []any `json:"default"`
	} `json:"startingEquipment"`
	Multiclassing struct {
		Requirements map[string]any `json:"requirements"`
	} `json:"multiclassing"`
	SubclassTitle string `json:"subclassTitle"`
	Fluff         struct {
		SubclassFluff []any `json:"subclassFluff"`
	} `json:"fluff"`
}

type HD struct {
	Number int `json:"number"`
	Faces  int `json:"faces"`
}

type FiveESubclass struct {
	Name              string   `json:"name"`
	ShortName         string   `json:"shortName"`
	Source            string   `json:"source"`
	ClassName         string   `json:"className"`
	ClassSource       string   `json:"classSource"`
	Page              int      `json:"page"`
	SRD               bool     `json:"srd"`
	SubclassFeatures  []string `json:"subclassFeatures"`
	AdditionalSpells  []any    `json:"additionalSpells"`
}

type FiveERace struct {
	Name           string `json:"name"`
	Source         string `json:"source"`
	Size           []string `json:"size"`
	Speed          any `json:"speed"`
	Ability        []map[string]any `json:"ability"`
	Entries        []any  `json:"entries"`
	Lineage        any `json:"lineage"`
	AdditionalSpells []any `json:"additionalSpells"`
}

type FiveEBackground struct {
	Name               string `json:"name"`
	Source             string `json:"source"`
	SkillProficiencies []map[string]any `json:"skillProficiencies"`
	ToolProficiencies  []map[string]any `json:"toolProficiencies"`
	LanguageProficiencies []map[string]any `json:"languageProficiencies"`
	StartingEquipment  []any `json:"startingEquipment"`
	Entries            []any `json:"entries"`
	Feats              []map[string]any `json:"feats"`
}

type FiveESpell struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Level       int    `json:"level"`
	School      string `json:"school"`
	Time        []struct {
		Number int    `json:"number"`
		Unit   string `json:"unit"`
	} `json:"time"`
	Range      map[string]any `json:"range"`
	Components any `json:"components"`
	Duration   []struct {
		Type     string `json:"type"`
		Duration struct {
			Type   string `json:"type"`
			Amount int    `json:"amount"`
		} `json:"duration"`
	} `json:"duration"`
	Entries         []any `json:"entries"`
	DamageInflict   []string `json:"damageInflict"`
	SavingThrow     []string `json:"savingThrow"`
	MiscTags        []string `json:"miscTags"`
	AreaTags        []string `json:"areaTags"`
	ScalingLevelDice any `json:"scalingLevelDice"`
	Classes         []struct {
		Name string `json:"name"`
		Source string `json:"source"`
	} `json:"classes"`
}

type FiveEFeat struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Page        int    `json:"page"`
	Prerequisite []any `json:"prerequisite"`
	Entries     []any  `json:"entries"`
}

type FiveEData struct {
	Classes     []FiveEClass     `json:"class"`
	Subclasses  []FiveESubclass  `json:"subclass"`
	Races       []FiveERace      `json:"race"`
	Backgrounds []FiveEBackground `json:"background"`
	Spells      []FiveESpell     `json:"spell"`
	Feats       []FiveEFeat      `json:"feat"`
}

func load5eToolsData(dir string) (*FiveEData, error) {
	data := &FiveEData{}

	classFiles, _ := filepath.Glob(filepath.Join(dir, "class", "class-*.json"))
	for _, f := range classFiles {
		var wrapper struct {
			Class []FiveEClass `json:"class"`
			Subclass []FiveESubclass `json:"subclass"`
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return nil, err
		}
		data.Classes = append(data.Classes, wrapper.Class...)
		data.Subclasses = append(data.Subclasses, wrapper.Subclass...)
	}

	racesRaw, err := os.ReadFile(filepath.Join(dir, "races.json"))
	if err != nil {
		return nil, err
	}
	var racesWrapper struct {
		Race []FiveERace `json:"race"`
	}
	if err := json.Unmarshal(racesRaw, &racesWrapper); err != nil {
		return nil, err
	}
	data.Races = racesWrapper.Race

	backgroundsRaw, err := os.ReadFile(filepath.Join(dir, "backgrounds.json"))
	if err != nil {
		return nil, err
	}
	var bgWrapper struct {
		Background []FiveEBackground `json:"background"`
	}
	if err := json.Unmarshal(backgroundsRaw, &bgWrapper); err != nil {
		return nil, err
	}
	data.Backgrounds = bgWrapper.Background

	spellFiles, _ := filepath.Glob(filepath.Join(dir, "spells", "spells-*.json"))
	slog.Info("found spell files", "count", len(spellFiles))
	for _, f := range spellFiles {
		var wrapper struct {
			Spell []FiveESpell `json:"spell"`
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			slog.Warn("read spell file", "file", f, "error", err)
			continue
		}
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			slog.Warn("parse spell file", "file", f, "error", err)
			continue
		}
		slog.Info("loaded spells", "file", f, "count", len(wrapper.Spell))
		data.Spells = append(data.Spells, wrapper.Spell...)
	}

	featsRaw, err := os.ReadFile(filepath.Join(dir, "feats.json"))
	if err != nil {
		return nil, err
	}
	var featsWrapper struct {
		Feat []FiveEFeat `json:"feat"`
	}
	if err := json.Unmarshal(featsRaw, &featsWrapper); err != nil {
		return nil, err
	}
	data.Feats = featsWrapper.Feat

	return data, nil
}

func convertToOnatar(data *FiveEData) (*content.Content, error) {
	onatar := &content.Content{}

	// Convert Classes
	classMap := make(map[string]bool)
	for _, c := range data.Classes {
		if classMap[c.Name] {
			continue
		}
		classMap[c.Name] = true

		id := toSlug(c.Name)
		hitDie := fmt.Sprintf("d%d", c.HD.Faces)

		suggestedSpecies := []string{}
		suggestedBackgrounds := []string{}

		onatar.Classes = append(onatar.Classes, content.Class{
			ID:                  id,
			Name:                c.Name,
			HitDie:              hitDie,
			Spellcaster:         c.SpellcastingAbility != "",
			SubclassLevel:       3,
			PrimaryAbility:      upperCaseFirst(c.SpellcastingAbility),
			SuggestedSpecies:    suggestedSpecies,
			SuggestedBackground: suggestedBackgrounds,
			Data: map[string]any{
				"description":       "",
				"savingThrows":      c.Proficiency,
				"skills":            c.StartingProficiencies.Skills,
				"weapons":           c.StartingProficiencies.Weapons,
				"startingEquipment": c.StartingEquipment.Default,
				"multiclassing":     c.Multiclassing.Requirements,
				"subclassTitle":     c.SubclassTitle,
			},
		})
	}

	// Convert Subclasses
	for _, s := range data.Subclasses {
		if s.ClassName == "" {
			continue
		}
		classID := toSlug(s.ClassName)
		subclassID := toSlug(s.Name)

		recommendedSpells := make(map[int][]string)
		if len(s.AdditionalSpells) > 0 {
			for _, asp := range s.AdditionalSpells {
				if spellMap, ok := asp.(map[string]any); ok {
					if known, ok := spellMap["known"].(map[string]any); ok {
						for levelStr, spells := range known {
							var level int
							if _, err := fmt.Sscanf(levelStr, "%d", &level); err == nil {
								if arr, ok := spells.([]any); ok {
									for _, sp := range arr {
										if spStr, ok := sp.(string); ok {
											recommendedSpells[level] = append(recommendedSpells[level], spStr)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		onatar.Subclasses = append(onatar.Subclasses, content.Subclass{
			ID:                subclassID,
			ClassID:           classID,
			Name:              s.Name,
			LevelRequired:     3,
			RecommendedSpells: recommendedSpells,
			Data: map[string]any{
				"description":      "",
				"subclassFeatures": s.SubclassFeatures,
			},
		})
	}

	// Convert Species (Races)
	raceMap := make(map[string]bool)
	for _, r := range data.Races {
		key := r.Name + "|" + r.Source
		if raceMap[key] {
			continue
		}
		raceMap[key] = true

		id := toSlug(r.Name)
		abilityScores := make(map[string]int)
		for _, a := range r.Ability {
			for k, v := range a {
				if val, ok := v.(float64); ok {
					abilityScores[upperCase(k)] += int(val)
				} else if val, ok := v.(int); ok {
					abilityScores[upperCase(k)] += val
				}
			}
		}

		traits := []map[string]any{}
		for _, e := range r.Entries {
			if entryMap, ok := e.(map[string]any); ok {
				if name, ok := entryMap["name"].(string); ok {
					traits = append(traits, map[string]any{
						"name":        name,
						"description": fmt.Sprintf("%v", entryMap["entries"]),
					})
				}
			}
		}

		onatar.Species = append(onatar.Species, content.Species{
			ID:            id,
			Name:          r.Name,
			AbilityScores: abilityScores,
			Data: map[string]any{
				"traits": traits,
				"size":   r.Size,
				"speed":  r.Speed,
				"lineage": r.Lineage,
			},
		})
	}

	// Convert Backgrounds
	bgMap := make(map[string]bool)
	for _, b := range data.Backgrounds {
		key := b.Name + "|" + b.Source
		if bgMap[key] {
			continue
		}
		bgMap[key] = true

		id := toSlug(b.Name)
		skills := []string{}
		for _, sp := range b.SkillProficiencies {
			for k := range sp {
				skills = append(skills, k)
			}
		}

		onatar.Backgrounds = append(onatar.Backgrounds, content.Background{
			ID:   id,
			Name: b.Name,
			Data: map[string]any{
				"skills":             skills,
				"toolProficiencies":  b.ToolProficiencies,
				"languageProficiencies": b.LanguageProficiencies,
				"startingEquipment":  b.StartingEquipment,
				"feature":            b.Entries,
				"feats":              b.Feats,
			},
		})
	}

	// Convert Spells
	spellMap := make(map[string]bool)
	for _, s := range data.Spells {
		key := s.Name + "|" + s.Source
		if spellMap[key] {
			continue
		}
		spellMap[key] = true

		id := toSlug(s.Name)
		school := schoolFullName(s.School)

		classes := []string{}
		for _, c := range s.Classes {
			classes = append(classes, toSlug(c.Name))
		}

		onatar.Spells = append(onatar.Spells, content.Spell{
			ID:     id,
			Name:   s.Name,
			Level:  s.Level,
			School: school,
			Data: map[string]any{
				"description":       fmt.Sprintf("%v", s.Entries),
				"classes":           classes,
				"components":        s.Components,
				"duration":          s.Duration,
				"range":             s.Range,
				"castingTime":       s.Time,
				"damageInflict":     s.DamageInflict,
				"savingThrow":       s.SavingThrow,
				"scalingLevelDice":  s.ScalingLevelDice,
				"miscTags":          s.MiscTags,
				"areaTags":          s.AreaTags,
			},
		})
	}

	// Convert Feats
	featMap := make(map[string]bool)
	for _, f := range data.Feats {
		key := f.Name + "|" + f.Source
		if featMap[key] {
			continue
		}
		featMap[key] = true

		id := toSlug(f.Name)
		prereqs := map[string]any{}
		if len(f.Prerequisite) > 0 {
			prereqs["raw"] = f.Prerequisite
		}

		onatar.Feats = append(onatar.Feats, content.Feat{
			ID:            id,
			Name:          f.Name,
			Prerequisites: prereqs,
			Data: map[string]any{
				"description": fmt.Sprintf("%v", f.Entries),
			},
		})
	}

	// Features are derived from class/subclass features during build, not seeded separately
	// for now we'll leave this empty

	return onatar, nil
}

func toSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "—", "-")
	s = strings.ReplaceAll(s, "–", "-")
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	return s
}

func upperCaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func upperCase(s string) string {
	return strings.ToUpper(s)
}

func schoolFullName(s string) string {
	switch s {
	case "A":
		return "abjuration"
	case "C":
		return "conjuration"
	case "D":
		return "divination"
	case "E":
		return "enchantment"
	case "V":
		return "evocation"
	case "I":
		return "illusion"
	case "N":
		return "necromancy"
	case "T":
		return "transmutation"
	default:
		return "unknown"
	}
}