package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hadnu/onatar/internal/store"
)

// handleListCharacters returns all characters for the authenticated user.
func (s *Server) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	chars, err := store.ListCharactersByUser(ctx, s.db, user.ID)
	if err != nil {
		s.logger.Error("list characters", "error", err)
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

// handleCreateCharacter creates a new character for the authenticated user.
func (s *Server) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	var req BuildRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DRAFT", "invalid request body", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), buildTimeout)
	defer cancel()

	sheet, err := s.buildSheet(ctx, req)
	if err != nil {
		s.logger.Error("build failed", "error", err)
		writeError(w, http.StatusUnprocessableEntity, "BUILD_ERROR", err.Error(), nil)
		return
	}

	// Create character in database
	draftBytes, _ := json.Marshal(req)
	sheetBytes, _ := json.Marshal(sheet)

	char := &store.Character{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Name:      req.Name,
		IsNpc:     req.IsNpc,
		Draft:     draftBytes,
		Sheet:     sheetBytes,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := store.CreateCharacter(r.Context(), s.db, char); err != nil {
		s.logger.Error("create character", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create character", nil)
		return
	}

	resp := CharacterResponse{
		ID:        char.ID,
		Name:      char.Name,
		IsNpc:     char.IsNpc,
		Draft:     req,
		Sheet:     sheet,
		CreatedAt: char.CreatedAt,
		UpdatedAt: char.UpdatedAt,
	}

	writeJSON(w, http.StatusCreated, resp)
}

// handleGetCharacter returns a character by ID.
func (s *Server) handleGetCharacter(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "character ID is required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	char, err := store.GetCharacter(ctx, s.db, id)
	if err != nil {
		s.logger.Error("get character", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get character", nil)
		return
	}

	if char == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "character not found", nil)
		return
	}

	// Check ownership
	if char.UserID != user.ID {
		writeError(w, http.StatusForbidden, "FORBIDDEN", "not the owner of this character", nil)
		return
	}

	var draft BuildRequest
	_ = json.Unmarshal(char.Draft, &draft)
	var sheet *Sheet
	if char.Sheet != nil {
		var s Sheet
		_ = json.Unmarshal(char.Sheet, &s)
		sheet = &s
	}
	var live *SheetLive
	if char.Live != nil {
		var l SheetLive
		_ = json.Unmarshal(char.Live, &l)
		live = &l
	}

	resp := CharacterResponse{
		ID:         char.ID,
		Name:       char.Name,
		IsNpc:      char.IsNpc,
		CampaignID: char.CampaignID.String,
		Draft:      draft,
		Sheet:      sheet,
		Live:       live,
		CreatedAt:  char.CreatedAt,
		UpdatedAt:  char.UpdatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleUpdateCharacter updates a character.
func (s *Server) handleUpdateCharacter(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "character ID is required", nil)
		return
	}

	var req BuildRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DRAFT", "invalid request body", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), buildTimeout)
	defer cancel()

	sheet, err := s.buildSheet(ctx, req)
	if err != nil {
		s.logger.Error("build failed", "error", err)
		writeError(w, http.StatusUnprocessableEntity, "BUILD_ERROR", err.Error(), nil)
		return
	}

	// Update character in database
	draftBytes, _ := json.Marshal(req)
	sheetBytes, _ := json.Marshal(sheet)

	char := &store.Character{
		ID:        id,
		UserID:    user.ID,
		Name:      req.Name,
		IsNpc:     req.IsNpc,
		Draft:     draftBytes,
		Sheet:     sheetBytes,
		UpdatedAt: time.Now().Unix(),
	}

	if err := store.UpdateCharacter(r.Context(), s.db, char); err != nil {
		s.logger.Error("update character", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update character", nil)
		return
	}

	resp := CharacterResponse{
		ID:        char.ID,
		Name:      char.Name,
		IsNpc:     char.IsNpc,
		Draft:     req,
		Sheet:     sheet,
		UpdatedAt: char.UpdatedAt,
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleDeleteCharacter deletes a character.
func (s *Server) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated", nil)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ID", "character ID is required", nil)
		return
	}

	if err := store.DeleteCharacter(r.Context(), s.db, user.ID, id); err != nil {
		s.logger.Error("delete character", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to delete character", nil)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}