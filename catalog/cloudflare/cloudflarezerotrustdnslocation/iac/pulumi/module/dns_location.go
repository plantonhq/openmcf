package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustdnslocationv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustdnslocation/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dnsLocation creates the Gateway DNS location. A plain CRUD resource (real
// create/update/delete; only the account forces replacement). Provider
// behaviors the mapping honors:
//   - Update is a full replace at the API: the spec's max_ttl is sent
//     whenever declared, and its omission genuinely resets the behavior to
//     inherit (documented on the spec field).
//   - dns_destination_ips_id is only sent when set -- unset lets Cloudflare
//     auto-assign the shared IPv4 destination pair.
//   - max_ttl and every networks list are ALWAYS sent as known values
//     (max_ttl coalesces to the documented server default {mode: inherit};
//     an absent networks list is sent empty). At provider v5.23.0/v5.24.0
//     the Go model types these computed-optional attributes as raw pointers
//     that cannot hold "unknown", so leaving any of them null plans it as
//     unknown and CRASHES the apply-time conversion ("Value Conversion
//     Error ... target type cannot handle unknown values"; measured live
//     2026-08-26 through this bridged provider too, unfixed on provider
//     main). The explicit server defaults keep the planned value known and
//     change nothing semantically.
func dnsLocation(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustDnsLocation.Spec

	args := &cloudflare.ZeroTrustDnsLocationArgs{
		AccountId: pulumi.String(spec.AccountId),
		Name:      pulumi.String(spec.Name),
	}

	if spec.ClientDefault != nil {
		args.ClientDefault = pulumi.BoolPtr(spec.GetClientDefault())
	}
	if spec.EcsSupport != nil {
		args.EcsSupport = pulumi.BoolPtr(spec.GetEcsSupport())
	}
	if spec.DnsDestinationIpsId != "" {
		args.DnsDestinationIpsId = pulumi.String(spec.DnsDestinationIpsId)
	}

	if spec.Endpoints != nil {
		args.Endpoints = buildEndpoints(spec.Endpoints)
	}

	networks := cloudflare.ZeroTrustDnsLocationNetworkArray{}
	for _, network := range spec.Networks {
		networks = append(networks, cloudflare.ZeroTrustDnsLocationNetworkArgs{
			Network: pulumi.String(network.Network),
		})
	}
	args.Networks = networks

	maxTtl := cloudflare.ZeroTrustDnsLocationMaxTtlArgs{
		Mode: pulumi.String("inherit"),
	}
	if spec.MaxTtl != nil {
		maxTtl.Mode = pulumi.String(spec.MaxTtl.Mode)
		if spec.MaxTtl.TtlSecs != nil {
			maxTtl.TtlSecs = pulumi.IntPtr(int(spec.MaxTtl.GetTtlSecs()))
		}
	}
	args.MaxTtl = maxTtl

	createdLocation, err := cloudflare.NewZeroTrustDnsLocation(
		ctx,
		"dns_location",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create dns location")
	}

	ctx.Export(OpLocationId, createdLocation.ID())
	ctx.Export(OpDohSubdomain, createdLocation.DohSubdomain)
	ctx.Export(OpIp, createdLocation.Ip)

	return nil
}

// buildEndpoints maps the all-four-types endpoint tree (the spec requires
// every type declared when the tree is set, mirroring the provider schema).
func buildEndpoints(endpoints *cloudflarezerotrustdnslocationv1alpha1.CloudflareZeroTrustDnsLocationEndpoints) cloudflare.ZeroTrustDnsLocationEndpointsPtrInput {
	networkRows := func(rows []*cloudflarezerotrustdnslocationv1alpha1.CloudflareZeroTrustDnsLocationNetwork) []string {
		networks := make([]string, 0, len(rows))
		for _, row := range rows {
			networks = append(networks, row.Network)
		}
		return networks
	}

	doh := cloudflare.ZeroTrustDnsLocationEndpointsDohArgs{}
	if endpoints.Doh.Enabled != nil {
		doh.Enabled = pulumi.BoolPtr(endpoints.Doh.GetEnabled())
	}
	if endpoints.Doh.RequireToken != nil {
		doh.RequireToken = pulumi.BoolPtr(endpoints.Doh.GetRequireToken())
	}
	dohNetworks := cloudflare.ZeroTrustDnsLocationEndpointsDohNetworkArray{}
	for _, network := range networkRows(endpoints.Doh.Networks) {
		dohNetworks = append(dohNetworks, cloudflare.ZeroTrustDnsLocationEndpointsDohNetworkArgs{
			Network: pulumi.String(network),
		})
	}
	doh.Networks = dohNetworks

	dot := cloudflare.ZeroTrustDnsLocationEndpointsDotArgs{}
	if endpoints.Dot.Enabled != nil {
		dot.Enabled = pulumi.BoolPtr(endpoints.Dot.GetEnabled())
	}
	dotNetworks := cloudflare.ZeroTrustDnsLocationEndpointsDotNetworkArray{}
	for _, network := range networkRows(endpoints.Dot.Networks) {
		dotNetworks = append(dotNetworks, cloudflare.ZeroTrustDnsLocationEndpointsDotNetworkArgs{
			Network: pulumi.String(network),
		})
	}
	dot.Networks = dotNetworks

	ipv4 := cloudflare.ZeroTrustDnsLocationEndpointsIpv4Args{}
	if endpoints.Ipv4.Enabled != nil {
		ipv4.Enabled = pulumi.BoolPtr(endpoints.Ipv4.GetEnabled())
	}

	ipv6 := cloudflare.ZeroTrustDnsLocationEndpointsIpv6Args{}
	if endpoints.Ipv6.Enabled != nil {
		ipv6.Enabled = pulumi.BoolPtr(endpoints.Ipv6.GetEnabled())
	}
	ipv6Networks := cloudflare.ZeroTrustDnsLocationEndpointsIpv6NetworkArray{}
	for _, network := range networkRows(endpoints.Ipv6.Networks) {
		ipv6Networks = append(ipv6Networks, cloudflare.ZeroTrustDnsLocationEndpointsIpv6NetworkArgs{
			Network: pulumi.String(network),
		})
	}
	ipv6.Networks = ipv6Networks

	return cloudflare.ZeroTrustDnsLocationEndpointsArgs{
		Doh:  doh,
		Dot:  dot,
		Ipv4: ipv4,
		Ipv6: ipv6,
	}
}
