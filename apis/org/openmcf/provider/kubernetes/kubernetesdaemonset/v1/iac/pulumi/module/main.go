package module

import (
	"github.com/pkg/errors"
	kubernetesdaemonsetv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes/kubernetesdaemonset/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the main entry point for creating all necessary Kubernetes resources
// for a DaemonSet, based on the KubernetesDaemonSet API resource definition.
func Resources(ctx *pulumi.Context, stackInput *kubernetesdaemonsetv1.KubernetesDaemonSetStackInput) error {
	// Initialize local references
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	// Create a Pulumi Kubernetes provider from the given credentials
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx,
		stackInput.ProviderConfig,
		"kubernetes",
	)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes provider")
	}

	// Conditionally create namespace based on create_namespace flag
	if stackInput.Target.Spec.CreateNamespace {
		_, err = kubernetescorev1.NewNamespace(ctx,
			locals.Namespace,
			&kubernetescorev1.NamespaceArgs{
				Metadata: kubernetesmetav1.ObjectMetaPtrInput(
					&kubernetesmetav1.ObjectMetaArgs{
						Name:   pulumi.String(locals.Namespace),
						Labels: pulumi.ToStringMap(locals.Labels),
					}),
			}, pulumi.Provider(kubernetesProvider))
		if err != nil {
			return errors.Wrapf(err, "failed to create %s namespace", locals.Namespace)
		}
	}

	// Create ConfigMaps
	_, err = configMaps(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create configmaps")
	}

	// Create ServiceAccount
	serviceAccountName, err := serviceAccount(ctx, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create service account")
	}

	// Create RBAC resources
	if err := rbac(ctx, locals, serviceAccountName, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create rbac resources")
	}

	// Create the main secret resource for environment secrets
	if err := secret(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	// Create an image pull secret if Docker credentials are provided
	if locals.ImagePullSecretData != nil {
		if err := imagePullSecret(ctx, locals, kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to create image pull secret")
		}
	}

	// Create the DaemonSet resource
	if err := daemonSet(ctx, locals, serviceAccountName, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create daemonset")
	}

	return nil
}
