package parser

import (
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/normalizer"
)

func ParseSpecies(root string, filter filter.Filter) ([]content.Species, error) {
	file, err := ParseFile(root, "races.json")
	if err != nil {
		return nil, err
	}

	raceArr, ok := UnwrapArray(file, "race")
	if !ok {
		return nil, nil
	}

	var species []content.Species
	for _, item := range raceArr {
		raceMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name := GetString(raceMap, "name")
		source := GetString(raceMap, "source")
		if !filter(name, source) {
			continue
		}

		id := normalizer.MakeSpeciesID(name, source)

		abilityScores := make(map[string]int)
		if abilityArr, ok := UnwrapArray(raceMap, "ability"); ok {
			for _, item := range abilityArr {
				if m, ok := item.(map[string]any); ok {
					for k, v := range m {
						if vFloat, ok := v.(float64); ok {
							abilityScores[strings.ToUpper(k)] = int(vFloat)
						}
					}
				}
			}
		}

		speed := make(map[string]any)
		if speedMap, ok := raceMap["speed"].(map[string]any); ok {
			for k, v := range speedMap {
				speed[k] = v
			}
		}

		traits := []map[string]any{}
		if entriesArr, ok := UnwrapArray(raceMap, "entries"); ok {
			for _, item := range entriesArr {
				if m, ok := item.(map[string]any); ok {
					if name := GetString(m, "name"); name != "" {
						entries := GetStringSlice(m, "entries")
						traits = append(traits, map[string]any{
							"name":        name,
							"description": strings.Join(entries, "\n\n"),
						})
					}
				}
			}
		}

		species = append(species, content.Species{
			ID:            id,
			Name:          name,
			AbilityScores: abilityScores,
			Data: map[string]any{
				"description": "",
				"traits":      traits,
				"size":        GetStringSlice(raceMap, "size"),
				"speed":       speed,
				"source":      source,
				"lineage":     GetString(raceMap, "lineage"),
			},
		})
	}

	return species, nil
}