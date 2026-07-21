package handlers

import (
	"github.com/eval-hub/eval-hub/internal/eval_hub/abstractions"
	"github.com/eval-hub/eval-hub/internal/eval_hub/config"
	"github.com/eval-hub/eval-hub/pkg/mlflowclient"
	"github.com/go-playground/validator/v10"
	"k8s.io/client-go/kubernetes"
)

type Handlers struct {
	storage       abstractions.Storage
	validate      *validator.Validate
	runtime       abstractions.Runtime
	mlflowClient  *mlflowclient.Client
	serviceConfig *config.Config
	k8sClient     kubernetes.Interface
}

func New(storage abstractions.Storage, validate *validator.Validate, runtime abstractions.Runtime, mlflowClient *mlflowclient.Client, serviceConfig *config.Config, k8sClient kubernetes.Interface) *Handlers {
	return &Handlers{
		storage:       storage,
		validate:      validate,
		runtime:       runtime,
		mlflowClient:  mlflowClient,
		serviceConfig: serviceConfig,
		k8sClient:     k8sClient,
	}
}
