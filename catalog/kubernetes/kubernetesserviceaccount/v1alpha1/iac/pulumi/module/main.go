// Orchestrator for the KubernetesServiceAccount Pulumi module:
// locals initialization, provider setup, resource creation, and output export.
package module

import (
	"github.com/pkg/errors"
	kubernetesserviceaccountv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesserviceaccount/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the main entry point for the Pulumi module.
// It orchestrates the creation of a Kubernetes ServiceAccount with image-pull
// secrets, the token-automount setting, and workload-identity annotations.
func Resources(ctx *pulumi.Context, stackInput *kubernetesserviceaccountv1alpha1.KubernetesServiceAccountStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Create Kubernetes provider from credentials
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx,
		stackInput.ProviderConfig,
		"kubernetes",
	)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	// Create the service account
	if _, err := createServiceAccount(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create service account")
	}

	// Export outputs
	if err := exportOutputs(ctx, locals); err != nil {
		return errors.Wrap(err, "failed to export outputs")
	}

	return nil
}
