package module

import (
	"github.com/pkg/errors"
	kubernetespriorityclassv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetespriorityclass/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the main entry point for the Pulumi module.
// It orchestrates the creation of a Kubernetes PriorityClass with its
// priority value, default flag, and preemption policy.
func Resources(ctx *pulumi.Context, stackInput *kubernetespriorityclassv1.KubernetesPriorityClassStackInput) error {
	// Initialize locals with derived values
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

	// Create the priority class
	if _, err := createPriorityClass(ctx, locals, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create priority class")
	}

	// Export outputs
	if err := exportOutputs(ctx, locals); err != nil {
		return errors.Wrap(err, "failed to export outputs")
	}

	return nil
}
