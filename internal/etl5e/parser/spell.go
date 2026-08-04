package parser

import (
	"path/filepath"
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
)

var schoolMap = map[string]string{
	"A": "abjuration",
	"C": "conjuration",
	"D": "divination",
	"E": "enchantment",
	"V": "evocation",
	"I": "illusion",
	"N": "necromancy",
	"T": "transmutation",
}

func ParseSpells(root string, filter filter.Filter) ([]content.Spell, error) {
	spellFiles, err := filepath.Glob(filepath.Join(root, "spells", "spells-*.json"))
	if err != nil {
		return nil, err
	}

	var spells []content.Spell
	seen := make(map[string]bool)

	for _, file := range spellFiles {
		data, err := ParseFileFullPath(file)
		if err != nil {
			continue
		}

		spellArr, ok := UnwrapArray(data, "spell")
		if !ok {
			continue
		}

		for _, item := range spellArr {
			spellMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			name := GetString(spellMap, "name")
			source := GetString(spellMap, "source")
			if !filter(name, source) {
				continue
			}

			id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
			if seen[id] {
				continue
			}
			seen[id] = true

			level := GetInt(spellMap, "level")
			schoolCode := GetString(spellMap, "school")
			school := schoolMap[schoolCode]
			if school == "" {
				school = "unknown"
			}

			// Parse classes that can cast this spell
			var classes []string
			if classArr, ok := UnwrapArray(spellMap, "classes"); ok {
				for _, item := range classArr {
					if m, ok := item.(map[string]any); ok {
						if fromClassList := GetStringSlice(m, "fromClassList"); len(fromClassList) > 0 {
							classes = append(classes, fromClassList...)
						}
					}
				}
			}

			// Parse components
			components := make(map[string]bool)
			if compMap, ok := spellMap["components"].(map[string]any); ok {
				for k, v := range compMap {
					if vBool, ok := v.(bool); ok && vBool {
						components[k] = true
					}
				}
			}

			spells = append(spells, content.Spell{
				ID:     strings.ToLower(strings.ReplaceAll(name, " ", "-")),
				Name:   name,
				Level:  level,
				School: school,
				Data: map[string]any{
					"description":       "",
					"classes":           classes,
					"components":        components,
					"source":            source,
				},
			})
		}
	}

	return spells, nil
}