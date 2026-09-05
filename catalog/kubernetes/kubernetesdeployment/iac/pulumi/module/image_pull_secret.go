package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// imagePullSecret materializes the registry logins the workload declares on
// pod.image_registries into ONE kubernetes.io/dockerconfigjson Secret named
// <workload>-image-pull, in the workload's namespace — the twin of the
// <workload>-env-secrets Secret built from the containers' env.secrets. The pod
// references it beside any spec-listed pod.image_pull_secrets.
//
// Returns nil when the pod declares no registry: a public image, a same-cloud
// registry the cluster's own identity reaches, or a Secret declared beside the
// workload need nothing from here. The data itself is built by
// workloadpod.BuildImagePullSecretData (see locals), so this file only owns the
// resource.
func imagePullSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, namespaceDeps []pulumi.ResourceOption) (*kubernetescorev1.Secret, error) {

	// If no image pull secret data is configured, return nil
	if locals.ImagePullSecretData == nil {
		return nil, nil
	}

	// Create image pull secret resource with computed name to avoid conflicts
	secretArgs := &kubernetescorev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.ImagePullSecretName),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Type:       pulumi.String("kubernetes.io/dockerconfigjson"),
		StringData: pulumi.ToStringMap(locals.ImagePullSecretData),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)
	createdImagePullSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.ImagePullSecretName,
		secretArgs,
		opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create image pull secret")
	}

	return createdImagePullSecret, nil
}
