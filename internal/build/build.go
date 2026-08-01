// Package build implements the POST /build rules engine (PRD §3.5).
// It is a pure function over content.Content: no I/O, fully unit-testable.
package build

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hadnu/onatar/internal/content"
)

const (
	// MaxClasses, MaxSpells, MaxFeats bound request sizes (PRD threat model §7).
	MaxClasses = 10
	MaxSpells  = 20
	MaxFeats   = 20
)

// ErrorKind maps to the API error schema (PRD §3.5 "Error Response Format").
type ErrorKind string

const (
	KindInvalidDraft  ErrorKind = "INVALID_DRAFT"
	KindInvalidSpell  ErrorKind = "INVALID_SPELL_SELECTION"
	KindBuildError    ErrorKind = "BUILD_ERROR"
)

// RuleError is a domain-level build failure with a stable error code.
type RuleError struct {
	Kind    ErrorKind
	Message string
	Details map[string]any
}

func (e *RuleError) Error() string { return e.Message }

func invalidDraftf(format string, args ...any) *RuleError {
	return &RuleError{Kind: KindInvalidDraft, Message: fmt.Sprintf(format, args...)}
}

func buildErrorf(format string, args ...any) *RuleError {
	return &RuleError{Kind: KindBuildError, Message: fmt.Sprintf(format, args...)}
}

func invalidSpellf(format string, args ...any) *RuleError {
	return &RuleError{Kind: KindInvalidSpell, Message: fmt.Sprintf(format, args...)}
}

// Request mirrors the BuildRequest schema.
type Request struct {
	Name          string       `json:"name"`
	Classes       []ClassInput `json:"classes"`
	SpeciesID     string       `json:"speciesId"`
	BackgroundID  string       `json:"backgroundId"`
	AbilityScores map[string]int `json:"abilityScores"`
	AbilityMethod string       `json:"abilityMethod"`
	Skills        []string     `json:"skills"`
	Spells        []string     `json:"spells"`
	Feats         []string     `json:"feats"`
	IsNPC         bool         `json:"isNpc"`
}

type ClassInput struct {
	ID         string `json:"id"`
	Level      int    `json:"level"`
	SubclassID string `json:"subclassId"`
}

// Response mirrors the BuildResponse schema.
type Response struct {
	Sheet Sheet `json:"sheet"`
}

type Sheet struct {
	Level            int                  `json:"level"`
	HP               HP                   `json:"hp"`
	AC               int                  `json:"ac"`
	ProficiencyBonus int                  `json:"proficiencyBonus"`
	SpellSlots       []int                `json:"spellSlots"`
	Abilities        map[string]Ability   `json:"abilities"`
	Features         []SheetFeature       `json:"features"`
	PendingChoices   []PendingChoice      `json:"pendingChoices"`
}

type HP struct {
	Max     int `json:"max"`
	Current int `json:"current"`
}

type Ability struct {
	Score int `json:"score"`
	Mod   int `json:"mod"`
}

type SheetFeature struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type PendingChoice struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// abilityOrder fixes the JSON key order of the six ability scores.
var abilityOrder = []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}

// spellSlotsByCasterLevel is the full-caster progression (PHB 2024).
var spellSlotsByCasterLevel = [21][]int{
	0:  {0, 0, 0, 0, 0, 0, 0, 0, 0},
	1:  {2, 0, 0, 0, 0, 0, 0, 0, 0},
	2:  {3, 0, 0, 0, 0, 0, 0, 0, 0},
	3:  {4, 2, 0, 0, 0, 0, 0, 0, 0},
	4:  {4, 3, 0, 0, 0, 0, 0, 0, 0},
	5:  {4, 3, 2, 0, 0, 0, 0, 0, 0},
	6:  {4, 3, 3, 0, 0, 0, 0, 0, 0},
	7:  {4, 3, 3, 1, 0, 0, 0, 0, 0},
	8:  {4, 3, 3, 2, 0, 0, 0, 0, 0},
	9:  {4, 3, 3, 3, 1, 0, 0, 0, 0},
	10: {4, 3, 3, 3, 2, 0, 0, 0, 0},
	11: {4, 3, 3, 3, 2, 1, 0, 0, 0},
	12: {4, 3, 3, 3, 2, 1, 0, 0, 0},
	13: {4, 3, 3, 3, 2, 1, 1, 0, 0},
	14: {4, 3, 3, 3, 2, 1, 1, 0, 0},
	15: {4, 3, 3, 3, 2, 1, 1, 1, 0},
	16: {4, 3, 3, 3, 2, 1, 1, 1, 0},
	17: {4, 3, 3, 3, 2, 1, 1, 1, 1},
	18: {4, 3, 3, 3, 3, 1, 1, 1, 1},
	19: {4, 3, 3, 3, 3, 2, 1, 1, 1},
	20: {4, 3, 3, 3, 3, 2, 2, 1, 1},
}

