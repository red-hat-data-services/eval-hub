package k8s

import (
	"log/slog"

	"github.com/eval-hub/eval-hub/internal/abstractions"
	"github.com/eval-hub/eval-hub/pkg/api"
)

type K8sRuntime struct {
	logger *slog.Logger
}

func NewK8sRuntime(logger *slog.Logger) (abstractions.Runtime, error) {
	return &K8sRuntime{logger: logger}, nil
}

func (r *K8sRuntime) RunEvaluationJob(evaluation *api.EvaluationJobResource, storage *abstractions.Storage) error {
	return nil
}

func (r *K8sRuntime) Name() string {
	return "kubernetes"
}
