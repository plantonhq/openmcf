package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// customProfile creates the targeted WARP device profile: the default
// profile's settings body, applied only to the devices matched by the
// wirefilter expression. A real object -- create, update, and delete all do
// what they say (deleting the profile returns its devices to the default
// profile).
func customProfile(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustDeviceCustomProfile.Spec

	args := &cloudflare.ZeroTrustDeviceCustomProfileArgs{
		AccountId:  pulumi.String(spec.AccountId),
		Name:       pulumi.String(spec.Name),
		Match:      pulumi.String(spec.Match),
		Precedence: pulumi.Float64Ptr(float64(spec.Precedence)),
	}

	if spec.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.AllowModeSwitch != nil {
		args.AllowModeSwitch = pulumi.BoolPtr(spec.GetAllowModeSwitch())
	}
	if spec.AllowUpdates != nil {
		args.AllowUpdates = pulumi.BoolPtr(spec.GetAllowUpdates())
	}
	if spec.AllowedToLeave != nil {
		args.AllowedToLeave = pulumi.BoolPtr(spec.GetAllowedToLeave())
	}
	if spec.AutoConnect != nil {
		args.AutoConnect = pulumi.Float64Ptr(float64(spec.GetAutoConnect()))
	}
	if spec.CaptivePortal != nil {
		args.CaptivePortal = pulumi.Float64Ptr(float64(spec.GetCaptivePortal()))
	}
	if spec.DisableAutoFallback != nil {
		args.DisableAutoFallback = pulumi.BoolPtr(spec.GetDisableAutoFallback())
	}
	if spec.ExcludeOfficeIps != nil {
		args.ExcludeOfficeIps = pulumi.BoolPtr(spec.GetExcludeOfficeIps())
	}
	if spec.RegisterInterfaceIpWithDns != nil {
		args.RegisterInterfaceIpWithDns = pulumi.BoolPtr(spec.GetRegisterInterfaceIpWithDns())
	}
	if spec.SccmVpnBoundarySupport != nil {
		args.SccmVpnBoundarySupport = pulumi.BoolPtr(spec.GetSccmVpnBoundarySupport())
	}
	if spec.SwitchLocked != nil {
		args.SwitchLocked = pulumi.BoolPtr(spec.GetSwitchLocked())
	}
	if spec.SupportUrl != "" {
		args.SupportUrl = pulumi.String(spec.SupportUrl)
	}
	if spec.TunnelProtocol != "" {
		args.TunnelProtocol = pulumi.String(spec.TunnelProtocol)
	}
	if spec.LanAllowMinutes != nil {
		args.LanAllowMinutes = pulumi.Float64Ptr(float64(spec.GetLanAllowMinutes()))
	}
	if spec.LanAllowSubnetSize != nil {
		args.LanAllowSubnetSize = pulumi.Float64Ptr(float64(spec.GetLanAllowSubnetSize()))
	}

	// Split tunnel: exclude and include are mutually exclusive (spec
	// validation enforces it). Each entry carries exactly one of address or
	// host; the unset one is never sent.
	if len(spec.Exclude) > 0 {
		excludes := cloudflare.ZeroTrustDeviceCustomProfileExcludeArray{}
		for _, entry := range spec.Exclude {
			row := cloudflare.ZeroTrustDeviceCustomProfileExcludeArgs{}
			if entry.Address != "" {
				row.Address = pulumi.String(entry.Address)
			}
			if entry.Host != "" {
				row.Host = pulumi.String(entry.Host)
			}
			if entry.Description != "" {
				row.Description = pulumi.String(entry.Description)
			}
			excludes = append(excludes, row)
		}
		args.Excludes = excludes
	}

	if len(spec.Include) > 0 {
		includes := cloudflare.ZeroTrustDeviceCustomProfileIncludeArray{}
		for _, entry := range spec.Include {
			row := cloudflare.ZeroTrustDeviceCustomProfileIncludeArgs{}
			if entry.Address != "" {
				row.Address = pulumi.String(entry.Address)
			}
			if entry.Host != "" {
				row.Host = pulumi.String(entry.Host)
			}
			if entry.Description != "" {
				row.Description = pulumi.String(entry.Description)
			}
			includes = append(includes, row)
		}
		args.Includes = includes
	}

	if spec.ServiceModeV2 != nil {
		serviceMode := cloudflare.ZeroTrustDeviceCustomProfileServiceModeV2Args{
			Mode: pulumi.StringPtr(spec.ServiceModeV2.Mode),
		}
		if spec.ServiceModeV2.Port != nil {
			serviceMode.Port = pulumi.Float64Ptr(float64(spec.ServiceModeV2.GetPort()))
		}
		args.ServiceModeV2 = serviceMode
	}

	if spec.VirtualNetworks != nil {
		allowed := make([]string, 0, len(spec.VirtualNetworks.Allowed))
		for _, vnet := range spec.VirtualNetworks.Allowed {
			allowed = append(allowed, vnet.GetValue())
		}
		args.VirtualNetworks = cloudflare.ZeroTrustDeviceCustomProfileVirtualNetworksArgs{
			Alloweds: pulumi.ToStringArray(allowed),
			Default:  pulumi.String(spec.VirtualNetworks.DefaultVirtualNetworkId.GetValue()),
		}
	}

	if len(spec.DnsSearchSuffixes) > 0 {
		suffixes := cloudflare.ZeroTrustDeviceCustomProfileDnsSearchSuffixArray{}
		for _, row := range spec.DnsSearchSuffixes {
			suffix := cloudflare.ZeroTrustDeviceCustomProfileDnsSearchSuffixArgs{
				Suffix: pulumi.String(row.Suffix),
			}
			if row.Description != "" {
				suffix.Description = pulumi.String(row.Description)
			}
			suffixes = append(suffixes, suffix)
		}
		args.DnsSearchSuffixes = suffixes
	}

	createdProfile, err := cloudflare.NewZeroTrustDeviceCustomProfile(
		ctx,
		"custom_profile",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create custom device profile")
	}

	// This profile's local-DNS fallback list. The profile resource reports
	// fallback_domains READ-ONLY; this dedicated per-profile companion is
	// the only write path, and it replaces the whole list on every apply
	// (declarative: what the spec lists is exactly what exists). The rows
	// ride the profile -- deleting the profile retires them with it.
	// Deployed only when the spec declares rows.
	if len(spec.FallbackDomains) > 0 {
		domains := cloudflare.ZeroTrustDeviceCustomProfileLocalDomainFallbackDomainArray{}
		for _, row := range spec.FallbackDomains {
			domain := cloudflare.ZeroTrustDeviceCustomProfileLocalDomainFallbackDomainArgs{
				Suffix: pulumi.String(row.Suffix),
			}
			if row.Description != "" {
				domain.Description = pulumi.String(row.Description)
			}
			if len(row.DnsServer) > 0 {
				domain.DnsServers = pulumi.ToStringArray(row.DnsServer)
			}
			domains = append(domains, domain)
		}

		_, err := cloudflare.NewZeroTrustDeviceCustomProfileLocalDomainFallback(
			ctx,
			"fallback_domains",
			&cloudflare.ZeroTrustDeviceCustomProfileLocalDomainFallbackArgs{
				AccountId: pulumi.String(spec.AccountId),
				PolicyId:  createdProfile.ID(),
				Domains:   domains,
			},
			pulumi.Provider(cloudflareProvider),
		)
		if err != nil {
			return errors.Wrap(err, "failed to apply fallback domains")
		}
	}

	ctx.Export(OpPolicyId, createdProfile.ID())
	ctx.Export(OpGatewayUniqueId, createdProfile.GatewayUniqueId)

	return nil
}
