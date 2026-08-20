package module

import (
	"github.com/pkg/errors"
	awsbedrockinvocationloggingv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockinvocationlogging/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources manages the region's Bedrock invocation logging
// configuration and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsbedrockinvocationloggingv1alpha1.AwsBedrockInvocationLoggingStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := invocationLogging(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "bedrock invocation logging")
	}

	return nil
}
