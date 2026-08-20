package module

import (
	"github.com/pkg/errors"
	awsrdsproxyv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrdsproxy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the proxy, its pool tuning,
// endpoints, and target, and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsrdsproxyv1alpha1.AwsRdsProxyStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := proxy(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "proxy")
	}

	return nil
}
