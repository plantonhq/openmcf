package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustaccessinfrastructuretargetv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustaccessinfrastructuretarget/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// target registers the infrastructure target: the hostname/IP inventory row
// that Access infrastructure applications select and broker SSH access to.
// A plain CRUD resource: real create, in-place updates (hostname, addresses,
// virtual networks), real delete. Only the account forces replacement.
func target(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustAccessInfrastructureTarget.Spec

	args := &cloudflare.ZeroTrustAccessInfrastructureTargetArgs{
		AccountId: pulumi.String(spec.AccountId),
		Hostname:  pulumi.String(spec.Hostname),
		Ip:        buildIp(spec.Ip),
	}

	createdTarget, err := cloudflare.NewZeroTrustAccessInfrastructureTarget(
		ctx,
		"target",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create infrastructure target")
	}

	ctx.Export(OpTargetId, createdTarget.ID())

	return nil
}

// buildIp maps the spec's per-family address blocks. An omitted
// virtual_network_id is not sent -- Cloudflare then assigns the account's
// default virtual network (the field is computed on the provider side, so
// the assigned value never drifts).
func buildIp(ip *cloudflarezerotrustaccessinfrastructuretargetv1alpha1.CloudflareZeroTrustAccessInfrastructureTargetIp) cloudflare.ZeroTrustAccessInfrastructureTargetIpArgs {
	ipArgs := cloudflare.ZeroTrustAccessInfrastructureTargetIpArgs{}

	if ip.Ipv4 != nil {
		ipv4 := cloudflare.ZeroTrustAccessInfrastructureTargetIpIpv4Args{
			IpAddr: pulumi.String(ip.Ipv4.IpAddr),
		}
		if ip.Ipv4.VirtualNetworkId.GetValue() != "" {
			ipv4.VirtualNetworkId = pulumi.String(ip.Ipv4.VirtualNetworkId.GetValue())
		}
		ipArgs.Ipv4 = ipv4
	}

	if ip.Ipv6 != nil {
		ipv6 := cloudflare.ZeroTrustAccessInfrastructureTargetIpIpv6Args{
			IpAddr: pulumi.String(ip.Ipv6.IpAddr),
		}
		if ip.Ipv6.VirtualNetworkId.GetValue() != "" {
			ipv6.VirtualNetworkId = pulumi.String(ip.Ipv6.VirtualNetworkId.GetValue())
		}
		ipArgs.Ipv6 = ipv6
	}

	return ipArgs
}
