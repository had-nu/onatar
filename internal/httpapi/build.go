package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hadnu/onatar/internal/build"
)

// maxBodyBytes bounds the /build request size.
const maxBodyBytes = 1 << 20 // 1 MiB

// handleBuild derives a character sheet from a draft (PRD §3.5 POST /build).
func (s *Server) handleBuild(w http.ResponseWriter, r *http.Request) {
	var req build.Request
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, string(build.KindInvalidDraft), "invalid JSON body", nil)
		return
	}
	if dec.More() {
		writeError(w, http.StatusBadRequest, string(build.KindInvalidDraft), "unexpected trailing data", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), buildTimeout)
	defer cancel()

	rules, err := s.contentLoader.LoadContent(ctx)
	if err != nil {
		s.logger.Error("load content for build", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
		return
	}

	resp, err := build.Compute(rules, req)
	if err != nil {
		var re *build.RuleError
		if errors.As(err, &re) {
			status := http.StatusUnprocessableEntity
			if re.Kind == build.KindInvalidDraft || re.Kind == build.KindInvalidSpell {
				status = http.StatusBadRequest
			}
			writeError(w, status, string(re.Kind), re.Message, re.Details)
			return
		}
		s.logger.Error("build", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
