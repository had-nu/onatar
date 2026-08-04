package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hadnu/onatar/internal/content"
)

// Content loads all rule content from MariaDB. It is the runtime source of the
// GET /content API (PRD §3.5): YAML files are the source of truth and the seed
// writes them here; the server never parses YAML directly.
func Content(ctx context.Context, db *sql.DB) (*content.Content, error) {
	out := &content.Content{}
	var err error

	if out.Classes, err = queryClasses(ctx, db); err != nil {
		return nil, fmt.Errorf("classes: %w", err)
	}
	if out.Subclasses, err = querySubclasses(ctx, db); err != nil {
		return nil, fmt.Errorf("subclasses: %w", err)
	}
	if out.Species, err = querySpecies(ctx, db); err != nil {
		return nil, fmt.Errorf("species: %w", err)
	}
	if out.Backgrounds, err = queryBackgrounds(ctx, db); err != nil {
		return nil, fmt.Errorf("backgrounds: %w", err)
	}
	if out.Spells, err = querySpells(ctx, db); err != nil {
		return nil, fmt.Errorf("spells: %w", err)
	}
	if out.Feats, err = queryFeats(ctx, db); err != nil {
		return nil, fmt.Errorf("feats: %w", err)
	}
	if out.Features, err = queryFeatures(ctx, db); err != nil {
		return nil, fmt.Errorf("features: %w", err)
	}
	if out.Items, err = queryItems(ctx, db); err != nil {
		return nil, fmt.Errorf("items: %w", err)
	}
	return out, nil
}

func queryClasses(ctx context.Context, db *sql.DB) ([]content.Class, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, hit_die, spellcaster, subclass_level, data
		FROM classes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Class
	for rows.Next() {
		var (
			c              content.Class
			dataRaw        []byte
			hitDie         string
			spellcaster    bool
			subclassLevel  int
		)
		if err := rows.Scan(&c.ID, &c.Name, &hitDie, &spellcaster, &subclassLevel, &dataRaw); err != nil {
			return nil, err
		}
		c.HitDie = hitDie
		c.Spellcaster = spellcaster
		c.SubclassLevel = subclassLevel
		c.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &c.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", c.ID, err)
		}
		c.SuggestedSpecies = toStringSlice(c.Data["suggestedSpecies"])
		c.SuggestedBackground = toStringSlice(c.Data["suggestedBackgrounds"])
		delete(c.Data, "suggestedSpecies")
		delete(c.Data, "suggestedBackgrounds")
		out = append(out, c)
	}
	return out, rows.Err()
}

func querySubclasses(ctx context.Context, db *sql.DB) ([]content.Subclass, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, class_id, name, level_required, data
		FROM subclasses ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Subclass
	for rows.Next() {
		var (
			s       content.Subclass
			dataRaw []byte
		)
		if err := rows.Scan(&s.ID, &s.ClassID, &s.Name, &s.LevelRequired, &dataRaw); err != nil {
			return nil, err
		}
		s.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &s.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", s.ID, err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func querySpecies(ctx context.Context, db *sql.DB) ([]content.Species, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, data FROM species ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Species
	for rows.Next() {
		var (
			s       content.Species
			dataRaw []byte
		)
		if err := rows.Scan(&s.ID, &s.Name, &dataRaw); err != nil {
			return nil, err
		}
		s.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &s.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", s.ID, err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func queryBackgrounds(ctx context.Context, db *sql.DB) ([]content.Background, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, data FROM backgrounds ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Background
	for rows.Next() {
		var (
			b       content.Background
			dataRaw []byte
		)
		if err := rows.Scan(&b.ID, &b.Name, &dataRaw); err != nil {
			return nil, err
		}
		b.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &b.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", b.ID, err)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func querySpells(ctx context.Context, db *sql.DB) ([]content.Spell, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, level, school, data FROM spells ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Spell
	for rows.Next() {
		var (
			s       content.Spell
			dataRaw []byte
		)
		if err := rows.Scan(&s.ID, &s.Name, &s.Level, &s.School, &dataRaw); err != nil {
			return nil, err
		}
		s.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &s.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", s.ID, err)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func queryFeats(ctx context.Context, db *sql.DB) ([]content.Feat, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, prerequisites, data FROM feats ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Feat
	for rows.Next() {
		var (
			f       content.Feat
			prereq  sql.NullString
			dataRaw []byte
		)
		if err := rows.Scan(&f.ID, &f.Name, &prereq, &dataRaw); err != nil {
			return nil, err
		}
		f.Prerequisites = map[string]any{}
		if prereq.Valid && prereq.String != "" {
			if err := json.Unmarshal([]byte(prereq.String), &f.Prerequisites); err != nil {
				return nil, fmt.Errorf("%s: %w", f.ID, err)
			}
		}
		f.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &f.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", f.ID, err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func queryFeatures(ctx context.Context, db *sql.DB) ([]content.Feature, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, COALESCE(class_id,''), COALESCE(subclass_id,''), name, level, data
		FROM features ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Feature
	for rows.Next() {
		var (
			f       content.Feature
			dataRaw []byte
		)
		if err := rows.Scan(&f.ID, &f.ClassID, &f.SubclassID, &f.Name, &f.Level, &dataRaw); err != nil {
			return nil, err
		}
		f.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &f.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", f.ID, err)
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func queryItems(ctx context.Context, db *sql.DB) ([]content.Item, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, type, COALESCE(rarity,''), source, edition, data
		FROM items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var out []content.Item
	for rows.Next() {
		var (
			i       content.Item
			dataRaw []byte
		)
		if err := rows.Scan(&i.ID, &i.Name, &i.Type, &i.Rarity, &i.Source, &i.Edition, &dataRaw); err != nil {
			return nil, err
		}
		i.Data = map[string]any{}
		if err := json.Unmarshal(dataRaw, &i.Data); err != nil {
			return nil, fmt.Errorf("%s: %w", i.ID, err)
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

// toStringSlice converts a JSON-decoded value ([]any or []string) into []string.
func toStringSlice(v any) []string {
	switch t := v.(type) {
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	default:
		return nil
	}
}
