package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
)

// HandleGetSystemMetrics handles GET /api/v1/metrics/system
func (h *Handlers) HandleGetSystemMetrics(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "System metrics not yet implemented",
	})
}
