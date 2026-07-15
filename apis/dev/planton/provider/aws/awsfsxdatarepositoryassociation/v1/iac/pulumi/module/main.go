package module

import (
	"github.com/pkg/errors"
	awsfsxdrav1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsfsxdatarepositoryassociation/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the primary entry point for the AwsFsxDataRepositoryAssociation
// Pulumi module. It creates the association and exports its outputs for
// downstream consumption.
func Resources(ctx *pulumi.Context, stackInput *awsfsxdrav1.AwsFsxDataRepositoryAssociationStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsFsxDataRepositoryAssociation.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdAssociation, err := dataRepositoryAssociation(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create data repository association")
	}

	ctx.Export(OpAssociationId, createdAssociation.AssociationId)
	ctx.Export(OpAssociationArn, createdAssociation.Arn)
	ctx.Export(OpFileSystemId, createdAssociation.FileSystemId)

	return nil
}
