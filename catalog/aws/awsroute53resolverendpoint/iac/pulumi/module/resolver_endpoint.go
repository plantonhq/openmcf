package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// resolverEndpoint creates the endpoint, its forwarding rules and
// their VPC associations, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - Direction and SecurityGroupIds are ForceNew on the endpoint;
//     IpAddresses churn in place (added before removed, so the
//     2-address floor never breaks);
//   - a rule's ResolverEndpointId updates in place EXCEPT detaching it
//     (endpoint id -> empty), which the provider forces to replace;
//     SYSTEM rules never carry one;
//   - rule associations are pure joins (rule, vpc) - every argument
//     ForceNew, no update path;
//   - AWS strips a trailing dot from rule domain names (the provider
//     normalizes both ways, so no drift);
//   - tags on a RAM-shared rule are never read back (the provider
//     skips them for SHARED_WITH_ME rules).
func resolverEndpoint(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	ipAddresses := route53.ResolverEndpointIpAddressArray{}
	for _, entry := range spec.IpAddresses {
		ipAddress := &route53.ResolverEndpointIpAddressArgs{
			SubnetId: pulumi.String(entry.SubnetId.GetValue()),
		}
		if entry.Ip != "" {
			ipAddress.Ip = pulumi.String(entry.Ip)
		}
		if entry.Ipv6 != "" {
			ipAddress.Ipv6 = pulumi.String(entry.Ipv6)
		}
		ipAddresses = append(ipAddresses, ipAddress)
	}

	securityGroupIds := pulumi.StringArray{}
	for _, securityGroupId := range spec.SecurityGroupIds {
		securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
	}

	endpointArgs := &route53.ResolverEndpointArgs{
		Name:             pulumi.String(locals.Target.Metadata.Name),
		Direction:        pulumi.String(spec.Direction),
		IpAddresses:      ipAddresses,
		SecurityGroupIds: securityGroupIds,
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.EndpointType != "" {
		endpointArgs.ResolverEndpointType = pulumi.String(spec.EndpointType)
	}
	if len(spec.Protocols) > 0 {
		endpointArgs.Protocols = pulumi.ToStringArray(spec.Protocols)
	}
	// Tri-state metrics toggles: nil leaves AWS's default, an explicit
	// value (true or false) is sent as stated.
	if spec.RniEnhancedMetricsEnabled != nil {
		endpointArgs.RniEnhancedMetricsEnabled = pulumi.Bool(*spec.RniEnhancedMetricsEnabled)
	}
	if spec.TargetNameServerMetricsEnabled != nil {
		endpointArgs.TargetNameServerMetricsEnabled = pulumi.Bool(*spec.TargetNameServerMetricsEnabled)
	}

	createdEndpoint, err := route53.NewResolverEndpoint(ctx, "endpoint", endpointArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create endpoint")
	}

	// Forwarding rules, one per spec entry. SYSTEM rules carry no
	// endpoint binding (they restore recursive resolution); FORWARD
	// and DELEGATE rules bind to this endpoint.
	ruleIds := pulumi.StringMap{}
	ruleAssociationIds := pulumi.StringMap{}
	for _, rule := range spec.Rules {
		ruleArgs := &route53.ResolverRuleArgs{
			Name:       pulumi.String(rule.Name),
			DomainName: pulumi.String(rule.DomainName),
			RuleType:   pulumi.String(rule.RuleType),
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}
		if rule.RuleType != "SYSTEM" {
			ruleArgs.ResolverEndpointId = createdEndpoint.ID()
		}
		if len(rule.TargetIps) > 0 {
			targetIps := route53.ResolverRuleTargetIpArray{}
			for _, target := range rule.TargetIps {
				targetIp := &route53.ResolverRuleTargetIpArgs{}
				if target.Ip != "" {
					targetIp.Ip = pulumi.String(target.Ip)
				}
				if target.Ipv6 != "" {
					targetIp.Ipv6 = pulumi.String(target.Ipv6)
				}
				if target.Port > 0 {
					targetIp.Port = pulumi.Int(int(target.Port))
				}
				if target.Protocol != "" {
					targetIp.Protocol = pulumi.String(target.Protocol)
				}
				targetIps = append(targetIps, targetIp)
			}
			ruleArgs.TargetIps = targetIps
		}

		createdRule, err := route53.NewResolverRule(ctx, fmt.Sprintf("rule-%s", rule.Name), ruleArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create rule %s", rule.Name)
		}
		ruleIds[rule.Name] = createdRule.ID().ToStringOutput()

		// One association per (rule, vpc) pair, keyed "rule//vpc" -
		// the same keying as the Terraform module. The cosmetic
		// association name is deterministic and capped at the
		// provider's 64-character wall, identical across both engines
		// (Pulumi would otherwise auto-generate one from the URN).
		for _, vpcId := range rule.VpcIds {
			resolvedVpcId := vpcId.GetValue()
			associationName := fmt.Sprintf("%s-%s", rule.Name, resolvedVpcId)
			if len(associationName) > 64 {
				associationName = associationName[:64]
			}
			createdAssociation, err := route53.NewResolverRuleAssociation(ctx,
				fmt.Sprintf("rule-association-%s-%s", rule.Name, resolvedVpcId),
				&route53.ResolverRuleAssociationArgs{
					Name:           pulumi.String(associationName),
					ResolverRuleId: createdRule.ID(),
					VpcId:          pulumi.String(resolvedVpcId),
				}, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "associate rule %s to vpc %s", rule.Name, resolvedVpcId)
			}
			ruleAssociationIds[fmt.Sprintf("%s//%s", rule.Name, resolvedVpcId)] = createdAssociation.ID().ToStringOutput()
		}
	}

	ctx.Export(OpEndpointId, createdEndpoint.ID())
	ctx.Export(OpEndpointArn, createdEndpoint.Arn)
	ctx.Export(OpHostVpcId, createdEndpoint.HostVpcId)
	ctx.Export(OpIpAddresses, createdEndpoint.IpAddresses.ApplyT(func(entries []route53.ResolverEndpointIpAddress) []string {
		ips := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.Ip != nil {
				ips = append(ips, *entry.Ip)
			}
		}
		return ips
	}).(pulumi.StringArrayOutput))
	ctx.Export(OpRuleIds, ruleIds)
	ctx.Export(OpRuleAssociationIds, ruleAssociationIds)
	return nil
}
