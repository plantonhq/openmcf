package module

import (
	"github.com/pkg/errors"
	awstgwrtv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awstransitgatewayroutetable/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the primary entry point for the AwsTransitGatewayRouteTable
// Pulumi module. It creates the route table with its folded routing domain
// (associations, propagations, static routes, prefix list references) and
// exports the table's outputs.
func Resources(ctx *pulumi.Context, stackInput *awstgwrtv1.AwsTransitGatewayRouteTableStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.RouteTable.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdRouteTable, err := routeTable(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create transit gateway route table")
	}

	ctx.Export(OpRouteTableId, createdRouteTable.ID())
	ctx.Export(OpRouteTableArn, createdRouteTable.Arn)

	return nil
}
