package etl5e

import (
	"github.com/hadnu/onatar/internal/content"
	"github.com/hadnu/onatar/internal/etl5e/filter"
	"github.com/hadnu/onatar/internal/etl5e/parser"
)

// LoadFrom5eTools loads and normalizes all 5etools data into Onatar's content.Content.
// The filter determines which sources to include.
func LoadFrom5eTools(root string, filter filter.Filter) (*content.Content, error) {
	c := &content.Content{}
	var err error

	// Order matters: classes first (features depend on them), then subclasses, features, etc.
	if c.Classes, err = parser.ParseClasses(root, filter); err != nil {
		return nil, err
	}
	if c.Subclasses, err = parser.ParseSubclasses(root, filter); err != nil {
		return nil, err
	}
	if c.Features, err = parser.ParseFeatures(root, filter); err != nil {
		return nil, err
	}
	if c.Species, err = parser.ParseSpecies(root, filter); err != nil {
		return nil, err
	}
	if c.Backgrounds, err = parser.ParseBackgrounds(root, filter); err != nil {
		return nil, err
	}
	if c.Spells, err = parser.ParseSpells(root, filter); err != nil {
		return nil, err
	}
	if c.Feats, err = parser.ParseFeats(root, filter); err != nil {
		return nil, err
	}
	if c.Items, err = parser.ParseItems(root, filter); err != nil {
		return nil, err
	}
	return c, nil
}