package httpapi

import (
	"context"
	"errors"

	"github.com/hadnu/onatar/internal/build"
)

// validateBuildRequest validates the build request.
func validateBuildRequest(req BuildRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Classes) == 0 {
		return errors.New("at least one class is required")
	}
	for _, c := range req.Classes {
		if c.ID == "" {
			return errors.New("class ID is required")
		}
		if c.Level < 1 || c.Level > 20 {
			return errors.New("class level must be between 1 and 20")
		}
	}
	return nil
}

// buildSheet builds a character sheet from a build request.
func (s *Server) buildSheet(ctx context.Context, req BuildRequest) (*Sheet, error) {
	rules, err := s.contentLoader.LoadContent(ctx)
	if err != nil {
		return nil, err
	}

	// Convert BuildRequest to build.Request
	buildReq := build.Request{
		Name:           req.Name,
		Classes:        make([]build.ClassInput, len(req.Classes)),
		SpeciesID:      req.SpeciesID,
		BackgroundID:   req.BackgroundID,
		AbilityScores:  req.AbilityScores,
		AbilityMethod:  req.AbilityMethod,
		Skills:         req.Skills,
		Spells:         req.Spells,
		Feats:          req.Feats,
		IsNPC:          req.IsNpc,
	}

	for i, c := range req.Classes {
		buildReq.Classes[i] = build.ClassInput{
			ID:         c.ID,
			Level:      c.Level,
			SubclassID: c.SubclassID,
		}
	}

	resp, err := build.Compute(rules, buildReq)
	if err != nil {
		var re *build.RuleError
		if errors.As(err, &re) {
			return nil, re
		}
		return nil, err
	}

	// Convert build.Response to httpapi.Sheet
	sheet := &Sheet{
		Level:            resp.Sheet.Level,
		HP:               struct{ Max, Current int }{Max: resp.Sheet.HP.Max, Current: resp.Sheet.HP.Current},
		AC:               resp.Sheet.AC,
		ProficiencyBonus: resp.Sheet.ProficiencyBonus,
		SpellSlots:       resp.Sheet.SpellSlots,
		Abilities:        make(map[string]AbilityScore),
		Features:         make([]SheetFeature, len(resp.Sheet.Features)),
		PendingChoices:   make([]PendingChoice, len(resp.Sheet.PendingChoices)),
	}

	for k, v := range resp.Sheet.Abilities {
		sheet.Abilities[k] = AbilityScore{Score: v.Score, Mod: v.Mod}
	}

	for i, f := range resp.Sheet.Features {
		sheet.Features[i] = SheetFeature{
			Name:        f.Name,
			Level:       f.Level,
			Description: f.Description,
		}
	}

	for i, pc := range resp.Sheet.PendingChoices {
		sheet.PendingChoices[i] = PendingChoice{
			Type:        pc.Type,
			Description: pc.Description,
		}
	}

	return sheet, nil
}