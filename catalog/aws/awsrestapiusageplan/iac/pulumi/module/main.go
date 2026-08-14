package module

import (
	"github.com/pkg/errors"
	awsrestapiusageplanv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapiusageplan/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the usage plan with its API keys
// and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsrestapiusageplanv1alpha1.AwsRestApiUsagePlanStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := usagePlan(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "rest api usage plan")
	}

	return nil
}
