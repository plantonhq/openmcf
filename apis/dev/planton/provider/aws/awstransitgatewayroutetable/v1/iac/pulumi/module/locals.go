package module

import (
	"strconv"

	awstgwrtv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awstransitgatewayroutetable/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the route table.
type Locals struct {
	RouteTable *awstgwrtv1.AwsTransitGatewayRouteTable
	AwsTags    map[string]string
}

// initializeLocals reads the stack input and builds the Locals instance.
func initializeLocals(ctx *pulumi.Context, stackInput *awstgwrtv1.AwsTransitGatewayRouteTableStackInput) *Locals {
	locals := &Locals{}

	locals.RouteTable = stackInput.Target

	// Identity tags match the Terraform module key-for-key. The Name tag IS
	// the route table's console identity -- route tables have no name
	// attribute of their own.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.RouteTable.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.RouteTable.Metadata.Org,
		awstagkeys.Environment:  locals.RouteTable.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsTransitGatewayRouteTable.String(),
		awstagkeys.ResourceId:   locals.RouteTable.Metadata.Id,
	}

	return locals
}
