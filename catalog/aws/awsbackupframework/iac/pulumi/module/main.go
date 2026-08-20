package module

import (
	"github.com/pkg/errors"
	awsbackupframeworkv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbackupframework/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of the Backup Audit Manager
// framework and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsbackupframeworkv1alpha1.AwsBackupFrameworkStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.Target.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	if err := framework(ctx, locals, provider); err != nil {
		return errors.Wrap(err, "framework")
	}

	return nil
}
