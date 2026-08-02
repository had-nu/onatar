package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Campaign represents a campaign in the database.
type Campaign struct {
	ID          string         `json:"id"`
	OwnerID     string         `json:"owner_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description,omitempty"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
}

// CampaignMember represents a member of a campaign.
type CampaignMember struct {
	CampaignID string `json:"campaign_id"`
	UserID     string `json:"user_id"`
	Role       string `json:"role"` // "dm" or "player"
	JoinedAt   int64  `json:"joined_at"`
}

// CreateCampaign creates a new campaign.
func CreateCampaign(ctx context.Context, db *sql.DB, ownerID, name, description string) (*Campaign, error) {
	c := &Campaign{
		ID:          uuid.New().String(),
		OwnerID:     ownerID,
		Name:        name,
		Description: sql.NullString{String: description, Valid: description != ""},
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	query := `
		INSERT INTO campaigns (id, owner_id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.ExecContext(ctx, query, c.ID, c.OwnerID, c.Name, c.Description, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert campaign: %w", err)
	}

	// Add owner as DM member
	if err := AddCampaignMember(ctx, db, c.ID, ownerID, "dm"); err != nil {
		return nil, fmt.Errorf("add owner as member: %w", err)
	}

	return c, nil
}

// GetCampaign retrieves a campaign by ID.
func GetCampaign(ctx context.Context, db *sql.DB, id string) (*Campaign, error) {
	query := `
		SELECT id, owner_id, name, description, created_at, updated_at
		FROM campaigns WHERE id = ?
	`
	row := db.QueryRowContext(ctx, query, id)

	var c Campaign
	var desc sql.NullString
	err := row.Scan(&c.ID, &c.OwnerID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan campaign: %w", err)
	}
	c.Description = desc
	return &c, nil
}

// ListCampaignsByUser retrieves all campaigns where the user is a member.
func ListCampaignsByUser(ctx context.Context, db *sql.DB, userID string) ([]*Campaign, error) {
	query := `
		SELECT c.id, c.owner_id, c.name, c.description, c.created_at, c.updated_at
		FROM campaigns c
		JOIN campaign_members m ON c.id = m.campaign_id
		WHERE m.user_id = ?
		ORDER BY c.updated_at DESC
	`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var campaigns []*Campaign
	for rows.Next() {
		var c Campaign
		var desc sql.NullString
		if err := rows.Scan(&c.ID, &c.OwnerID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Description = desc
		campaigns = append(campaigns, &c)
	}
	return campaigns, rows.Err()
}

// GetCampaignWithMembers retrieves a campaign with its members.
func GetCampaignWithMembers(ctx context.Context, db *sql.DB, id string) (*Campaign, []CampaignMember, error) {
	campaign, err := GetCampaign(ctx, db, id)
	if err != nil {
		return nil, nil, err
	}
	if campaign == nil {
		return nil, nil, nil
	}

	members, err := GetCampaignMembers(ctx, db, id)
	if err != nil {
		return nil, nil, err
	}
	return campaign, members, nil
}

// AddCampaignMember adds a user to a campaign.
func AddCampaignMember(ctx context.Context, db *sql.DB, campaignID, userID, role string) error {
	query := `
		INSERT INTO campaign_members (campaign_id, user_id, role, joined_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE role = VALUES(role), joined_at = VALUES(joined_at)
	`
	_, err := db.ExecContext(ctx, query, campaignID, userID, role, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("add campaign member: %w", err)
	}
	return nil
}

// RemoveCampaignMember removes a user from a campaign.
func RemoveCampaignMember(ctx context.Context, db *sql.DB, campaignID, userID string) error {
	result, err := db.ExecContext(ctx, `DELETE FROM campaign_members WHERE campaign_id = ? AND user_id = ?`, campaignID, userID)
	if err != nil {
		return fmt.Errorf("remove campaign member: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

// GetCampaignMembers retrieves all members of a campaign.
func GetCampaignMembers(ctx context.Context, db *sql.DB, campaignID string) ([]CampaignMember, error) {
	query := `
		SELECT campaign_id, user_id, role, joined_at
		FROM campaign_members WHERE campaign_id = ?
	`
	rows, err := db.QueryContext(ctx, query, campaignID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var members []CampaignMember
	for rows.Next() {
		var m CampaignMember
		if err := rows.Scan(&m.CampaignID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// IsUserDM checks if a user is the DM of a campaign.
func IsUserDM(ctx context.Context, db *sql.DB, campaignID, userID string) (bool, error) {
	query := `SELECT role FROM campaign_members WHERE campaign_id = ? AND user_id = ?`
	var role string
	err := db.QueryRowContext(ctx, query, campaignID, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return role == "dm", nil
}

// GetCampaignCharacters retrieves all characters in a campaign.
func GetCampaignCharacters(ctx context.Context, db *sql.DB, campaignID string) ([]*Character, error) {
	query := `
		SELECT id, user_id, campaign_id, name, is_npc, draft, sheet, live, created_at, updated_at
		FROM characters WHERE campaign_id = ?
	`
	rows, err := db.QueryContext(ctx, query, campaignID)
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

// IsCampaignOwner checks if a user is the owner of a campaign.
func IsCampaignOwner(ctx context.Context, db *sql.DB, campaignID, userID string) (bool, error) {
	query := `SELECT owner_id FROM campaigns WHERE id = ?`
	var ownerID string
	err := db.QueryRowContext(ctx, query, campaignID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil
}