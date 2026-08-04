package parser

import (
	"fmt"
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
)

func ParseItems(root string, filter filter.Filter) ([]content.Item, error) {
	file, err := ParseFile(root, "items.json")
	if err != nil {
		file, err = ParseFile(root, "items-base.json")
		if err != nil {
			return nil, err
		}
	}

	itemArr, ok := UnwrapArray(file, "item")
	if !ok {
		return nil, nil
	}

	var items []content.Item
	seen := make(map[string]bool)

	for _, item := range itemArr {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name := GetString(itemMap, "name")
		source := GetString(itemMap, "source")
		if !filter(name, source) {
			continue
		}

id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	originalLen := len(id)
	if len(id) > 64 {
		id = id[:64]
		fmt.Printf("DEBUG: item=%s original_len=%d truncated_id=%s len=%d\n", name, originalLen, id, len(id))
	}
	if len(id) > 64 {
		fmt.Printf("ERROR: Still too long after truncation! id=%s len=%d\n", id, len(id))
	}
	if len(id) > 64 {
		fmt.Printf("FATAL: ID still too long after truncation! item=%s id=%s len=%d\n", name, id, len(id))
	}
	if seen[id] {
		fmt.Printf("COLLISION: item=%s id=%s\n", name, id)
	}
	if seen[id] {
			continue
		}
		seen[id] = true

		itemType := "gear"
		if t := GetString(itemMap, "type"); t != "" {
			itemType = t
		}

		rarity := ""
		if r := GetString(itemMap, "rarity"); r != "" {
			rarity = r
		}

		weight := 0.0
		if w := GetString(itemMap, "weight"); w != "" {
			if f, ok := itemMap["weight"].(float64); ok {
				weight = f
			}
		}

		reqAttune := GetString(itemMap, "reqAttune")
		var reqAttuneTags []string
		if attuneTagsArr, ok := UnwrapArray(itemMap, "reqAttuneTags"); ok {
			for _, item := range attuneTagsArr {
				if m, ok := item.(map[string]any); ok {
					if class := GetString(m, "class"); class != "" {
						reqAttuneTags = append(reqAttuneTags, class)
					}
				}
			}
		}

		bonusSpellAttack := GetString(itemMap, "bonusSpellAttack")
		bonusSpellSaveDc := GetString(itemMap, "bonusSpellSaveDc")

		var focus []string
		if focusArr, ok := UnwrapArray(itemMap, "focus"); ok {
			for _, item := range focusArr {
				if s, ok := item.(string); ok {
					focus = append(focus, s)
				}
			}
		}

		wondrous := GetBool(itemMap, "wondrous")

		description := ""
		if entriesArr, ok := UnwrapArray(itemMap, "entries"); ok {
			var parts []string
			for _, e := range entriesArr {
				if s, ok := e.(string); ok {
					parts = append(parts, s)
				}
			}
			description = strings.Join(parts, "\n\n")
		}

		var properties []string
		if propArr, ok := UnwrapArray(itemMap, "property"); ok {
			for _, item := range propArr {
				if s, ok := item.(string); ok {
					properties = append(properties, s)
				}
			}
		}

		var damage map[string]any
		if dmg, ok := itemMap["dmg"].(map[string]any); ok {
			damage = dmg
		}

		var rangeData map[string]any
		if r, ok := itemMap["range"].(map[string]any); ok {
			rangeData = r
		}

		ac := 0
		if acVal := itemMap["ac"]; acVal != nil {
			if f, ok := acVal.(float64); ok {
				ac = int(f)
			}
		}

		stealthDisadv := false
		if sd := itemMap["stealth"]; sd != nil {
			if b, ok := sd.(bool); ok {
				stealthDisadv = b
			}
		}

		baseItem := GetString(itemMap, "baseItem")
		category := GetString(itemMap, "category")

		items = append(items, content.Item{
			ID:       id,
			Name:     name,
			Type:     itemType,
			Rarity:   rarity,
			Source:   source,
			Edition:  "2024",
			Data: map[string]any{
				"description":       description,
				"source":            source,
				"weight":            weight,
				"rarity":            rarity,
				"type":              itemType,
				"category":          category,
				"baseItem":          baseItem,
				"wondrous":          wondrous,
				"reqAttune":         reqAttune,
				"reqAttuneTags":     reqAttuneTags,
				"bonusSpellAttack":  bonusSpellAttack,
				"bonusSpellSaveDc":  bonusSpellSaveDc,
				"focus":             focus,
				"properties":        properties,
				"damage":            damage,
				"range":             rangeData,
				"ac":                ac,
				"stealthDisadv":     stealthDisadv,
				"entries":           itemMap["entries"],
			},
		})
	}

	return items, nil
}