// Compute validates the request against the rules content and derives the sheet.
func Compute(rules *content.Content, req Request) (Response, error) {
	if err := validate(rules, req); err != nil {
		return Response{}, err
	}

	level := 0
	for _, c := range req.Classes {
		level += c.Level
	}

	conMod := modFor(req.AbilityScores["CON"])
	dexMod := modFor(req.AbilityScores["DEX"])

	primary := classByID(rules.Classes, req.Classes[0].ID)
	hpMax := hitDieValue(primary.HitDie) + conMod
	if level > 1 {
		hpMax += (avgHitDie(primary.HitDie) + conMod) * (level - 1)
	}

	spellSlots := make([]int, 9)
	if casterLevel := casterLevelOf(rules, req.Classes); casterLevel > 0 {
		copy(spellSlots, spellSlotsByCasterLevel[min(casterLevel, 20)])
	}

	abilities := map[string]Ability{}
	for _, a := range abilityOrder {
		abilities[a] = Ability{Score: req.AbilityScores[a], Mod: modFor(req.AbilityScores[a])}
	}

	features, pending := collectFeatures(rules, req, level)

	return Response{Sheet: Sheet{
		Level:            level,
		HP:               HP{Max: hpMax, Current: hpMax},
		AC:               10 + dexMod,
		ProficiencyBonus: proficiencyBonus(level),
		SpellSlots:       spellSlots,
		Abilities:        abilities,
		Features:         features,
		PendingChoices:   pending,
	}}, nil
}

func validate(rules *content.Content, req Request) error {
	if len(req.Classes) == 0 {
		return invalidDraftf("classes: at least one class is required")
	}
	if len(req.Classes) > MaxClasses {
		return invalidDraftf("classes: at most %d allowed", MaxClasses)
	}

	classIDs := map[string]bool{}
	casterLevel := 0
	for _, ci := range req.Classes {
		classIDs[ci.ID] = true
		cls := classByID(rules.Classes, ci.ID)
		if cls == nil {
			return buildErrorf("unknown class %q", ci.ID)
		}
		if ci.Level < 1 || ci.Level > 20 {
			return invalidDraftf("class %q level %d out of range 1-20", ci.ID, ci.Level)
		}
		if cls.Spellcaster {
			casterLevel += ci.Level
		}
		if ci.SubclassID != "" {
			sub := subclassByID(rules.Subclasses, ci.SubclassID)
			if sub == nil {
				return buildErrorf("unknown subclass %q", ci.SubclassID)
			}
			if sub.ClassID != ci.ID {
				return buildErrorf("subclass %q does not belong to class %q", ci.SubclassID, ci.ID)
			}
			if sub.LevelRequired > ci.Level {
				return buildErrorf("subclass %q requires level %d (class at level %d)",
					ci.SubclassID, sub.LevelRequired, ci.Level)
			}
		}
	}
	totalLevel := 0
	for _, c := range req.Classes {
		totalLevel += c.Level
	}
	if totalLevel > 20 {
		return invalidDraftf("total level %d exceeds 20", totalLevel)
	}

	if req.SpeciesID != "" && speciesByID(rules.Species, req.SpeciesID) == nil {
		return buildErrorf("unknown species %q", req.SpeciesID)
	}
	if req.BackgroundID != "" && backgroundByID(rules.Backgrounds, req.BackgroundID) == nil {
		return buildErrorf("unknown background %q", req.BackgroundID)
	}

	for _, a := range abilityOrder {
		score, ok := req.AbilityScores[a]
		if !ok {
			return invalidDraftf("abilityScores: missing %s", a)
		}
		if score < 3 || score > 20 {
			return invalidDraftf("abilityScores: %s=%d out of range 3-20", a, score)
		}
	}

	if len(req.Spells) > MaxSpells {
		return invalidDraftf("spells: at most %d allowed", MaxSpells)
	}
	spellIDs := map[string]bool{}
	for _, id := range req.Spells {
		if id == "" {
			return invalidDraftf("spells: empty id")
		}
		if spellIDs[id] {
			return invalidDraftf("spells: duplicate %q", id)
		}
		spellIDs[id] = true
		sp := spellByID(rules.Spells, id)
		if sp == nil {
			return buildErrorf("unknown spell %q", id)
		}
		if casterLevel == 0 {
			return invalidSpellf("spell %q requires a spellcaster class", id)
		}
		if !spellForClasses(sp, classIDs) {
			return invalidSpellf("spell %q is not on the class list", id)
		}
		if sp.Level > maxSlotLevelFor(casterLevel) {
			return invalidSpellf("spell %q is level %d, not yet available", id, sp.Level)
		}
	}

	if len(req.Feats) > MaxFeats {
		return invalidDraftf("feats: at most %d allowed", MaxFeats)
	}
	featIDs := map[string]bool{}
	for _, id := range req.Feats {
		if id == "" {
			return invalidDraftf("feats: empty id")
		}
		if featIDs[id] {
			return invalidDraftf("feats: duplicate %q", id)
		}
		featIDs[id] = true
		ft := featByID(rules.Feats, id)
		if ft == nil {
			return buildErrorf("unknown feat %q", id)
		}
		if requiresSpellcasting(ft) && casterLevel == 0 {
			return buildErrorf("feat %q requires spellcasting", id)
		}
	}

	return nil
}

