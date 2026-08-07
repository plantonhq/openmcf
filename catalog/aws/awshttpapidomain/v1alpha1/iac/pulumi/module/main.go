package module

import (
	"github.com/pkg/errors"
	awshttpapidomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awshttpapidomain/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates custom domain + API mapping creation and exports
// outputs.
func Resources(ctx *pulumi.Context, stackInput *awshttpapidomainv1alpha1.AwsHttpApiDomainStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsHttpApiDomain.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := domainName(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "domain name")
	}

	return nil
}
