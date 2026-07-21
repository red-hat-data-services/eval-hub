package handlers

import (
	"context"
	"fmt"
	"strings"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/eval-hub/eval-hub/pkg/api"
)

func collectSecretRefs(evaluation *api.EvaluationJobConfig) []string {
	var names []string
	if evaluation.Model.Auth != nil {
		if s := strings.TrimSpace(evaluation.Model.Auth.SecretRef); s != "" {
			names = append(names, s)
		}
	}
	if evaluation.Exports != nil && evaluation.Exports.OCI != nil && evaluation.Exports.OCI.K8s != nil {
		if s := strings.TrimSpace(evaluation.Exports.OCI.K8s.Connection); s != "" {
			names = append(names, s)
		}
	}
	for _, b := range evaluation.Benchmarks {
		if b.TestDataRef != nil && b.TestDataRef.S3 != nil {
			if s := strings.TrimSpace(b.TestDataRef.S3.SecretRef); s != "" {
				names = append(names, s)
			}
		}
	}
	return names
}

func checkSecretAccess(ctx context.Context, client kubernetes.Interface, username, namespace string, secretNames []string) error {
	for _, name := range secretNames {
		sar := &authorizationv1.SubjectAccessReview{
			Spec: authorizationv1.SubjectAccessReviewSpec{
				User: username,
				ResourceAttributes: &authorizationv1.ResourceAttributes{
					Namespace: namespace,
					Verb:      "get",
					Group:     "",
					Resource:  "secrets",
					Name:      name,
				},
			},
		}
		result, err := client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("secret access review failed: %w", err)
		}
		if !result.Status.Allowed {
			return fmt.Errorf("caller %q is not permitted to access secret %q in namespace %q", username, name, namespace)
		}
	}
	return nil
}