func collectFeatures(rules *content.Content, req Request, level int) ([]SheetFeature, []PendingChoice) {
	classIDs := map[string]bool{}
	subclassIDs := map[string]bool{}
	for _, ci := range req.Classes {
		classIDs[ci.ID] = true
		if ci.SubclassID != "" {
			subclassIDs[ci.SubclassID] = true
		}
	}

	feats := make([]SheetFeature, 0, len(rules.Features))
	pending := make([]PendingChoice, 0)
	for _, f := range rules.Features {
		if f.Level > level {
			continue
		}
		if f.ClassID != "" && !classIDs[f.ClassID] {
			continue
		}
		if f.SubclassID != "" && !subclassIDs[f.SubclassID] {
			continue
		}
		feats = append(feats, SheetFeature{
			Name:        f.Name,
			Level:       f.Level,
			Description: descriptionOf(f.Data),
		})
		if opts, ok := stringOptions(f.Data["options"]); ok && len(opts) > 0 {
			n := intOption(f.Data["choose"], 1)
			pending = append(pending, PendingChoice{
				Type:        f.ID,
				Description: choiceDescription(f.Name, n, opts),
			})
		}
	}
	return feats, pending
}

func descriptionOf(data map[string]any) string {
	if s, ok := data["description"].(string); ok {
		return s
	}
	return ""
}

func stringOptions(v any) ([]string, bool) {
	switch t := v.(type) {
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out, true
	case []string:
		return t, true
	default:
		return nil, false
	}
}

func intOption(v any, def int) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case int64:
		return int(t)
	default:
		return def
	}
}

func choiceDescription(featureName string, n int, opts []string) string {
	noun := "options"
	if n == 1 {
		noun = "option"
	}
	sort.Strings(opts)
	return fmt.Sprintf("Choose %d %s %s: %s", n, featureName, noun, strings.Join(opts, ", "))
}

func casterLevelOf(rules *content.Content, classes []ClassInput) int {
	total := 0
	for _, ci := range classes {
		if cls := classByID(rules.Classes, ci.ID); cls != nil && cls.Spellcaster {
			total += ci.Level
		}
	}
	return total
}

func maxSlotLevelFor(casterLevel int) int {
	slots := spellSlotsByCasterLevel[min(casterLevel, 20)]
	for i := 8; i >= 0; i-- {
		if slots[i] > 0 {
			return i + 1
		}
	}
	return 0
}

func spellForClasses(sp *content.Spell, classIDs map[string]bool) bool {
	classes, ok := stringOptions(sp.Data["classes"])
	if !ok || len(classes) == 0 {
		return true
	}
	for _, c := range classes {
		if classIDs[c] {
			return true
		}
	}
	return false
}

func requiresSpellcasting(f *content.Feat) bool {
	if v, ok := f.Prerequisites["spellcasting"]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// proficiencyBonus follows the 5e progression: 2 + (level-1)/4.
func proficiencyBonus(level int) int { return (level + 7) / 4 }

func modFor(score int) int {
	if score >= 10 {
		return (score - 10) / 2
	}
	return (score - 11) / 2 // floor for scores below 10
}

// hitDieValue parses "d8" -> 8.
func hitDieValue(die string) int {
	var d int
	_, _ = fmt.Sscanf(die, "d%d", &d)
	return d
}

func avgHitDie(die string) int { return hitDieValue(die)/2 + 1 }

func classByID(classes []content.Class, id string) *content.Class {
	for i := range classes {
		if classes[i].ID == id {
			return &classes[i]
		}
	}
	return nil
}

func subclassByID(subclasses []content.Subclass, id string) *content.Subclass {
	for i := range subclasses {
		if subclasses[i].ID == id {
			return &subclasses[i]
		}
	}
	return nil
}

func speciesByID(species []content.Species, id string) *content.Species {
	for i := range species {
		if species[i].ID == id {
			return &species[i]
		}
	}
	return nil
}

func backgroundByID(bg []content.Background, id string) *content.Background {
	for i := range bg {
		if bg[i].ID == id {
			return &bg[i]
		}
	}
	return nil
}

func spellByID(spells []content.Spell, id string) *content.Spell {
	for i := range spells {
		if spells[i].ID == id {
			return &spells[i]
		}
	}
	return nil
}

func featByID(feats []content.Feat, id string) *content.Feat {
	for i := range feats {
		if feats[i].ID == id {
			return &feats[i]
		}
	}
	return nil
}
