package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
	"github.com/eval-hub/eval-hub/internal/serialization"
	"github.com/eval-hub/eval-hub/pkg/api"
)

// BackendSpec represents the backend specification
type BackendSpec struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// BenchmarkSpec represents the benchmark specification
type BenchmarkSpec struct {
	BenchmarkID string                 `json:"benchmark_id"`
	ProviderID  string                 `json:"provider_id"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// HandleCreateEvaluation handles POST /api/v1/evaluations/jobs
func (h *Handlers) HandleCreateEvaluation(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodPost, w) {
		return
	}
	// get the body bytes from the context
	bodyBytes, err := ctx.GetBodyAsBytes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	evaluation := &api.EvaluationJobConfig{}
	err = serialization.Unmarshal(h.validate, ctx, bodyBytes, evaluation)
	if err != nil {
		h.serializationError(ctx, w, err, http.StatusBadRequest)
		return
	}

	response, err := h.storage.CreateEvaluationJob(ctx, evaluation)
	if err != nil {
		h.errorResponse(ctx, w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.successResponse(ctx, w, response, http.StatusAccepted)
}

// HandleListEvaluations handles GET /api/v1/evaluations/jobs
func (h *Handlers) HandleListEvaluations(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":       []interface{}{},
		"total_count": 0,
		"limit":       50,
		"first":       map[string]string{"href": ""},
		"next":        nil,
	})
}

// HandleGetEvaluation handles GET /api/v1/evaluations/jobs/{id}
func (h *Handlers) HandleGetEvaluation(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	// Extract ID from path
	pathParts := strings.Split(ctx.URI, "/")
	id := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation retrieval not yet implemented",
		"id":      id,
	})
}

// HandleCancelEvaluation handles DELETE /api/v1/evaluations/jobs/{id}
func (h *Handlers) HandleCancelEvaluation(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodDelete, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation cancellation not yet implemented",
	})
}

// HandleGetEvaluationSummary handles GET /api/v1/evaluations/jobs/{id}/summary
func (h *Handlers) HandleGetEvaluationSummary(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Evaluation summary not yet implemented",
	})
}

// HandleListBenchmarks handles GET /api/v1/evaluations/benchmarks
func (h *Handlers) HandleListBenchmarks(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"benchmarks":         []interface{}{},
		"total_count":        0,
		"providers_included": []string{},
	})
}

// HandleListCollections handles GET /api/v1/evaluations/collections
func (h *Handlers) HandleListCollections(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"collections":       []interface{}{},
		"total_collections": 0,
	})
}

// HandleCreateCollection handles POST /api/v1/evaluations/collections
func (h *Handlers) HandleCreateCollection(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodPost, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection creation not yet implemented",
	})
}

// HandleGetCollection handles GET /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleGetCollection(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	// Extract collection_id from path
	pathParts := strings.Split(ctx.URI, "/")
	collectionID := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Collection retrieval not yet implemented",
		"collection_id": collectionID,
	})
}

// HandleUpdateCollection handles PUT /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleUpdateCollection(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodPut, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection update not yet implemented",
	})
}

// HandlePatchCollection handles PATCH /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandlePatchCollection(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodPatch, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection patch not yet implemented",
	})
}

// HandleDeleteCollection handles DELETE /api/v1/evaluations/collections/{collection_id}
func (h *Handlers) HandleDeleteCollection(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodDelete, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Collection deletion not yet implemented",
	})
}

// HandleListProviders handles GET /api/v1/evaluations/providers
func (h *Handlers) HandleListProviders(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"providers":        []interface{}{},
		"total_providers":  0,
		"total_benchmarks": 0,
	})
}

// HandleGetProvider handles GET /api/v1/evaluations/providers/{provider_id}
func (h *Handlers) HandleGetProvider(ctx *executioncontext.ExecutionContext, w http.ResponseWriter) {
	if !h.checkMethod(ctx, http.MethodGet, w) {
		return
	}

	// Extract provider_id from path
	pathParts := strings.Split(ctx.URI, "/")
	providerID := pathParts[len(pathParts)-1]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Provider retrieval not yet implemented",
		"provider_id": providerID,
	})
}
