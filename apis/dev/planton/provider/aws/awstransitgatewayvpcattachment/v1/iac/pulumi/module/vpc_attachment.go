package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2transitgateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vpcAttachment creates the Transit Gateway VPC attachment -- one spoke of
// the hub. AWS provisions an ENI in each listed subnet; traffic between the
// VPC and the gateway flows through those ENIs. The gateway and the VPC are
// create-time immutable (changing either replaces the attachment), while
// the subnet set updates in place.
func vpcAttachment(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ec2transitgateway.VpcAttachment, error) {
	spec := locals.VpcAttachment.Spec

	// Referenced IDs arrive pre-resolved as plain strings.
	subnetIds := make([]string, 0, len(spec.SubnetIds))
	for _, subnetId := range spec.SubnetIds {
		subnetIds = append(subnetIds, subnetId.GetValue())
	}

	args := &ec2transitgateway.VpcAttachmentArgs{
		TransitGatewayId: pulumi.String(spec.TransitGatewayId.GetValue()),
		VpcId:            pulumi.String(spec.VpcId.GetValue()),
		SubnetIds:        pulumi.ToStringArray(subnetIds),

		// Tri-state options: nil falls through to the provider default
		// ("enable" for DNS support) or, for security group referencing, to
		// the value the attachment inherits from the gateway (AWS computes
		// it).
		DnsSupport:                      enableDisableTriState(spec.DnsSupport),
		SecurityGroupReferencingSupport: enableDisableTriState(spec.SecurityGroupReferencingSupport),

		// Plain-bool options whose spec default (false) matches the AWS
		// default ("disable"): always sent explicitly. Appliance mode keeps
		// a flow's return traffic in the AZ it entered through -- required
		// for stateful inspection VPCs.
		Ipv6Support:          pulumi.StringPtr(enableDisable(spec.Ipv6Support)),
		ApplianceModeSupport: pulumi.StringPtr(enableDisable(spec.ApplianceModeSupport)),

		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Default route table membership. Left nil, AWS derives the value from
	// the GATEWAY's own default-association/propagation dials -- the
	// provider marks both Optional+Computed for exactly this inheritance.
	// Only a spec that pins them sends a value: false detaches this spoke
	// from the default table (the segmented-topology posture, where a
	// custom AwsTransitGatewayRouteTable owns the association instead).
	if spec.DefaultRouteTableAssociation != nil {
		args.TransitGatewayDefaultRouteTableAssociation = pulumi.BoolPtr(*spec.DefaultRouteTableAssociation)
	}
	if spec.DefaultRouteTablePropagation != nil {
		args.TransitGatewayDefaultRouteTablePropagation = pulumi.BoolPtr(*spec.DefaultRouteTablePropagation)
	}

	// The cloud name is carried by the Name tag (set in locals); the Pulumi
	// resource name below is only the URN segment.
	createdVpcAttachment, err := ec2transitgateway.NewVpcAttachment(
		ctx,
		locals.VpcAttachment.Metadata.Name,
		args,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transit gateway vpc attachment")
	}

	return createdVpcAttachment, nil
}
