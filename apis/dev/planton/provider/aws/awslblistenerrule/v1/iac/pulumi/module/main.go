package module

import (
	"github.com/pkg/errors"
	awslblistenerrulev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awslblistenerrule/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awslblistenerrulev1.AwsLbListenerRuleStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which
	// resolves the right credential mechanism (static keys, keyless web identity,
	// or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsLbListenerRule.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := listenerRule(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "failed to create listener rule")
	}

	return nil
}
