package module

import (
	"github.com/pkg/errors"
	awstgwv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awstransitgateway/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the primary entry point for the AwsTransitGateway Pulumi
// module. It creates the Transit Gateway hub and exports its outputs for
// downstream consumption; VPC attachments and route tables are their own
// resource kinds composing onto the exported gateway ID.
func Resources(ctx *pulumi.Context, stackInput *awstgwv1.AwsTransitGatewayStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.TransitGateway.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdTransitGateway, err := transitGateway(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create transit gateway")
	}

	ctx.Export(OpTransitGatewayId, createdTransitGateway.ID())
	ctx.Export(OpTransitGatewayArn, createdTransitGateway.Arn)
	ctx.Export(OpOwnerId, createdTransitGateway.OwnerId)
	ctx.Export(OpAssociationDefaultRouteTableId, createdTransitGateway.AssociationDefaultRouteTableId)
	ctx.Export(OpPropagationDefaultRouteTableId, createdTransitGateway.PropagationDefaultRouteTableId)

	return nil
}
