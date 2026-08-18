package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// peering renders whichever arm the spec configured (the spec's CEL
// guarantees exactly one) and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - VpcId / PeerVpcId / PeerOwnerId / PeerRegion are fixed for life
//     (replace-on-change);
//   - AutoAccept works same-account, same-region only (the provider
//     hard-errors on PeerRegion + AutoAccept; the spec's CELs
//     front-load both walls);
//   - DNS-resolution options are managed in-line as the single owner
//     (the standalone options resource fights this form) and need an
//     ACTIVE connection;
//   - the accept arm's destroy is a NO-OP at AWS - it abandons
//     management without deleting the peering (only the requester
//     side deletes);
//   - one VPC pair supports at most one peering - AWS returns the
//     EXISTING connection id for a duplicate request.
func peering(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	if spec.Request != nil {
		return requestArm(ctx, locals, provider)
	}
	return acceptArm(ctx, locals, provider)
}

func requestArm(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	request := locals.Spec.Request

	args := &ec2.VpcPeeringConnectionArgs{
		VpcId:     pulumi.String(request.VpcId.GetValue()),
		PeerVpcId: pulumi.String(request.PeerVpcId.GetValue()),
		Tags:      pulumi.ToStringMap(locals.AwsTags),
	}
	if request.PeerOwnerId != "" {
		args.PeerOwnerId = pulumi.String(request.PeerOwnerId)
	}
	if request.PeerRegion != "" {
		args.PeerRegion = pulumi.String(request.PeerRegion)
	}
	if request.AutoAccept {
		args.AutoAccept = pulumi.Bool(true)
	}
	if request.RequesterAllowRemoteVpcDnsResolution {
		args.Requester = &ec2.VpcPeeringConnectionRequesterArgs{
			AllowRemoteVpcDnsResolution: pulumi.Bool(true),
		}
	}
	if request.AccepterAllowRemoteVpcDnsResolution {
		args.Accepter = &ec2.VpcPeeringConnectionAccepterTypeArgs{
			AllowRemoteVpcDnsResolution: pulumi.Bool(true),
		}
	}

	createdConnection, err := ec2.NewVpcPeeringConnection(ctx, "peering_connection", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create peering connection (request arm)")
	}

	ctx.Export(OpPeeringConnectionId, createdConnection.ID())
	ctx.Export(OpAcceptStatus, createdConnection.AcceptStatus)
	return nil
}

func acceptArm(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	accept := locals.Spec.Accept

	args := &ec2.VpcPeeringConnectionAccepterArgs{
		VpcPeeringConnectionId: pulumi.String(accept.VpcPeeringConnectionId.GetValue()),
		Tags:                   pulumi.ToStringMap(locals.AwsTags),
	}
	if accept.AutoAccept {
		args.AutoAccept = pulumi.Bool(true)
	}
	if accept.AccepterAllowRemoteVpcDnsResolution {
		args.Accepter = &ec2.VpcPeeringConnectionAccepterAccepterArgs{
			AllowRemoteVpcDnsResolution: pulumi.Bool(true),
		}
	}

	createdAccepter, err := ec2.NewVpcPeeringConnectionAccepter(ctx, "peering_accepter", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "accept peering connection (accept arm)")
	}

	ctx.Export(OpPeeringConnectionId, createdAccepter.ID())
	ctx.Export(OpAcceptStatus, createdAccepter.AcceptStatus)
	return nil
}
