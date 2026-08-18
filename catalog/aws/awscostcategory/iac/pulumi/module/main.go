package module

import (
	"github.com/pkg/errors"
	awscostcategoryv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscostcategory/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the cost category and exports
// outputs.
func Resources(ctx *pulumi.Context, stackInput *awscostcategoryv1alpha1.AwsCostCategoryStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := costCategory(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "cost category")
	}

	return nil
}
