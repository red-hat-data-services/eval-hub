package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
)

func (h *Handlers) HandleStatus(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service":   "eval-hub",
		"version":   "1.0.0",
		"status":    "running",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
