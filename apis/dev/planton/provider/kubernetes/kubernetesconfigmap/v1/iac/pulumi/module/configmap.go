package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createConfigMap creates the Kubernetes ConfigMap resource.
//
// Data carries UTF-8 entries as-is. BinaryData carries values that arrive
// already base64-encoded — Kubernetes stores binaryData as base64 on the wire,
// so the strings pass through unchanged and the API server decodes them when
// serving the ConfigMap to consumers. When Immutable is true, the cluster
// rejects any update to data/binaryData after creation; changing content then
// requires delete-and-recreate.
func createConfigMap(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.ConfigMap, error) {
	configMap, err := kubernetescorev1.NewConfigMap(
		ctx,
		locals.ConfigMapName,
		&kubernetescorev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.ConfigMapName),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Data:       pulumi.ToStringMap(locals.Data),
			BinaryData: pulumi.ToStringMap(locals.BinaryData),
			Immutable:  pulumi.Bool(locals.Immutable),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create configmap %s", locals.ConfigMapName)
	}

	return configMap, nil
}
