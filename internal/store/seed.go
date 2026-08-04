package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hadnu/onatar/internal/content"
)

// Seed upserts all content entities into MariaDB (idempotent). The source of
// truth is data/**/*.yaml; the data JSON column is never hand-edited (PRD §8.1).
func Seed(ctx context.Context, db *sql.DB, c *content.Content) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := seedClasses(ctx, tx, c.Classes); err != nil {
		return err
	}
	if err := seedSubclasses(ctx, tx, c.Subclasses); err != nil {
		return err
	}
	if err := seedSpecies(ctx, tx, c.Species); err != nil {
		return err
	}
	if err := seedBackgrounds(ctx, tx, c.Backgrounds); err != nil {
		return err
	}
	if err := seedSpells(ctx, tx, c.Spells); err != nil {
		return err
	}
	if err := seedFeats(ctx, tx, c.Feats); err != nil {
		return err
	}
	if err := seedFeatures(ctx, tx, c.Features); err != nil {
		return err
	}
	if err := seedItems(ctx, tx, c.Items); err != nil {
		return err
	}
	return tx.Commit()
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func seedClasses(ctx context.Context, tx execer, items []content.Class) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO classes (id, name, hit_die, spellcaster, subclass_level, data)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), hit_die=VALUES(hit_die),
				spellcaster=VALUES(spellcaster), subclass_level=VALUES(subclass_level),
				data=VALUES(data)`,
			it.ID, it.Name, it.HitDie, it.Spellcaster, it.SubclassLevel, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedSubclasses(ctx context.Context, tx execer, items []content.Subclass) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO subclasses (id, class_id, name, level_required, data)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE class_id=VALUES(class_id), name=VALUES(name),
				level_required=VALUES(level_required), data=VALUES(data)`,
			it.ID, it.ClassID, it.Name, it.LevelRequired, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedSpecies(ctx context.Context, tx execer, items []content.Species) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO species (id, name, data)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), data=VALUES(data)`,
			it.ID, it.Name, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedBackgrounds(ctx context.Context, tx execer, items []content.Background) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO backgrounds (id, name, data)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), data=VALUES(data)`,
			it.ID, it.Name, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedSpells(ctx context.Context, tx execer, items []content.Spell) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO spells (id, name, level, school, data)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), level=VALUES(level),
				school=VALUES(school), data=VALUES(data)`,
			it.ID, it.Name, it.Level, it.School, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedFeats(ctx context.Context, tx execer, items []content.Feat) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO feats (id, name, prerequisites, data)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name),
				prerequisites=VALUES(prerequisites), data=VALUES(data)`,
			it.ID, it.Name, mustJSON(it.Prerequisites), data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedFeatures(ctx context.Context, tx execer, items []content.Feature) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO features (id, class_id, subclass_id, name, level, data)
			VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE class_id=VALUES(class_id),
				subclass_id=VALUES(subclass_id), name=VALUES(name),
				level=VALUES(level), data=VALUES(data)`,
			it.ID, nullIfEmpty(it.ClassID), nullIfEmpty(it.SubclassID), it.Name, it.Level, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func seedItems(ctx context.Context, tx execer, items []content.Item) error {
	for _, it := range items {
		data, err := it.JSONData()
		if err != nil {
			return err
		}
		fmt.Printf("SEED ITEM: id=%s name=%s\n", it.ID, it.Name)
		_, err = tx.ExecContext(ctx, `
			INSERT INTO items (id, name, type, rarity, source, edition, data)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), type=VALUES(type),
				rarity=VALUES(rarity), source=VALUES(source), edition=VALUES(edition), data=VALUES(data)`,
			it.ID, it.Name, it.Type, it.Rarity, it.Source, it.Edition, data)
		if err != nil {
			fmt.Printf("SEED ITEM ERROR: id=%s err=%v\n", it.ID, err)
			return err
		}
	}
	return nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func mustJSON(v map[string]any) any {
	if len(v) == 0 {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return data
}
