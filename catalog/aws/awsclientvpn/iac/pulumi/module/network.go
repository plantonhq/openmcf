package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2clientvpn"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// networkAssociations attaches the endpoint to each subnet in
// spec.subnet_ids. Associations are the slow part of a Client VPN deploy
// (AWS takes minutes per attach/detach), keyed per subnet so adding one
// never disturbs the others. The map of created associations is returned so
// routes can depend on THEIR subnet's association -- AWS rejects a route
// whose target subnet is not yet associated.
func networkAssociations(
	ctx *pulumi.Context,
	locals *Locals,
	provider pulumi.ProviderResource,
	createdEndpoint *ec2clientvpn.Endpoint,
) (map[string]*ec2clientvpn.NetworkAssociation, error) {
	spec := locals.AwsClientVpn.Spec

	createdAssociations := make(map[string]*ec2clientvpn.NetworkAssociation, len(spec.SubnetIds))
	for _, subnetRef := range spec.SubnetIds {
		subnetId := subnetRef.GetValue()
		if subnetId == "" {
			continue
		}

		createdAssociation, err := ec2clientvpn.NewNetworkAssociation(ctx,
			fmt.Sprintf("association-%s", subnetId),
			&ec2clientvpn.NetworkAssociationArgs{
				ClientVpnEndpointId: createdEndpoint.ID(),
				SubnetId:            pulumi.String(subnetId),
			}, pulumi.Provider(provider), pulumi.Parent(createdEndpoint))
		if err != nil {
			return nil, errors.Wrapf(err, "associate subnet %s", subnetId)
		}
		createdAssociations[subnetId] = createdAssociation

		// Dot-suffixed exports flatten into the subnet_association_ids map
		// output, keyed by subnet ID verbatim.
		ctx.Export(fmt.Sprintf("%s.%s", OpSubnetAssociationIds, subnetId), createdAssociation.AssociationId)
	}

	return createdAssociations, nil
}

// routes creates the endpoint's additional route-table entries. Each route
// depends on its target subnet's association: AWS returns
// InvalidClientVpnActiveAssociationNotFound for a route whose subnet is
// still associating, so the edge makes ordering explicit instead of leaning
// on provider-side retries.
func routes(
	ctx *pulumi.Context,
	locals *Locals,
	provider pulumi.ProviderResource,
	createdEndpoint *ec2clientvpn.Endpoint,
	createdAssociations map[string]*ec2clientvpn.NetworkAssociation,
) error {
	spec := locals.AwsClientVpn.Spec

	for _, route := range spec.Routes {
		subnetId := route.TargetSubnetId.GetValue()

		args := &ec2clientvpn.RouteArgs{
			ClientVpnEndpointId:  createdEndpoint.ID(),
			DestinationCidrBlock: pulumi.String(route.DestinationCidrBlock),
			TargetVpcSubnetId:    pulumi.String(subnetId),
		}
		if route.Description != "" {
			args.Description = pulumi.String(route.Description)
		}

		// AWS deletes a route through the endpoint's associated target
		// network slowly -- observed live 'deleting' past the provider's
		// 4-minute default delete timeout (twice, 2026-08-13); the delete
		// does complete, just late. Both timeouts are pinned above the
		// observed worst case (mirrors the Terraform module's timeouts).
		opts := []pulumi.ResourceOption{
			pulumi.Provider(provider),
			pulumi.Parent(createdEndpoint),
			pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "10m", Delete: "10m"}),
		}
		if createdAssociation, ok := createdAssociations[subnetId]; ok {
			opts = append(opts, pulumi.DependsOn([]pulumi.Resource{createdAssociation}))
		}

		// Routes are keyed by (destination, subnet) -- the same pair AWS
		// uses -- so a destination can be reached through several subnets.
		if _, err := ec2clientvpn.NewRoute(ctx,
			fmt.Sprintf("route-%s-%s", route.DestinationCidrBlock, subnetId),
			args, opts...); err != nil {
			return errors.Wrapf(err, "route %s via %s", route.DestinationCidrBlock, subnetId)
		}
	}

	return nil
}
