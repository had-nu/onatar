package httpapi

import (
	"context"
	"net/http"
)

// handleContent returns all rule content from MariaDB (PRD §3.5 GET /content).
func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), contentTimeout)
	defer cancel()

	c, err := s.contentLoader.LoadContent(ctx)
	if err != nil {
		s.logger.Error("load content", "error", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error", nil)
		return
	}
	writeJSON(w, http.StatusOK, c)
}
