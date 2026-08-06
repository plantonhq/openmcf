package module

import (
	"github.com/pkg/errors"
	awsapprunnervpcconnectorv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsapprunnervpcconnector/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates VPC connector creation and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsapprunnervpcconnectorv1alpha1.AwsAppRunnerVpcConnectorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsAppRunnerVpcConnector.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := vpcConnector(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "vpc connector")
	}

	return nil
}
