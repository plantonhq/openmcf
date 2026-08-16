package module

import (
	"github.com/pkg/errors"
	awsebsvolumev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsebsvolume/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the volume (either arm), its
// attachments, and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsebsvolumev1alpha1.AwsEbsVolumeStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := volume(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "volume")
	}

	return nil
}
