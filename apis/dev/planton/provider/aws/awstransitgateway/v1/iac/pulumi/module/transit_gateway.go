package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2transitgateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// transitGateway creates the AWS Transit Gateway -- the pure hub. Everything
// that composes onto it (VPC attachments, custom route tables) is its own
// resource kind referencing this gateway's outputs, so spokes and routing
// domains come and go without touching the hub.
func transitGateway(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ec2transitgateway.TransitGateway, error) {
	spec := locals.TransitGateway.Spec

	args := &ec2transitgateway.TransitGatewayArgs{
		// Tri-state dials: nil falls through to the provider default
		// ("enable" for DNS/ECMP and the default-table pair). AWS QUIRK on
		// the default-table pair: flipping disable -> enable REPLACES the
		// gateway (the provider's asymmetric ForceNew), while
		// enable -> disable updates in place -- state a topology's posture
		// up front rather than tightening later.
		DefaultRouteTableAssociation: enableDisableTriState(spec.DefaultRouteTableAssociation),
		DefaultRouteTablePropagation: enableDisableTriState(spec.DefaultRouteTablePropagation),
		DnsSupport:                   enableDisableTriState(spec.DnsSupport),
		VpnEcmpSupport:               enableDisableTriState(spec.VpnEcmpSupport),

		// Left nil unless the spec pins it: AWS computes the effective
		// in-transit encryption posture on its own when unset.
		EncryptionSupport: enableDisableTriState(spec.EncryptionSupport),

		// Plain-bool dials whose spec default (false) matches the AWS
		// default ("disable"): always sent explicitly, so the applied state
		// is exactly the spec with no inheritance ambiguity.
		// multicast_support is create-time immutable -- changing it
		// replaces the gateway.
		AutoAcceptSharedAttachments:     pulumi.StringPtr(enableDisable(spec.AutoAcceptSharedAttachments)),
		SecurityGroupReferencingSupport: pulumi.StringPtr(enableDisable(spec.SecurityGroupReferencingSupport)),
		MulticastSupport:                pulumi.StringPtr(enableDisable(spec.MulticastSupport)),

		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// 0 means "not set" (proto3 scalar zero value): fall through to the
	// provider default of 64512 rather than sending an out-of-range ASN.
	// Changing the ASN after creation replaces the gateway.
	if spec.AmazonSideAsn != 0 {
		args.AmazonSideAsn = pulumi.IntPtr(int(spec.AmazonSideAsn))
	}

	// Only needed for TGW Connect (GRE appliance integration); empty for
	// the overwhelming majority of hubs.
	if len(spec.TransitGatewayCidrBlocks) > 0 {
		args.TransitGatewayCidrBlocks = pulumi.ToStringArray(spec.TransitGatewayCidrBlocks)
	}

	// The cloud name is carried by the Name tag (set in locals); the Pulumi
	// resource name below is only the URN segment.
	createdTransitGateway, err := ec2transitgateway.NewTransitGateway(
		ctx,
		locals.TransitGateway.Metadata.Name,
		args,
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transit gateway")
	}

	return createdTransitGateway, nil
}
