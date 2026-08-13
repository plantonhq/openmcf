package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2clientvpn"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// endpoint provisions the aws_ec2_client_vpn_endpoint. The endpoint carries
// everything decided at create time (authentication, client CIDR, transport,
// IP address types, transit-gateway attachment) plus the in-place dials
// (sessions, banner, connect handler, logging, DNS).
func endpoint(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*ec2clientvpn.Endpoint, error) {
	spec := locals.AwsClientVpn.Spec

	// One block per authentication option; the per-type identity source is
	// CEL-guaranteed, so each arm maps straight onto the provider block.
	// AWS accepts up to two options and a client passes on ANY of them.
	authOptions := ec2clientvpn.EndpointAuthenticationOptionArray{}
	for _, opt := range spec.AuthenticationOptions {
		args := &ec2clientvpn.EndpointAuthenticationOptionArgs{
			Type: pulumi.String(opt.Type),
		}
		if opt.RootCertificateChainArn.GetValue() != "" {
			// The client CA chain -- deliberately its own ref instead of
			// silently reusing the server certificate: the two play
			// different roles even when a self-signed setup points both
			// at the same imported certificate.
			args.RootCertificateChainArn = pulumi.String(opt.RootCertificateChainArn.GetValue())
		}
		if opt.ActiveDirectoryId != "" {
			args.ActiveDirectoryId = pulumi.String(opt.ActiveDirectoryId)
		}
		if opt.SamlProviderArn != "" {
			args.SamlProviderArn = pulumi.String(opt.SamlProviderArn)
		}
		if opt.SelfServiceSamlProviderArn != "" {
			args.SelfServiceSamlProviderArn = pulumi.String(opt.SelfServiceSamlProviderArn)
		}
		authOptions = append(authOptions, args)
	}

	// Connection logging: presence of the block is the switch. The provider
	// requires the block either way, so absence maps to enabled=false --
	// there is no separate boolean for a manifest to contradict.
	logOptions := &ec2clientvpn.EndpointConnectionLogOptionsArgs{
		Enabled: pulumi.Bool(spec.ConnectionLog != nil),
	}
	if spec.ConnectionLog != nil {
		logOptions.CloudwatchLogGroup = pulumi.String(spec.ConnectionLog.CloudwatchLogGroup.GetValue())
		if spec.ConnectionLog.CloudwatchLogStream != "" {
			logOptions.CloudwatchLogStream = pulumi.String(spec.ConnectionLog.CloudwatchLogStream)
		}
	}

	args := &ec2clientvpn.EndpointArgs{
		AuthenticationOptions: authOptions,
		ServerCertificateArn:  pulumi.String(spec.ServerCertificateArn.GetValue()),
		ConnectionLogOptions:  logOptions,
		// Split tunnel is a plain bool sent explicitly: AWS defaults to
		// full tunnel (false), and both engines must agree on the sent
		// value rather than each relying on provider defaults.
		SplitTunnel: pulumi.Bool(spec.SplitTunnel),
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Client addressing: required except for pure-IPv6 tunnel traffic
	// (CEL enforces the coupling; AWS assigns addressing there).
	if spec.ClientCidrBlock != "" {
		args.ClientCidrBlock = pulumi.String(spec.ClientCidrBlock)
	}

	if spec.TransportProtocol != "" {
		args.TransportProtocol = pulumi.String(spec.TransportProtocol)
	}
	if spec.VpnPort != nil {
		args.VpnPort = pulumi.Int(int(*spec.VpnPort))
	}
	if spec.EndpointIpAddressType != "" {
		args.EndpointIpAddressType = pulumi.String(spec.EndpointIpAddressType)
	}
	if spec.TrafficIpAddressType != "" {
		args.TrafficIpAddressType = pulumi.String(spec.TrafficIpAddressType)
	}

	// VPC attachment surface (CEL guarantees it never coexists with the
	// transit-gateway arm).
	if spec.VpcId.GetValue() != "" {
		args.VpcId = pulumi.String(spec.VpcId.GetValue())
	}
	var sgIds pulumi.StringArray
	for _, sg := range spec.SecurityGroupIds {
		if sg.GetValue() != "" {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
	}
	if len(sgIds) > 0 {
		args.SecurityGroupIds = sgIds
	}

	// Transit-gateway attachment surface.
	if spec.TransitGatewayConfiguration != nil {
		tgwArgs := &ec2clientvpn.EndpointTransitGatewayConfigurationArgs{
			TransitGatewayId: pulumi.String(spec.TransitGatewayConfiguration.TransitGatewayId.GetValue()),
		}
		if len(spec.TransitGatewayConfiguration.AvailabilityZones) > 0 {
			tgwArgs.AvailabilityZones = pulumi.ToStringArray(spec.TransitGatewayConfiguration.AvailabilityZones)
		}
		if len(spec.TransitGatewayConfiguration.AvailabilityZoneIds) > 0 {
			tgwArgs.AvailabilityZoneIds = pulumi.ToStringArray(spec.TransitGatewayConfiguration.AvailabilityZoneIds)
		}
		args.TransitGatewayConfiguration = tgwArgs
	}

	// Sessions and client experience.
	if spec.SessionTimeoutHours != nil {
		args.SessionTimeoutHours = pulumi.Int(int(*spec.SessionTimeoutHours))
	}
	// Always sent explicitly in both states (matching the Terraform module):
	// the provider declares this Optional+Computed, so an omitted value
	// absorbs whatever AWS holds and a true->false manifest flip would never
	// turn the setting off.
	args.DisconnectOnSessionTimeout = pulumi.Bool(spec.DisconnectOnSessionTimeout)
	if spec.SelfServicePortalEnabled {
		args.SelfServicePortal = pulumi.String("enabled")
	}
	// The three option blocks below are ALWAYS materialized with an explicit
	// enabled/enforced value: the provider declares them Optional+Computed and
	// diff-suppresses a missing block, so removing a once-enabled block from
	// the manifest would otherwise be a silent no-op. The disabled state
	// deliberately drops the payload, matching the provider's own expander.
	connectOpts := &ec2clientvpn.EndpointClientConnectOptionsArgs{
		Enabled: pulumi.Bool(spec.ClientConnectOptions != nil),
	}
	if spec.ClientConnectOptions != nil {
		connectOpts.LambdaFunctionArn = pulumi.String(spec.ClientConnectOptions.LambdaFunctionArn.GetValue())
	}
	args.ClientConnectOptions = connectOpts

	bannerOpts := &ec2clientvpn.EndpointClientLoginBannerOptionsArgs{
		Enabled: pulumi.Bool(spec.ClientLoginBanner != nil),
	}
	if spec.ClientLoginBanner != nil {
		bannerOpts.BannerText = pulumi.String(spec.ClientLoginBanner.BannerText)
	}
	args.ClientLoginBannerOptions = bannerOpts

	args.ClientRouteEnforcementOptions = &ec2clientvpn.EndpointClientRouteEnforcementOptionsArgs{
		Enforced: pulumi.Bool(spec.ClientRouteEnforcementEnabled),
	}
	if len(spec.DnsServers) > 0 {
		args.DnsServers = pulumi.ToStringArray(spec.DnsServers)
	}

	created, err := ec2clientvpn.NewEndpoint(ctx, locals.EndpointName, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create endpoint")
	}

	ctx.Export(OpClientVpnEndpointId, created.ID())
	ctx.Export(OpClientVpnEndpointArn, created.Arn)
	ctx.Export(OpEndpointDnsName, created.DnsName)
	ctx.Export(OpSelfServicePortalUrl, created.SelfServicePortalUrl)
	// Exported in both arms (empty for VPC-attached endpoints) so both
	// engines emit the same output set.
	if spec.TransitGatewayConfiguration != nil {
		ctx.Export(OpTransitGatewayAttachmentId,
			created.TransitGatewayConfiguration.TransitGatewayAttachmentId())
	} else {
		ctx.Export(OpTransitGatewayAttachmentId, pulumi.String(""))
	}

	return created, nil
}
