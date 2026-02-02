package k8s

// Contains the configuration logic that prepares the data needed by the builders
import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eval-hub/eval-hub/pkg/api"
)

const (
	defaultCPURequest      = "250m"
	defaultMemoryRequest   = "512Mi"
	defaultCPULimit        = "1"
	defaultMemoryLimit     = "2Gi"
	defaultNamespace       = "default"
	evalHubServiceEnv      = "EVALHUB_SERVICE_URL"
	inClusterNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

type jobConfig struct {
	jobID             string
	namespace         string
	providerID        string
	benchmarkID       string
	retryAttempts     int
	adapterImage      string
	evalHubServiceURL string
	defaultEnv        []api.EnvVar
	cpuRequest        string
	memoryRequest     string
	cpuLimit          string
	memoryLimit       string
	jobSpecJSON       string
}

func buildJobConfig(evaluation *api.EvaluationJobResource, provider *api.ProviderResource, benchmarkID string) (*jobConfig, error) {
	specJSON, err := json.MarshalIndent(evaluation, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal job spec: %w", err)
	}

	runtime := provider.Runtime
	if runtime == nil || runtime.K8s == nil {
		return nil, fmt.Errorf("provider %q missing runtime configuration", provider.ProviderID)
	}

	cpuRequest := defaultIfEmpty(runtime.K8s.CPURequest, defaultCPURequest)
	memoryRequest := defaultIfEmpty(runtime.K8s.MemoryRequest, defaultMemoryRequest)
	cpuLimit := defaultIfEmpty(runtime.K8s.CPULimit, defaultCPULimit)
	memoryLimit := defaultIfEmpty(runtime.K8s.MemoryLimit, defaultMemoryLimit)

	if runtime.K8s.Image == "" {
		return nil, fmt.Errorf("runtime adapter image is required")
	}
	evalHubServiceURL := os.Getenv(evalHubServiceEnv)
	if evalHubServiceURL == "" {
		return nil, fmt.Errorf("%s is required", evalHubServiceEnv)
	}

	retryAttempts := 0
	if evaluation.RetryAttempts != nil {
		if *evaluation.RetryAttempts < 0 {
			return nil, fmt.Errorf("retry attempts cannot be negative")
		}
		retryAttempts = *evaluation.RetryAttempts
	}
	namespace := resolveNamespace("")

	return &jobConfig{
		jobID:             evaluation.Resource.ID,
		namespace:         namespace,
		providerID:        provider.ProviderID,
		benchmarkID:       benchmarkID,
		retryAttempts:     retryAttempts,
		adapterImage:      runtime.K8s.Image,
		evalHubServiceURL: evalHubServiceURL,
		defaultEnv:        runtime.K8s.Env,
		cpuRequest:        cpuRequest,
		memoryRequest:     memoryRequest,
		cpuLimit:          cpuLimit,
		memoryLimit:       memoryLimit,
		jobSpecJSON:       string(specJSON),
	}, nil
}

func defaultIfEmpty(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func resolveNamespace(configured string) string {
	if configured != "" {
		return configured
	}
	inClusterNamespace := readInClusterNamespace()
	if inClusterNamespace != "" {
		return inClusterNamespace
	}
	return defaultNamespace
}

func readInClusterNamespace() string {
	content, err := os.ReadFile(inClusterNamespaceFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}
