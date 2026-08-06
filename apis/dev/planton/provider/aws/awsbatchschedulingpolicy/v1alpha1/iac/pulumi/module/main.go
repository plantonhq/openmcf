package module

import (
	"github.com/pkg/errors"
	awsbatchschedulingpolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsbatchschedulingpolicy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates the AWS Batch scheduling policy and exports its outputs.
func Resources(ctx *pulumi.Context, stackInput *awsbatchschedulingpolicyv1alpha1.AwsBatchSchedulingPolicyStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsBatchSchedulingPolicy.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdPolicy, err := schedulingPolicy(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "scheduling policy")
	}

	// --- Exports ---
	ctx.Export(OpSchedulingPolicyArn, createdPolicy.Arn)
	ctx.Export(OpSchedulingPolicyName, createdPolicy.Name)

	return nil
}
