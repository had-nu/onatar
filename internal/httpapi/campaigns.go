package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hadnu/onatar/internal/store"
)

// handleListCampaigns returns all campaigns for the authenticated user.
func (s *Server) handleListCampaigns(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	campaigns, err := store.ListCampaignsByUser(ctx, s.db, user.ID)
	if err != nil {
		s.logger.Error("list campaigns", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list campaigns", nil)
		return
	}

	var resp []CampaignResponse
	for _, c := range campaigns {
		resp = append(resp, CampaignResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description.String,
			OwnerID:     c.OwnerID,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleCreateCampaign creates a new campaign for the authenticated user (DM).
func (s *Server) handleCreateCampaign(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	var req CreateCampaignRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "INVALID_NAME", "campaign name is required", nil)
		return
	}

	campaign, err := store.CreateCampaign(r.Context(), s.db, user.ID, req.Name, req.Description)
	if err != nil {
		s.logger.Error("create campaign", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create campaign", nil)
		return
	}

	resp := CampaignResponse{
		ID:        campaign.ID,
		Name:      campaign.Name,
		Description: campaign.Description.String,
		OwnerID:   campaign.OwnerID,
		CreatedAt: campaign.CreatedAt,
		UpdatedAt: campaign.UpdatedAt,
	}
	writeJSON(w, http.StatusCreated, resp)
}

// handleGetCampaign returns a campaign by ID.
func (s *Server) handleGetCampaign(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "campaign ID is required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	campaign, members, err := store.GetCampaignWithMembers(ctx, s.db, id)
	if err != nil {
		s.logger.Error("get campaign", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get campaign", nil)
		return
	}

	if campaign == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "campaign not found", nil)
		return
	}

	// Check if user is a member
	isMember := false
	for _, m := range members {
		if m.UserID == user.ID {
			isMember = true
			break
		}
	}
	if !isMember {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "not a member of this campaign", nil)
		return
	}

	var memberList []Member
	for _, m := range members {
		memberList = append(memberList, Member{UserID: m.UserID, Role: m.Role})
	}

	resp := CampaignResponse{
		ID:          campaign.ID,
		Name:        campaign.Name,
		Description: campaign.Description.String,
		OwnerID:     campaign.OwnerID,
		Members:     memberList,
		CreatedAt:   campaign.CreatedAt,
		UpdatedAt:   campaign.UpdatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleListCampaignCharacters returns all characters in a campaign (DM only).
func (s *Server) handleListCampaignCharacters(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "campaign ID is required", nil)
		return
	}

	// Check if user is DM
	isDM, err := store.IsUserDM(r.Context(), s.db, id, user.ID)
	if err != nil {
		s.logger.Error("check DM", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check permissions", nil)
		return
	}
	if !isDM {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "only DM can view campaign characters", nil)
		return
	}

	chars, err := store.GetCampaignCharacters(r.Context(), s.db, id)
	if err != nil {
		s.logger.Error("list campaign characters", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list characters", nil)
		return
	}

	var resp []CharacterResponse
	for _, c := range chars {
		var draft BuildRequest
		_ = json.Unmarshal(c.Draft, &draft)
		var sheet *Sheet
		if c.Sheet != nil {
			var s Sheet
			_ = json.Unmarshal(c.Sheet, &s)
			sheet = &s
		}
		var live *SheetLive
		if c.Live != nil {
			var l SheetLive
			_ = json.Unmarshal(c.Live, &l)
			live = &l
		}
		resp = append(resp, CharacterResponse{
			ID:         c.ID,
			Name:       c.Name,
			IsNpc:      c.IsNpc,
			CampaignID: c.CampaignID.String,
			Draft:      draft,
			Sheet:      sheet,
			Live:       live,
			CreatedAt:  c.CreatedAt,
			UpdatedAt:  c.UpdatedAt,
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleAddCampaignMember adds a member to a campaign (DM only).
func (s *Server) handleAddCampaignMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "campaign ID is required", nil)
		return
	}

	// Check if user is DM
	isDM, err := store.IsUserDM(r.Context(), s.db, id, user.ID)
	if err != nil {
		s.logger.Error("check DM", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check permissions", nil)
		return
	}
	if !isDM {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "only DM can add members", nil)
		return
	}

	var req AddMemberRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body", nil)
		return
	}

	if req.UserID == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "user_id and role are required", nil)
		return
	}

	if req.Role != "player" && req.Role != "dm" {
		writeError(w, http.StatusBadRequest, "INVALID_ROLE", "role must be 'player' or 'dm'", nil)
		return
	}

	if err := store.AddCampaignMember(r.Context(), s.db, id, req.UserID, req.Role); err != nil {
		s.logger.Error("add member", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add member", nil)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "member added successfully"})
}

// handleRemoveCampaignMember removes a member from a campaign (DM only).
func (s *Server) handleRemoveCampaignMember(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "campaign ID is required", nil)
		return
	}

	userID := chi.URLParam(r, "userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "user ID is required", nil)
		return
	}

	// Check if user is DM
	isDM, err := store.IsUserDM(r.Context(), s.db, id, user.ID)
	if err != nil {
		s.logger.Error("check DM", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check permissions", nil)
		return
	}
	if !isDM {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "only DM can remove members", nil)
		return
	}

	if err := store.RemoveCampaignMember(r.Context(), s.db, id, userID); err != nil {
		s.logger.Error("remove member", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to remove member", nil)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}