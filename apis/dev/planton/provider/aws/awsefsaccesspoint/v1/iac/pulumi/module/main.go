package module

import (
	"github.com/pkg/errors"
	awsefsaccesspointv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsefsaccesspoint/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *awsefsaccesspointv1.AwsEfsAccessPointStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsEfsAccessPoint.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	result, err := accessPoint(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create access point")
	}

	// --- Exports ---
	ctx.Export(OpAccessPointId, result.AccessPoint.ID())
	ctx.Export(OpAccessPointArn, result.AccessPoint.Arn)
	ctx.Export(OpFileSystemId, result.AccessPoint.FileSystemId)
	ctx.Export(OpFileSystemArn, result.AccessPoint.FileSystemArn)

	return nil
}
