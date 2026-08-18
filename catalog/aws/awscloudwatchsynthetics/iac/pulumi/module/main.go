package module

import (
	"github.com/pkg/errors"
	awscloudwatchsyntheticsv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudwatchsynthetics/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the canary, owned groups, and the
// canary's group joins, and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awscloudwatchsyntheticsv1alpha1.AwsCloudwatchSyntheticsStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdGroups, err := groups(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "groups")
	}

	if err := canary(ctx, locals, provider, createdGroups); err != nil {
		return errors.Wrap(err, "canary")
	}

	return nil
}
