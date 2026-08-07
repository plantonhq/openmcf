package module

import (
	"github.com/pkg/errors"
	awssagemakerdomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerdomain/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of AWS SageMaker Domain related resources and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awssagemakerdomainv1alpha1.AwsSagemakerDomainStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	// SageMaker Domain
	createdDomain, err := domain(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "sagemaker domain")
	}

	// Export outputs
	outputs(ctx, createdDomain)

	return nil
}
