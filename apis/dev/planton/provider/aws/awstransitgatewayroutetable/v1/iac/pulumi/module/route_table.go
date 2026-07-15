package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2transitgateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// routeTable creates the Transit Gateway route table and materializes its
// folded routing domain: associations, propagations, static routes, and
// prefix list references. Every folded member is its own provider resource
// keyed by a value that is unique within the table (attachment ID,
// destination CIDR, prefix list ID -- the spec CELs enforce uniqueness), so
// adding or removing one member never churns its neighbors.
func routeTable(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ec2transitgateway.RouteTable, error) {
	spec := locals.RouteTable.Spec

	// The table itself is a tiny resource; the routing domain lives in the
	// folded members below. Moving the table to another gateway replaces it
	// (and with it, everything folded).
	createdRouteTable, err := ec2transitgateway.NewRouteTable(
		ctx,
		locals.RouteTable.Metadata.Name,
		&ec2transitgateway.RouteTableArgs{
			TransitGatewayId: pulumi.String(spec.TransitGatewayId.GetValue()),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transit gateway route table")
	}

	// Associations: the attachments whose OUTBOUND traffic is looked up in
	// this table. AWS allows an attachment at most ONE association across
	// the whole gateway -- an attachment listed here must have its default
	// association turned off and must not appear in another table's
	// associations. Documents cannot see each other, so that gateway-wide
	// uniqueness is enforced by AWS at apply time, not by validation.
	for _, association := range spec.Associations {
		attachmentId := association.GetValue()
		_, err := ec2transitgateway.NewRouteTableAssociation(
			ctx,
			fmt.Sprintf("%s-association-%s", locals.RouteTable.Metadata.Name, attachmentId),
			&ec2transitgateway.RouteTableAssociationArgs{
				TransitGatewayAttachmentId: pulumi.String(attachmentId),
				TransitGatewayRouteTableId: createdRouteTable.ID(),
			},
			pulumi.Provider(provider),
			pulumi.Parent(createdRouteTable),
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to associate attachment %s", attachmentId)
		}
	}

	// Propagations: the attachments that ADVERTISE their routes into this
	// table. VPC attachments propagate their VPC CIDRs; VPN and Direct
	// Connect attachments propagate BGP-learned routes. Unlike
	// associations, an attachment can propagate to any number of tables.
	for _, propagation := range spec.Propagations {
		attachmentId := propagation.GetValue()
		_, err := ec2transitgateway.NewRouteTablePropagation(
			ctx,
			fmt.Sprintf("%s-propagation-%s", locals.RouteTable.Metadata.Name, attachmentId),
			&ec2transitgateway.RouteTablePropagationArgs{
				TransitGatewayAttachmentId: pulumi.String(attachmentId),
				TransitGatewayRouteTableId: createdRouteTable.ID(),
			},
			pulumi.Provider(provider),
			pulumi.Parent(createdRouteTable),
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to propagate attachment %s", attachmentId)
		}
	}

	// Static routes: longest-prefix match beats propagated routes on ties,
	// so these steer specific prefixes (an inspection detour, a default
	// route toward an egress VPC) or blackhole traffic that must never
	// cross the hub. The spec CEL guarantees exactly one of
	// attachment_id/blackhole per route; AWS expects the attachment to be
	// ABSENT (not empty) for a blackhole.
	for _, route := range spec.Routes {
		args := &ec2transitgateway.RouteArgs{
			DestinationCidrBlock:       pulumi.String(route.DestinationCidrBlock),
			TransitGatewayRouteTableId: createdRouteTable.ID(),
		}
		if route.Blackhole {
			args.Blackhole = pulumi.BoolPtr(true)
		} else {
			args.TransitGatewayAttachmentId = pulumi.StringPtr(route.AttachmentId.GetValue())
		}
		_, err := ec2transitgateway.NewRoute(
			ctx,
			fmt.Sprintf("%s-route-%s", locals.RouteTable.Metadata.Name, route.DestinationCidrBlock),
			args,
			pulumi.Provider(provider),
			pulumi.Parent(createdRouteTable),
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create route %s", route.DestinationCidrBlock)
		}
	}

	// Prefix list references: route a managed prefix list's whole CIDR set
	// via one attachment (or blackhole it), tracking the list's membership
	// as it changes -- one reference instead of N hand-maintained statics.
	for _, prefixListReference := range spec.PrefixListReferences {
		args := &ec2transitgateway.PrefixListReferenceArgs{
			PrefixListId:               pulumi.String(prefixListReference.PrefixListId),
			TransitGatewayRouteTableId: createdRouteTable.ID(),
		}
		if prefixListReference.Blackhole {
			args.Blackhole = pulumi.BoolPtr(true)
		} else {
			args.TransitGatewayAttachmentId = pulumi.StringPtr(prefixListReference.AttachmentId.GetValue())
		}
		_, err := ec2transitgateway.NewPrefixListReference(
			ctx,
			fmt.Sprintf("%s-plref-%s", locals.RouteTable.Metadata.Name, prefixListReference.PrefixListId),
			args,
			pulumi.Provider(provider),
			pulumi.Parent(createdRouteTable),
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create prefix list reference %s", prefixListReference.PrefixListId)
		}
	}

	return createdRouteTable, nil
}
