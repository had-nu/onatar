package parser

import (
	"path/filepath"
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/normalizer"
)

func ParseFeatures(root string, filter filter.Filter) ([]content.Feature, error) {
	// Parse class features
	featureFiles, err := filepath.Glob(filepath.Join(root, "class", "class-features-*.json"))
	if err != nil {
		return nil, err
	}

	var features []content.Feature
	seen := make(map[string]bool)

	for _, file := range featureFiles {
		data, err := ParseFileFullPath(file)
		if err != nil {
			continue
		}

		featArr, ok := UnwrapArray(data, "classFeature")
		if !ok {
			continue
		}

		for _, item := range featArr {
			featMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			name := GetString(featMap, "name")
			source := GetString(featMap, "source")
			className := GetString(featMap, "className")
			classSource := GetString(featMap, "classSource")
			level := GetInt(featMap, "level")

			if !filter(name, source) {
				continue
			}

			classID := normalizer.MakeClassID(className, classSource)
			id := normalizer.MakeFeatureID(classID, name)

			if seen[id] {
				continue
			}
			seen[id] = true

			// Parse description from entries
			desc := ""
			if entriesArr, ok := UnwrapArray(featMap, "entries"); ok {
				var parts []string
				for _, e := range entriesArr {
					if s, ok := e.(string); ok {
						parts = append(parts, s)
					} else if m, ok := e.(map[string]any); ok {
						if name := GetString(m, "name"); name != "" {
							entries := GetStringSlice(m, "entries")
							parts = append(parts, name+": "+strings.Join(entries, "\n\n"))
						}
					}
				}
				desc = strings.Join(parts, "\n\n")
			}

			features = append(features, content.Feature{
				ID:         id,
				ClassID:    classID,
				SubclassID: "",
				Name:       name,
				Level:      level,
				Data: map[string]any{
					"description": desc,
					"source":      source,
				},
			})
		}
	}

	// Parse subclass features
	subclassFeatureFiles, err := filepath.Glob(filepath.Join(root, "subclass", "subclass-features-*.json"))
	if err != nil {
		return nil, err
	}

	for _, file := range subclassFeatureFiles {
		data, err := ParseFileFullPath(file)
		if err != nil {
			continue
		}

		featArr, ok := UnwrapArray(data, "subclassFeature")
		if !ok {
			continue
		}

		for _, item := range featArr {
			featMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			name := GetString(featMap, "name")
			source := GetString(featMap, "source")
			className := GetString(featMap, "className")
			classSource := GetString(featMap, "classSource")
			subclassShortName := GetString(featMap, "subclassShortName")
			subclassSource := GetString(featMap, "subclassSource")
			level := GetInt(featMap, "level")

			if !filter(name, source) {
				continue
			}

			classID := normalizer.MakeClassID(className, classSource)
			subclassID := normalizer.MakeSubclassID(subclassShortName, subclassSource)
			id := normalizer.MakeFeatureID(classID+"-"+subclassShortName, name)

			if seen[id] {
				continue
			}
			seen[id] = true

			desc := ""
			if entriesArr, ok := UnwrapArray(featMap, "entries"); ok {
				var parts []string
				for _, e := range entriesArr {
					if s, ok := e.(string); ok {
						parts = append(parts, s)
					}
				}
				desc = strings.Join(parts, "\n\n")
			}

			features = append(features, content.Feature{
				ID:         id,
				ClassID:    classID,
				SubclassID: subclassID,
				Name:       name,
				Level:      level,
				Data: map[string]any{
					"description": desc,
					"source":      source,
				},
			})
		}
	}

	return features, nil
}