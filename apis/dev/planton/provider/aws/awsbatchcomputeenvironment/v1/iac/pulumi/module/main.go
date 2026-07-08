package module

import (
	"github.com/pkg/errors"
	awsbatchcomputeenvironmentv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsbatchcomputeenvironment/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates the AWS Batch compute environment and exports its
// outputs. Job queues and scheduling policies are separate resources
// (AwsBatchJobQueue / AwsBatchSchedulingPolicy) that compose onto the
// environment through its exported ARN.
func Resources(ctx *pulumi.Context, stackInput *awsbatchcomputeenvironmentv1.AwsBatchComputeEnvironmentStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsBatchComputeEnvironment.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdCe, err := computeEnvironment(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "compute environment")
	}

	// --- Exports ---
	ctx.Export(OpComputeEnvironmentArn, createdCe.Arn)
	ctx.Export(OpComputeEnvironmentName, createdCe.Name)
	ctx.Export(OpEcsClusterArn, createdCe.EcsClusterArn)
	ctx.Export(OpStatus, createdCe.Status)

	return nil
}
