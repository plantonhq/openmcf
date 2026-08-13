package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vpcEndpoint creates the endpoint and exports its outputs.
func vpcEndpoint(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsVpcEndpoint.Spec

	args := &ec2.VpcEndpointArgs{
		VpcId: pulumi.String(spec.VpcId.GetValue()),
		Tags:  pulumi.ToStringMap(locals.AwsTags),
	}

	// Empty endpoint_type keeps AWS's own default (Gateway) -- passing
	// nothing and passing "Gateway" are equivalent, so only a non-empty
	// spec value is forwarded, keeping the diff surface minimal.
	if spec.EndpointType != "" {
		args.VpcEndpointType = pulumi.StringPtr(spec.EndpointType)
	}

	// Exactly one of the three service targets is set (CEL-enforced);
	// forward whichever one carries the target.
	if spec.ServiceName != "" {
		args.ServiceName = pulumi.StringPtr(spec.ServiceName)
	}
	if spec.ResourceConfigurationArn != "" {
		args.ResourceConfigurationArn = pulumi.StringPtr(spec.ResourceConfigurationArn)
	}
	if spec.ServiceNetworkArn != "" {
		args.ServiceNetworkArn = pulumi.StringPtr(spec.ServiceNetworkArn)
	}

	// Gateway endpoints attach through route tables; ENI-based types
	// attach through subnets (one ENI per subnet) and, for Interface,
	// security groups. The spec's CEL gating guarantees only the arms
	// matching the endpoint type are populated.
	if len(spec.RouteTableIds) > 0 {
		routeTableIds := make(pulumi.StringArray, 0, len(spec.RouteTableIds))
		for _, routeTable := range spec.RouteTableIds {
			routeTableIds = append(routeTableIds, pulumi.String(routeTable.GetValue()))
		}
		args.RouteTableIds = routeTableIds
	}

	if len(spec.SubnetIds) > 0 {
		subnetIds := make(pulumi.StringArray, 0, len(spec.SubnetIds))
		for _, subnet := range spec.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnet.GetValue()))
		}
		args.SubnetIds = subnetIds
	}

	if len(spec.SecurityGroupIds) > 0 {
		securityGroupIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, securityGroup := range spec.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroup.GetValue()))
		}
		args.SecurityGroupIds = securityGroupIds
	}

	// private_dns_enabled is only expressible on Interface endpoints (CEL),
	// where it updates in place. Tri-state send-when-set: the provider
	// attribute is Optional+Computed, so an omitted value keeps an existing
	// endpoint's current setting — an EXPLICIT false is the only way to
	// disable private DNS once enabled. Unset is never sent, so a gateway
	// endpoint's create call carries no DNS argument at all.
	if spec.PrivateDnsEnabled != nil {
		args.PrivateDnsEnabled = pulumi.BoolPtr(spec.GetPrivateDnsEnabled())
	}

	if spec.DnsOptions != nil {
		dnsOptions := &ec2.VpcEndpointDnsOptionsArgs{}
		if spec.DnsOptions.DnsRecordIpType != "" {
			dnsOptions.DnsRecordIpType = pulumi.StringPtr(spec.DnsOptions.DnsRecordIpType)
		}
		if spec.DnsOptions.PrivateDnsOnlyForInboundResolverEndpoint {
			dnsOptions.PrivateDnsOnlyForInboundResolverEndpoint = pulumi.BoolPtr(true)
		}
		if spec.DnsOptions.PrivateDnsPreference != "" {
			dnsOptions.PrivateDnsPreference = pulumi.StringPtr(spec.DnsOptions.PrivateDnsPreference)
		}
		if len(spec.DnsOptions.PrivateDnsSpecifiedDomains) > 0 {
			dnsOptions.PrivateDnsSpecifiedDomains = pulumi.ToStringArray(spec.DnsOptions.PrivateDnsSpecifiedDomains)
		}
		args.DnsOptions = dnsOptions
	}

	if spec.IpAddressType != "" {
		args.IpAddressType = pulumi.StringPtr(spec.IpAddressType)
	}

	// Empty policy means AWS's full-access default; forwarding nothing lets
	// AWS attach its default document. The spec carries the policy as a
	// structured document; it is serialized to JSON here.
	if spec.Policy != nil {
		policyJson, err := json.Marshal(spec.Policy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal endpoint policy")
		}
		args.Policy = pulumi.StringPtr(string(policyJson))
	}

	if len(spec.SubnetConfigurations) > 0 {
		subnetConfigurations := make(ec2.VpcEndpointSubnetConfigurationArray, 0, len(spec.SubnetConfigurations))
		for _, subnetConfiguration := range spec.SubnetConfigurations {
			configurationArgs := &ec2.VpcEndpointSubnetConfigurationArgs{
				SubnetId: pulumi.StringPtr(subnetConfiguration.SubnetId.GetValue()),
			}
			if subnetConfiguration.Ipv4 != "" {
				configurationArgs.Ipv4 = pulumi.StringPtr(subnetConfiguration.Ipv4)
			}
			if subnetConfiguration.Ipv6 != "" {
				configurationArgs.Ipv6 = pulumi.StringPtr(subnetConfiguration.Ipv6)
			}
			subnetConfigurations = append(subnetConfigurations, configurationArgs)
		}
		args.SubnetConfigurations = subnetConfigurations
	}

	if spec.ServiceRegion != "" {
		args.ServiceRegion = pulumi.StringPtr(spec.ServiceRegion)
	}

	if spec.AutoAccept {
		args.AutoAccept = pulumi.BoolPtr(true)
	}

	created, err := ec2.NewVpcEndpoint(ctx, locals.AwsVpcEndpoint.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create VPC endpoint")
	}

	ctx.Export(OpVpcEndpointId, created.ID())
	ctx.Export(OpArn, created.Arn)
	ctx.Export(OpState, created.State)
	// prefix_list_id is only populated for gateway endpoints; the
	// provider reports null for ENI-based types, exported as "".
	ctx.Export(OpPrefixListId, created.PrefixListId)

	// The first DNS entry is the endpoint's primary regional private DNS
	// name (interface endpoints only). Gateway endpoints have no DNS
	// presence -- export empty strings so downstream references resolve
	// deterministically on both endpoint types.
	ctx.Export(OpDnsName, created.DnsEntries.ApplyT(func(entries []ec2.VpcEndpointDnsEntry) string {
		if len(entries) == 0 || entries[0].DnsName == nil {
			return ""
		}
		return *entries[0].DnsName
	}).(pulumi.StringOutput))
	ctx.Export(OpHostedZoneId, created.DnsEntries.ApplyT(func(entries []ec2.VpcEndpointDnsEntry) string {
		if len(entries) == 0 || entries[0].HostedZoneId == nil {
			return ""
		}
		return *entries[0].HostedZoneId
	}).(pulumi.StringOutput))

	ctx.Export(OpNetworkInterfaceIds, created.NetworkInterfaceIds)

	return nil
}
