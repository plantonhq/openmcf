package module

import (
	"github.com/pkg/errors"
	awsvpcv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsvpc/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awsvpcv1alpha1.AwsVpcStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsVpc.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := vpc(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create vpc")
	}

	return nil
}
