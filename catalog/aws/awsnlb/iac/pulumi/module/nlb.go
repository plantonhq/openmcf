package module

import (
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// nlb provisions the Network Load Balancer. The NLB carries no routing
// configuration by design: listeners and target groups are separate
// resources that attach to it by ARN, so this module owns only what is truly
// load-balancer-wide -- node placement with optional static IPs, security
// groups, and traffic distribution behavior. Changing "internal" replaces
// the load balancer.
func nlb(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*lb.LoadBalancer, error) {
	spec := locals.Nlb.Spec

	// AWS limits load balancer names to 32 characters; truncate
	// deterministically so the same manifest always yields the same name.
	nlbName := truncateName(locals.Nlb.Metadata.Name, 32)

	// Each mapping pins one NLB node to a subnet, optionally with a static
	// Elastic IP (internet-facing), a fixed private IPv4 address (internal),
	// or a fixed IPv6 address (dualstack) -- the static-IP story that
	// differentiates NLB from ALB. Provider-verified: modifying an EXISTING
	// mapping replaces the load balancer; pure additions do not.
	subnetMappings := lb.LoadBalancerSubnetMappingArray{}
	for _, subnetMapping := range spec.SubnetMappings {
		mapping := lb.LoadBalancerSubnetMappingArgs{
			SubnetId: pulumi.String(subnetMapping.SubnetId.GetValue()),
		}
		if subnetMapping.AllocationId.GetValue() != "" {
			mapping.AllocationId = pulumi.StringPtr(subnetMapping.AllocationId.GetValue())
		}
		if subnetMapping.PrivateIpv4Address != "" {
			mapping.PrivateIpv4Address = pulumi.StringPtr(subnetMapping.PrivateIpv4Address)
		}
		if subnetMapping.Ipv6Address != "" {
			mapping.Ipv6Address = pulumi.StringPtr(subnetMapping.Ipv6Address)
		}
		subnetMappings = append(subnetMappings, mapping)
	}

	args := &lb.LoadBalancerArgs{
		Name:                     pulumi.String(nlbName),
		LoadBalancerType:         pulumi.String("network"),
		Internal:                 pulumi.Bool(spec.Internal),
		EnableDeletionProtection: pulumi.Bool(spec.DeleteProtectionEnabled),
		SubnetMappings:           subnetMappings,
		// The Name tag carries the truncated NLB name so both engines tag the
		// resource identically (the Terraform module does the same in its
		// locals).
		Tags: pulumi.ToStringMap(stringmaps.AddEntry(locals.AwsTags, "Name", nlbName)),
	}

	// Cross-zone distribution is a real cost decision on NLB (inter-AZ data
	// transfer is billed), which is why AWS defaults it off and the spec
	// makes it an explicit opt-in.
	if spec.CrossZoneLoadBalancingEnabled {
		args.EnableCrossZoneLoadBalancing = pulumi.Bool(true)
	}

	// Only explicitly set attributes are sent, so AWS keeps its own defaults
	// for the rest -- the module never bakes in opinions the spec does not
	// express.
	if spec.IpAddressType != "" {
		args.IpAddressType = pulumi.StringPtr(spec.IpAddressType)
	}
	if spec.DnsRecordClientRoutingPolicy != "" {
		args.DnsRecordClientRoutingPolicy = pulumi.StringPtr(spec.DnsRecordClientRoutingPolicy)
	}
	if spec.ZonalShiftEnabled {
		args.EnableZonalShift = pulumi.BoolPtr(true)
	}
	if spec.EnforceSecurityGroupInboundRulesOnPrivateLinkTraffic != "" {
		args.EnforceSecurityGroupInboundRulesOnPrivateLinkTraffic = pulumi.StringPtr(spec.EnforceSecurityGroupInboundRulesOnPrivateLinkTraffic)
	}

	// Prefix-based IPv6 source NAT (dualstack only; required for UDP
	// listeners on a dualstack NLB). "on"/"off" strings mirror the
	// provider's enum.
	if spec.EnablePrefixForIpv6SourceNat != "" {
		args.EnablePrefixForIpv6SourceNat = pulumi.StringPtr(spec.EnablePrefixForIpv6SourceNat)
	}

	// Secondary private IPv4 addresses AWS auto-assigns per subnet (0-7),
	// raising the source-port budget for very high connection counts.
	// Provider-verified: DECREASING this on a live NLB forces replacement
	// (AWS cannot release secondary IPs in place).
	if spec.SecondaryIpsAutoAssignedPerSubnet != nil {
		args.SecondaryIpsAutoAssignedPerSubnet = pulumi.IntPtr(int(spec.GetSecondaryIpsAutoAssignedPerSubnet()))
	}

	// Reserved capacity (LCU reservation). Sent only when the spec asks for
	// a reservation; unset keeps normal on-demand scaling and no
	// reservation billing.
	if spec.MinimumLoadBalancerCapacityUnits != nil {
		args.MinimumLoadBalancerCapacity = &lb.LoadBalancerMinimumLoadBalancerCapacityArgs{
			CapacityUnits: pulumi.Int(int(spec.GetMinimumLoadBalancerCapacityUnits())),
		}
	}

	// Optional for NLB (unlike ALB) -- and once attached, AWS never allows
	// removing the last one, so attaching any group is a one-way door.
	if len(spec.SecurityGroups) > 0 {
		args.SecurityGroups = pulumi.ToStringArray(valuefrom.ToStringArray(spec.SecurityGroups))
	}

	// NLB access logs only capture TLS-listener traffic (an AWS limitation);
	// "enabled" is implied by the block's presence in the spec. The bucket
	// must carry the ELB log-delivery bucket policy or delivery fails
	// silently.
	if spec.AccessLogs != nil {
		accessLogs := &lb.LoadBalancerAccessLogsArgs{
			Bucket:  pulumi.String(spec.AccessLogs.Bucket.GetValue()),
			Enabled: pulumi.BoolPtr(true),
		}
		if spec.AccessLogs.Prefix != "" {
			accessLogs.Prefix = pulumi.StringPtr(spec.AccessLogs.Prefix)
		}
		args.AccessLogs = accessLogs
	}

	createdNlb, err := lb.NewLoadBalancer(ctx, nlbName, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Network Load Balancer")
	}

	ctx.Export(OpLoadBalancerArn, createdNlb.Arn)
	ctx.Export(OpLoadBalancerName, createdNlb.Name)
	ctx.Export(OpLoadBalancerDnsName, createdNlb.DnsName)
	ctx.Export(OpLoadBalancerHostedZoneId, createdNlb.ZoneId)

	return createdNlb, nil
}

// truncateName enforces AWS's 32-character load balancer name limit.
func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen]
}
