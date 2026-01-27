package server

import (
	"context"
	"net/http"
	"time"

	"github.com/eval-hub/eval-hub/internal/executioncontext"
)

// newExecutionContext creates a new ExecutionContext with default values. This function
// is called at the route level before invoking evaluation-related handlers to set up
// request-scoped context.
//
// The function automatically:
//   - Enhances the logger with request-specific fields via logging.LoggerWithRequest
//   - Sets default timeout (60 minutes) and retry attempts (3)
//   - Initializes an empty metadata map
//
// This enables automatic request ID tracking (from X-Global-Transaction-Id header or
// auto-generated UUID) and structured logging with consistent request metadata.
//
// Parameters:
//   - r: The HTTP request to extract context from
//   - logger: The base logger to enhance with request fields
//   - serviceConfig: The service configuration to include in the context
//
// Returns:
//   - *ExecutionContext: A new execution context ready for use in handlers
func (s *Server) newExecutionContext(r *http.Request) *executioncontext.ExecutionContext {
	// Enhance logger with request-specific fields
	requestID, enhancedLogger := s.loggerWithRequest(r)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + r.Host

	return executioncontext.NewExecutionContext(
		context.Background(),
		requestID,
		enhancedLogger,
		r.Method,
		r.URL.Path,
		baseURL,
		r.URL.RawQuery,
		r.Header,
		r.Body,
		"",
		"",
		"",
		time.Minute*60,
		3,
		make(map[string]interface{}),
		nil,
		"",
	)
}
