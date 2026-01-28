package handlers_test

import (
	"testing"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
	"github.com/eval-hub/eval-hub/internal/handlers"
)

func TestNew(t *testing.T) {
	h := handlers.New(nil, nil)
	if h == nil {
		t.Error("New() returned nil")
	}
}

func createExecutionContext(method string, uri string) *executioncontext.ExecutionContext {
	return &executioncontext.ExecutionContext{
		Method: method,
		URI:    uri,
	}
}
