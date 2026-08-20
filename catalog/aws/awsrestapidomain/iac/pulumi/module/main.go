package module

import (
	"github.com/pkg/errors"
	awsrestapidomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapidomain/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the custom domain, its base-path
// mappings, and any private access associations, and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsrestapidomainv1alpha1.AwsRestApiDomainStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := domain(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "rest api domain")
	}

	return nil
}
