package parser

import (
	"fmt"
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/normalizer"
)

// ParseClasses parses all class files from the 5etools class directory.
func ParseClasses(root string, filter filter.Filter) ([]content.Class, error) {
	files, err := ParseFilesInDir(root, "class", "class-")
	if err != nil {
		return nil, err
	}

	var classes []content.Class
	seen := make(map[string]bool)

	for _, file := range files {
		classArr, ok := UnwrapArray(file, "class")
		if !ok {
			continue
		}

		for _, item := range classArr {
			classMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			name := GetString(classMap, "name")
			source := GetString(classMap, "source")
			if !filter(name, source) {
				continue
			}

			id := normalizer.MakeClassID(name, source)
			if seen[id] {
				continue
			}
			seen[id] = true

			faces := GetInt(classMap, "hd.faces")
			hitDie := fmt.Sprintf("d%d", faces)
			spellcaster := GetString(classMap, "spellcastingAbility") != ""
			subclassLevel := 3 // default, may be overridden by subclass features

			// Extract primary ability
			primaryAbility := strings.ToUpper(GetString(classMap, "spellcastingAbility"))

			// Extract saving throws
			savingThrows := GetStringSlice(classMap, "proficiency")
			for i, s := range savingThrows {
				savingThrows[i] = strings.ToUpper(s)
			}

			// Extract skill choices
			var skillChoices []string
			if sp := GetStringSlice(classMap, "startingProficiencies.skills"); len(sp) > 0 {
				skillChoices = sp
			}

			classes = append(classes, content.Class{
				ID:                   id,
				Name:                 name,
				HitDie:               hitDie,
				Spellcaster:          spellcaster,
				SubclassLevel:        subclassLevel,
				PrimaryAbility:       primaryAbility,
				SuggestedSpecies:     []string{}, // 5etools doesn't provide this directly
				SuggestedBackground:  []string{},
				Data: map[string]any{
					"description":       "", // filled from fluff later
					"savingThrows":      savingThrows,
					"skillChoices":      skillChoices,
					"source":            source,
				},
			})
		}
	}

	return classes, nil
}