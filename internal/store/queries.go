package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/hadnu/onatar/internal/auth"
)

// Queries wraps all store functions for use in HTTP handlers.
type Queries struct {
	db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{db: db}
}

// User operations
func (q *Queries) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	return auth.GetUserByID(ctx, q.db, id)
}

func (q *Queries) UpsertUser(ctx context.Context, user *auth.User) error {
	return auth.UpsertUser(ctx, q.db, user)
}

// Session operations
func (q *Queries) CreateSession(ctx context.Context, sessionID, userID string, expiresAt time.Time) error {
	return auth.CreateSession(ctx, q.db, sessionID, userID, expiresAt)
}

func (q *Queries) GetSession(ctx context.Context, sessionID string) (*auth.Session, error) {
	return auth.GetSession(ctx, q.db, sessionID)
}

func (q *Queries) DeleteSession(ctx context.Context, sessionID string) error {
	return auth.DeleteSession(ctx, q.db, sessionID)
}

// Character operations
func (q *Queries) CreateCharacter(ctx context.Context, c *Character) error {
	return CreateCharacter(ctx, q.db, c)
}

func (q *Queries) GetCharacter(ctx context.Context, id string) (*Character, error) {
	return GetCharacter(ctx, q.db, id)
}

func (q *Queries) ListCharactersByUser(ctx context.Context, userID string) ([]*Character, error) {
	return ListCharactersByUser(ctx, q.db, userID)
}

func (q *Queries) UpdateCharacter(ctx context.Context, c *Character) error {
	return UpdateCharacter(ctx, q.db, c)
}

func (q *Queries) DeleteCharacter(ctx context.Context, userID, id string) error {
	return DeleteCharacter(ctx, q.db, userID, id)
}

// Campaign operations
func (q *Queries) CreateCampaign(ctx context.Context, ownerID, name, description string) (*Campaign, error) {
	return CreateCampaign(ctx, q.db, ownerID, name, description)
}

func (q *Queries) GetCampaign(ctx context.Context, id string) (*Campaign, error) {
	return GetCampaign(ctx, q.db, id)
}

func (q *Queries) ListCampaignsByUser(ctx context.Context, userID string) ([]*Campaign, error) {
	return ListCampaignsByUser(ctx, q.db, userID)
}

func (q *Queries) GetCampaignWithMembers(ctx context.Context, id string) (*Campaign, []CampaignMember, error) {
	return GetCampaignWithMembers(ctx, q.db, id)
}

func (q *Queries) GetCampaignMembers(ctx context.Context, campaignID string) ([]CampaignMember, error) {
	return GetCampaignMembers(ctx, q.db, campaignID)
}

func (q *Queries) AddCampaignMember(ctx context.Context, campaignID, userID, role string) error {
	return AddCampaignMember(ctx, q.db, campaignID, userID, role)
}

func (q *Queries) RemoveCampaignMember(ctx context.Context, campaignID, userID string) error {
	return RemoveCampaignMember(ctx, q.db, campaignID, userID)
}

func (q *Queries) GetCampaignCharacters(ctx context.Context, campaignID string) ([]*Character, error) {
	return GetCampaignCharacters(ctx, q.db, campaignID)
}

func (q *Queries) IsUserDM(ctx context.Context, campaignID, userID string) (bool, error) {
	return IsUserDM(ctx, q.db, campaignID, userID)
}

func (q *Queries) IsCampaignOwner(ctx context.Context, campaignID, userID string) (bool, error) {
	return IsCampaignOwner(ctx, q.db, campaignID, userID)
}