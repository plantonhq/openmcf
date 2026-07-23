// Creates the kubernetes.core.v1.ServiceAccount resource from computed locals.
package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createServiceAccount creates the Kubernetes ServiceAccount resource
func createServiceAccount(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.ServiceAccount, error) {
	serviceAccountArgs := &kubernetescorev1.ServiceAccountArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:        pulumi.String(locals.ServiceAccountName),
			Namespace:   pulumi.String(locals.Namespace),
			Labels:      pulumi.ToStringMap(locals.Labels),
			Annotations: pulumi.ToStringMap(locals.Annotations),
		},
	}

	// Attach pull secrets as LocalObjectReferences (name-only refs to secrets in the
	// same namespace). Left nil when empty so the field is absent from the manifest.
	if len(locals.ImagePullSecretNames) > 0 {
		imagePullSecrets := kubernetescorev1.LocalObjectReferenceArray{}
		for _, secretName := range locals.ImagePullSecretNames {
			imagePullSecrets = append(imagePullSecrets, &kubernetescorev1.LocalObjectReferenceArgs{
				Name: pulumi.String(secretName),
			})
		}
		serviceAccountArgs.ImagePullSecrets = imagePullSecrets
	}

	// Tri-state automount: set the field only when the user made an explicit choice.
	// When unset (nil) the field is omitted from the ServiceAccount entirely, so the
	// cluster/pod-level default applies — which in Kubernetes is "mount the token".
	// Note: the Terraform provider cannot omit this attribute (its schema defaults it
	// to true); both modules converge on identical observable behavior because the
	// cluster default for an absent field IS true. Documented in both modules.
	if locals.AutomountServiceAccountToken != nil {
		serviceAccountArgs.AutomountServiceAccountToken = pulumi.BoolPtr(*locals.AutomountServiceAccountToken)
	}

	serviceAccount, err := kubernetescorev1.NewServiceAccount(
		ctx,
		locals.ServiceAccountName,
		serviceAccountArgs,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service account %s", locals.ServiceAccountName)
	}

	return serviceAccount, nil
}
