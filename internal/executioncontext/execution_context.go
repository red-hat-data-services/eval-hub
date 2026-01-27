package executioncontext

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// ExecutionContext contains execution context for API operations. This pattern enables
// type-safe passing of configuration and state to evaluation-related handlers, which
// receive an ExecutionContext instead of a raw http.Request.
//
// The ExecutionContext contains:
//   - Logger: A request-scoped logger with enriched fields (request_id, method, uri, etc.)
//   - Config: The service configuration
//   - Evaluation-specific state: model info, timeouts, retries, metadata
type ExecutionContext struct {
	Ctx           context.Context
	RequestID     string
	Logger        *slog.Logger
	Method        string
	URI           string
	BaseURL       string
	RawQuery      string
	headers       map[string][]string
	body          io.ReadCloser
	bodyBytesRead bool
	EvaluationID  string
	ModelURL      string
	ModelName     string
	//BackendSpec    BackendSpec
	//BenchmarkSpec  BenchmarkSpec
	Timeout        time.Duration
	RetryAttempts  int
	StartedAt      time.Time
	Metadata       map[string]interface{}
	MLflowClient   interface{}
	ExperimentName string
}

func NewExecutionContext(
	ctx context.Context,
	requestID string,
	logger *slog.Logger,
	method string,
	uri string,
	baseURL string,
	rawQuery string,
	headers map[string][]string,
	body io.ReadCloser,
	evaluationID string,
	modelURL string,
	modelName string,
	timeout time.Duration,
	retryAttempts int,
	metadata map[string]interface{},
	mlflowClient interface{},
	experimentName string,
) *ExecutionContext {
	return &ExecutionContext{
		Ctx:            ctx,
		RequestID:      requestID,
		Logger:         logger,
		Method:         method,
		URI:            uri,
		BaseURL:        baseURL,
		RawQuery:       rawQuery,
		headers:        headers,
		body:           body,
		bodyBytesRead:  false,
		EvaluationID:   evaluationID,
		ModelURL:       modelURL,
		ModelName:      modelName,
		Timeout:        timeout,
		RetryAttempts:  retryAttempts,
		StartedAt:      time.Now(),
		Metadata:       metadata,
		MLflowClient:   mlflowClient,
		ExperimentName: experimentName,
	}
}

func (ctx *ExecutionContext) GetHeader(key string) string {
	if (ctx.headers != nil) && (ctx.headers[key] != nil) && len(ctx.headers[key]) > 0 {
		return ctx.headers[key][0]
	}
	return ""
}

func (ctx *ExecutionContext) SetHeader(key string, value string) {
	if ctx.headers == nil {
		ctx.headers = make(map[string][]string)
	}
	ctx.headers[key] = []string{value}
}

func (ctx *ExecutionContext) GetBody() io.ReadCloser {
	return ctx.body
}

func (ctx *ExecutionContext) GetBodyAsBytes() ([]byte, error) {
	// we could save the body bytes if we need to read this multiple times
	// but for now we just allow a single read
	if ctx.bodyBytesRead {
		return nil, fmt.Errorf("body bytes already read")
	}
	bodyBytes, err := io.ReadAll(ctx.body)
	if err != nil {
		return nil, err
	}
	ctx.bodyBytesRead = true
	return bodyBytes, nil
}
