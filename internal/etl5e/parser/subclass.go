package parser

import (
	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/normalizer"
)

func ParseSubclasses(root string, filter filter.Filter) ([]content.Subclass, error) {
	files, err := ParseFilesInDir(root, "class", "class-")
	if err != nil {
		return nil, err
	}

	var subclasses []content.Subclass
	seen := make(map[string]bool)

	for _, file := range files {
		subArr, ok := UnwrapArray(file, "subclass")
		if !ok {
			continue
		}

		for _, item := range subArr {
			subMap, ok := item.(map[string]any)
			if !ok {
				continue
			}

			name := GetString(subMap, "name")
			source := GetString(subMap, "source")
			if !filter(name, source) {
				continue
			}

			id := normalizer.MakeSubclassID(name, source)
			if seen[id] {
				continue
			}
			seen[id] = true

			classID := normalizer.MakeClassID(GetString(subMap, "className"), GetString(subMap, "source"))

			subclasses = append(subclasses, content.Subclass{
				ID:            id,
				ClassID:       classID,
				Name:          name,
				LevelRequired: 1, // will be updated from features
				RecommendedSpells: map[int][]string{},
				Data: map[string]any{
					"description":        "",
					"source":             source,
					"subclassFeatures":   GetStringSlice(subMap, "subclassFeatures"),
				},
			})
		}
	}

	return subclasses, nil
}