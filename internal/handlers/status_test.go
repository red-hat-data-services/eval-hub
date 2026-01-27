package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eval-hub/eval-hub/internal/handlers"
)

func TestHandleStatus(t *testing.T) {
	h := handlers.New(nil, nil)

	t.Run("GET request returns status information", func(t *testing.T) {
		ctx := createExecutionContext(http.MethodGet, "/api/v1/status")
		w := httptest.NewRecorder()

		h.HandleStatus(ctx, w)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		expectedFields := map[string]interface{}{
			"service": "eval-hub",
			"version": "1.0.0",
			"status":  "running",
		}

		for key, expectedValue := range expectedFields {
			if response[key] != expectedValue {
				t.Errorf("Expected %s to be %v, got %v", key, expectedValue, response[key])
			}
		}

		if _, ok := response["timestamp"]; !ok {
			t.Error("Response missing timestamp field")
		}

		// Verify timestamp is valid RFC3339 format
		if timestamp, ok := response["timestamp"].(string); ok {
			if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
				t.Errorf("Invalid timestamp format: %v", err)
			}
		}
	})

	t.Run("POST request returns method not allowed", func(t *testing.T) {
		ctx := createExecutionContext(http.MethodPost, "/api/v1/status")
		w := httptest.NewRecorder()

		h.HandleStatus(ctx, w)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})
}
