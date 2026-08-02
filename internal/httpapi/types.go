package httpapi

// BuildRequest mirrors the API contract for POST /build.
type BuildRequest struct {
	Name           string            `json:"name"`
	Classes        []ClassInput      `json:"classes"`
	SpeciesID      string            `json:"speciesId,omitempty"`
	BackgroundID   string            `json:"backgroundId,omitempty"`
	AbilityScores  map[string]int    `json:"abilityScores"`
	AbilityMethod  string            `json:"abilityMethod,omitempty"`
	Skills         []string          `json:"skills,omitempty"`
	Spells         []string          `json:"spells,omitempty"`
	Feats          []string          `json:"feats,omitempty"`
	Equipment      []string          `json:"equipment,omitempty"`
	IsNpc          bool              `json:"isNpc,omitempty"`
}

type ClassInput struct {
	ID         string `json:"id"`
	Level      int    `json:"level"`
	SubclassID string `json:"subclassId,omitempty"`
}

type AbilityScore struct {
	Score int `json:"score"`
	Mod   int `json:"mod"`
}

type SheetFeature struct {
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type PendingChoice struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Sheet struct {
	Level            int                    `json:"level"`
	HP               struct{ Max, Current int } `json:"hp"`
	AC               int                    `json:"ac"`
	ProficiencyBonus int                    `json:"proficiency_bonus"`
	SpellSlots       []int                  `json:"spell_slots"`
	Abilities        map[string]AbilityScore `json:"abilities"`
	Features         []SheetFeature         `json:"features"`
	PendingChoices   []PendingChoice        `json:"pending_choices"`
}

type BuildResponse struct {
	Sheet Sheet `json:"sheet"`
}

type SheetLive struct {
	HPCurrent  int                    `json:"hp_current"`
	SlotsUsed  []int                  `json:"slots_used"`
	Conditions []string               `json:"conditions"`
	Resources  map[string]int         `json:"resources"`
}

type CampaignResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	OwnerID     string   `json:"owner_id"`
	Members     []Member `json:"members,omitempty"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

type Member struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type CreateCampaignRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type CharacterResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	IsNpc      bool         `json:"is_npc"`
	CampaignID string       `json:"campaign_id,omitempty"`
	Draft      BuildRequest `json:"draft"`
	Sheet      *Sheet       `json:"sheet,omitempty"`
	Live       *SheetLive   `json:"live,omitempty"`
	CreatedAt  int64        `json:"created_at"`
	UpdatedAt  int64        `json:"updated_at"`
}