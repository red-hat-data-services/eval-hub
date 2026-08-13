package server_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eval-hub/eval-hub/auth"
	"github.com/eval-hub/eval-hub/internal/eval_hub/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// newTestClientset returns a *kubernetes.Clientset which the API server is a local
// httptest server. The server is closed when the test ends. Because the tests
// below never reach AuthenticateRequest, the handler content is not material.
func newTestClientset(t *testing.T) *kubernetes.Clientset {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(ts.Close)
	cfg := &rest.Config{Host: ts.URL}
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Fatalf("kubernetes.NewForConfig: %v", err)
	}
	return cs
}

func minimalAuthConfig() *auth.AuthConfig {
	cfg := &auth.AuthConfig{}
	cfg.Authorization.Endpoints = []auth.Endpoint{
		{
			Path:      "/api/v1/evaluations/jobs",
			PathParts: []string{"", "api", "v1", "evaluations", "jobs"},
			Mappings: []auth.Mapping{
				{
					Methods: []string{"post", "get"},
					Resources: []auth.ResourceRule{
						{
							ResourceAttributes: auth.ResourceAttributes{
								Namespace: "default",
								APIGroup:  "trustyai.opendatahub.io",
								Resource:  "evaluations",
								Verb:      "get",
							},
						},
					},
				},
			},
		},
	}
	return cfg
}

func TestWithAuthentication_EmptyEndpoints(t *testing.T) {
	logger := slog.Default()
	cfg := &auth.AuthConfig{}

	_, err := server.WithAuthentication(http.NotFoundHandler(), logger, nil, cfg)
	if err == nil {
		t.Fatal("expected error when auth config has no endpoints, got nil")
	}
}

func TestWithAuthentication_PublicPathPassesThrough(t *testing.T) {
	logger := slog.Default()
	cs := newTestClientset(t)

	reached := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})

	handler, err := server.WithAuthentication(inner, logger, cs, minimalAuthConfig())
	if err != nil {
		t.Fatalf("WithAuthentication: %v", err)
	}

	for _, path := range []string{"/api/v1/health", "/metrics", "/openapi.yaml", "/docs"} {
		reached = false
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if !reached {
			t.Errorf("path %q: expected inner handler to be called, but it was not", path)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("path %q: expected 200, got %d", path, rec.Code)
		}
	}
}

func TestWithAuthentication_UnlistedPathDenied(t *testing.T) {
	logger := slog.Default()
	cs := newTestClientset(t)

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler, err := server.WithAuthentication(inner, logger, cs, minimalAuthConfig())
	if err != nil {
		t.Fatalf("WithAuthentication: %v", err)
	}

	// A path that exists in the mux but is not listed in auth.yaml should be
	// denied rather than passed through unauthenticated.
	for _, path := range []string{"/api/v1/unlisted", "/api/v1/evaluations/jobs/123/status"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("path %q: expected 401, got %d", path, rec.Code)
		}
	}
}
