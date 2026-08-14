package module

import (
	"github.com/pkg/errors"
	awssagemakermodelv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakermodel/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the SageMaker model and exports
// outputs.
func Resources(ctx *pulumi.Context, stackInput *awssagemakermodelv1alpha1.AwsSagemakerModelStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := model(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "sagemaker model")
	}

	return nil
}
