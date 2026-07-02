package module

import (
	"github.com/pkg/errors"
	awslblistenerv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslblistener/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awslblistenerv1.AwsLbListenerStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsLbListener.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := listener(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create listener")
	}

	return nil
}
