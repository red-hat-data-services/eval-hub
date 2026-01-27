package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
)

func (h *Handlers) HandleHealth(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
