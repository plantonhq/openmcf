package module

import (
	"github.com/pkg/errors"
	kubernetesdaemonsetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesdaemonset/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources deploys a KubernetesDaemonSet: an optional namespace, the env and
// image-pull satellite Secrets, and the DaemonSet itself.
//
// DaemonSets have no Service, HPA, or PDB — node membership is the replica
// count, and clients reach agents on their node via host ports or host
// networking. Identity is composed, not created: pods run as the
// ServiceAccount referenced in spec.pod.service_account, and API permissions
// come from KubernetesRbac grants targeting that identity. This module never
// creates ServiceAccounts, RBAC objects, certificates, gateways, or routes.
func Resources(ctx *pulumi.Context, stackInput *kubernetesdaemonsetv1.KubernetesDaemonSetStackInput) error {
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	createdNamespace, err := namespace(ctx, stackInput, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	// Conditional namespace dependency (Pulumi equivalent of Terraform depends_on):
	// empty when the namespace pre-exists or is owned by a KubernetesNamespace resource.
	var namespaceDeps []pulumi.ResourceOption
	if createdNamespace != nil {
		namespaceDeps = append(namespaceDeps, pulumi.DependsOn([]pulumi.Resource{createdNamespace}))
	}

	// Satellite secrets are created BEFORE the DaemonSet: pods reference them by
	// name at startup, and a pod that starts before its env secret exists crashes.
	if err := secret(ctx, locals, kubernetesProvider, namespaceDeps); err != nil {
		return errors.Wrap(err, "failed to create env secret")
	}

	createdImagePullSecret, err := imagePullSecret(ctx, locals, kubernetesProvider, namespaceDeps)
	if err != nil {
		return errors.Wrap(err, "failed to create image pull secret")
	}

	if err := daemonSet(ctx, locals, kubernetesProvider, createdImagePullSecret, namespaceDeps); err != nil {
		return errors.Wrap(err, "failed to create daemon set")
	}

	return nil
}
