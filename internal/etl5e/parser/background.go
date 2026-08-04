package parser

import (
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/normalizer"
)

func ParseBackgrounds(root string, filter filter.Filter) ([]content.Background, error) {
	file, err := ParseFile(root, "backgrounds.json")
	if err != nil {
		return nil, err
	}

	bgArr, ok := UnwrapArray(file, "background")
	if !ok {
		return nil, nil
	}

	var backgrounds []content.Background
	for _, item := range bgArr {
		bgMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name := GetString(bgMap, "name")
		source := GetString(bgMap, "source")
		if !filter(name, source) {
			continue
		}

		id := normalizer.MakeBackgroundID(name, source)

		// Parse skill proficiencies
		var skills []string
		if skillArr, ok := UnwrapArray(bgMap, "skillProficiencies"); ok {
			for _, item := range skillArr {
				if m, ok := item.(map[string]any); ok {
					for k, v := range m {
						if vBool, ok := v.(bool); ok && vBool {
							skills = append(skills, k)
						}
					}
				}
			}
		}

		// Parse tool proficiencies
		var tools []string
		if toolArr, ok := UnwrapArray(bgMap, "toolProficiencies"); ok {
			for _, item := range toolArr {
				if m, ok := item.(map[string]any); ok {
					for k, v := range m {
						if vBool, ok := v.(bool); ok && vBool {
							tools = append(tools, k)
						}
					}
				}
			}
		}

		// Parse languages
		languages := 0
		if langArr, ok := UnwrapArray(bgMap, "languageProficiencies"); ok {
			for _, item := range langArr {
				if m, ok := item.(map[string]any); ok {
					for _, v := range m {
						if vBool, ok := v.(bool); ok && vBool {
							languages++
						}
					}
				}
			}
		}

		// Parse features from entries
		var features []map[string]any
		if entriesArr, ok := UnwrapArray(bgMap, "entries"); ok {
			for _, item := range entriesArr {
				if m, ok := item.(map[string]any); ok {
					if name := GetString(m, "name"); name != "" {
						entries := GetStringSlice(m, "entries")
						features = append(features, map[string]any{
							"name":        name,
							"description": strings.Join(entries, "\n\n"),
						})
					}
				}
			}
		}

		// Parse feats
		var feats []string
		if featArr, ok := UnwrapArray(bgMap, "feats"); ok {
			for _, item := range featArr {
				if m, ok := item.(map[string]any); ok {
					for k := range m {
						feats = append(feats, k)
					}
				}
			}
		}

		backgrounds = append(backgrounds, content.Background{
			ID:   id,
			Name: name,
			Data: map[string]any{
				"description":            "",
				"skills":                 skills,
				"toolProficiencies":      tools,
				"languages":              languages,
				"features":               features,
				"feats":                  feats,
				"source":                 source,
			},
		})
	}

	return backgrounds, nil
}