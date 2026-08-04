package parser

import (
	"strings"

	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
)

func ParseFeats(root string, filter filter.Filter) ([]content.Feat, error) {
	file, err := ParseFile(root, "feats.json")
	if err != nil {
		return nil, err
	}

	featArr, ok := UnwrapArray(file, "feat")
	if !ok {
		return nil, nil
	}

	var feats []content.Feat
	for _, item := range featArr {
		featMap, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name := GetString(featMap, "name")
		source := GetString(featMap, "source")
		if !filter(name, source) {
			continue
		}

		id := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

		// Parse prerequisites
		prereqs := make(map[string]any)
		if prereqArr, ok := UnwrapArray(featMap, "prerequisite"); ok {
			for _, item := range prereqArr {
				if m, ok := item.(map[string]any); ok {
					for k, v := range m {
						prereqs[k] = v
					}
				}
			}
		}

		feats = append(feats, content.Feat{
			ID:            id,
			Name:          name,
			Prerequisites: prereqs,
			Data: map[string]any{
				"description": "",
				"source":      source,
			},
		})
	}

	return feats, nil
}