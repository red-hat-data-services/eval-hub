package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/eval-hub/eval-hub/auth"
	"github.com/eval-hub/eval-hub/internal/eval_hub/messages"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/kubernetes"
)

// publicPaths are served without authentication regardless of auth config.
// These are infrastructure endpoints that must remain accessible to health
// checkers, monitoring, and API documentation consumers.
var publicPaths = []string{
	"/api/v1/health",
	"/metrics",
	"/openapi.yaml",
	"/docs",
}

func isPublicPath(path string) bool {
	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// WithAuthentication authenticates requests for endpoints listed in the auth
// config. Public infrastructure paths pass through unauthenticated. Any path
// not in the config and not in the public allowlist is rejected with 401.
func WithAuthentication(next http.Handler, logger *slog.Logger, client *kubernetes.Clientset, config *auth.AuthConfig) (http.Handler, error) {
	if len(config.Authorization.Endpoints) == 0 {
		return nil, fmt.Errorf("auth config has no endpoints defined: refusing to start with authentication enabled but no authorization rules")
	}

	authn, err := auth.NewAuthenticator(client, logger)
	if err != nil {
		logger.Error("Error creating authenticator", "error", err)
		return nil, err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Authenticating request", "path", r.URL.Path, "method", r.Method)

		if isPublicPath(r.URL.Path) {
			logger.Debug("Public path, skipping authentication", "path", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		rules := auth.FindRules(r, config)
		if len(rules) == 0 {
			logger.Warn("No auth rules for request, denying", "path", r.URL.Path, "method", r.Method)
			writeError(w, messages.Unauthorized)
			return
		}

		resp, ok, err := authn.AuthenticateRequest(r)
		if err != nil {
			logger.Error("Error authenticating request", "error", err)
			writeError(w, messages.Unauthorized, "Error", err.Error())
			return
		}
		if !ok {
			logger.Error("Request not authenticated", "path", r.URL.Path, "method", r.Method)
			writeError(w, messages.Unauthorized)
			return
		}

		r = r.WithContext(request.WithUser(r.Context(), resp.User))
		next.ServeHTTP(w, r)
	})

	return handler, nil
}
