package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Character represents a character stored in the database.
type Character struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	CampaignID sql.NullString  `json:"campaign_id,omitempty"`
	Name       string          `json:"name"`
	IsNpc      bool            `json:"is_npc"`
	Draft      json.RawMessage `json:"draft"`
	Sheet      json.RawMessage `json:"sheet,omitempty"`
	Live       json.RawMessage `json:"live,omitempty"`
	CreatedAt  int64           `json:"created_at"`
	UpdatedAt  int64           `json:"updated_at"`
}

// CreateCharacter inserts a new character into the database.
func CreateCharacter(ctx context.Context, db *sql.DB, c *Character) error {
	draft, _ := json.Marshal(c.Draft)
	sheet, _ := json.Marshal(c.Sheet)
	live, _ := json.Marshal(c.Live)

	query := `
		INSERT INTO characters (id, user_id, campaign_id, name, is_npc, draft, sheet, live, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query,
		c.ID, c.UserID, c.CampaignID, c.Name, c.IsNpc,
		draft, sheet, live, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert character: %w", err)
	}
	return nil
}

// GetCharacter retrieves a character by ID.
func GetCharacter(ctx context.Context, db *sql.DB, id string) (*Character, error) {
	query := `
		SELECT id, user_id, campaign_id, name, is_npc, draft, sheet, live, created_at, updated_at
		FROM characters WHERE id = ?
	`
	row := db.QueryRowContext(ctx, query, id)

	var c Character
	err := row.Scan(
		&c.ID, &c.UserID, &c.CampaignID, &c.Name, &c.IsNpc,
		&c.Draft, &c.Sheet, &c.Live, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan character: %w", err)
	}
	return &c, nil
}

// ListCharactersByUser retrieves all characters for a user.
func ListCharactersByUser(ctx context.Context, db *sql.DB, userID string) ([]*Character, error) {
	query := `
		SELECT id, user_id, campaign_id, name, is_npc, draft, sheet, live, created_at, updated_at
		FROM characters WHERE user_id = ? ORDER BY updated_at DESC
	`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var chars []*Character
	for rows.Next() {
		var c Character
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.CampaignID, &c.Name, &c.IsNpc,
			&c.Draft, &c.Sheet, &c.Live, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		chars = append(chars, &c)
	}
	return chars, rows.Err()
}

// UpdateCharacter updates an existing character.
func UpdateCharacter(ctx context.Context, db *sql.DB, c *Character) error {
	draft, _ := json.Marshal(c.Draft)
	sheet, _ := json.Marshal(c.Sheet)
	live, _ := json.Marshal(c.Live)

	query := `
		UPDATE characters
		SET name = ?, is_npc = ?, draft = ?, sheet = ?, live = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`
	result, err := db.ExecContext(ctx, query,
		c.Name, c.IsNpc, draft, sheet, live, c.UpdatedAt,
		c.ID, c.UserID,
	)
	if err != nil {
		return fmt.Errorf("update character: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("character not found or not owned by user")
	}
	return nil
}

// DeleteCharacter deletes a character owned by the user.
func DeleteCharacter(ctx context.Context, db *sql.DB, userID, id string) error {
	result, err := db.ExecContext(ctx, `DELETE FROM characters WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return fmt.Errorf("delete character: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("character not found or not owned by user")
	}
	return nil
}

// CreateCharacterWithDraft creates a new character with generated ID and timestamps.
func CreateCharacterWithDraft(ctx context.Context, db *sql.DB, userID string, draft json.RawMessage, isNpc bool, name string) (*Character, error) {
	c := &Character{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		IsNpc:     isNpc,
		Draft:     draft,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	draftBytes, _ := json.Marshal(draft)
	query := `
		INSERT INTO characters (id, user_id, name, is_npc, draft, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query, c.ID, c.UserID, c.Name, c.IsNpc, draftBytes, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert character: %w", err)
	}
	return c, nil
}