package module

import (
	"strconv"

	"github.com/pkg/errors"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// firewall provisions the firewall and exports its ID.
func firewall(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.Firewall, error) {
	spec := locals.DigitalOceanFirewall.Spec

	inboundRules := make(digitalocean.FirewallInboundRuleArray, 0, len(spec.InboundRules))
	for _, rule := range spec.InboundRules {
		sourceDropletIds, err := refsToIntArray(rule.SourceDropletIds)
		if err != nil {
			return nil, errors.Wrap(err, "inbound rule source_droplet_ids")
		}
		args := digitalocean.FirewallInboundRuleArgs{
			Protocol:               pulumi.String(rule.Protocol),
			SourceAddresses:        stringArrayOrNil(rule.SourceAddresses),
			SourceDropletIds:       sourceDropletIds,
			SourceTags:             stringArrayOrNil(rule.SourceTags),
			SourceKubernetesIds:    refsToStringArray(rule.SourceKubernetesIds),
			SourceLoadBalancerUids: refsToStringArray(rule.SourceLoadBalancerUids),
		}
		// port_range stays unset for icmp rules: the provider drops any icmp
		// port_range on read-back, so sending one only creates rule diffs.
		if rule.PortRange != "" {
			args.PortRange = pulumi.StringPtr(rule.PortRange)
		}
		inboundRules = append(inboundRules, args)
	}

	outboundRules := make(digitalocean.FirewallOutboundRuleArray, 0, len(spec.OutboundRules))
	for _, rule := range spec.OutboundRules {
		destinationDropletIds, err := refsToIntArray(rule.DestinationDropletIds)
		if err != nil {
			return nil, errors.Wrap(err, "outbound rule destination_droplet_ids")
		}
		args := digitalocean.FirewallOutboundRuleArgs{
			Protocol:                    pulumi.String(rule.Protocol),
			DestinationAddresses:        stringArrayOrNil(rule.DestinationAddresses),
			DestinationDropletIds:       destinationDropletIds,
			DestinationTags:             stringArrayOrNil(rule.DestinationTags),
			DestinationKubernetesIds:    refsToStringArray(rule.DestinationKubernetesIds),
			DestinationLoadBalancerUids: refsToStringArray(rule.DestinationLoadBalancerUids),
		}
		if rule.PortRange != "" {
			args.PortRange = pulumi.StringPtr(rule.PortRange)
		}
		outboundRules = append(outboundRules, args)
	}

	dropletIds, err := refsToIntArray(spec.DropletIds)
	if err != nil {
		return nil, errors.Wrap(err, "droplet_ids")
	}

	createdFirewall, err := digitalocean.NewFirewall(
		ctx,
		"firewall",
		&digitalocean.FirewallArgs{
			Name:          pulumi.String(spec.FirewallName),
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
			DropletIds:    dropletIds,
			Tags:          stringArrayOrNil(spec.Tags),
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean firewall")
	}

	ctx.Export(OpFirewallId, createdFirewall.ID())

	return createdFirewall, nil
}

// refsToIntArray resolves reference fields carrying numeric Droplet IDs. The
// orchestrator resolves valueFrom references to literal values before the
// module runs, so GetValue() always carries the final string here.
func refsToIntArray(refs []*foreignkeyv1.StringValueOrRef) (pulumi.IntArray, error) {
	if len(refs) == 0 {
		return nil, nil
	}
	ids := make(pulumi.IntArray, 0, len(refs))
	for _, ref := range refs {
		id, err := strconv.Atoi(ref.GetValue())
		if err != nil {
			return nil, errors.Wrapf(err, "entry %q is not a numeric Droplet ID", ref.GetValue())
		}
		ids = append(ids, pulumi.Int(id))
	}
	return ids, nil
}

// refsToStringArray resolves reference fields carrying string IDs (Kubernetes
// cluster UUIDs, load balancer UIDs).
func refsToStringArray(refs []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	if len(refs) == 0 {
		return nil
	}
	values := make(pulumi.StringArray, 0, len(refs))
	for _, ref := range refs {
		values = append(values, pulumi.String(ref.GetValue()))
	}
	return values
}

// stringArrayOrNil keeps empty collections unset (nil, not []): the provider
// omits absent collections when it reads state back, so sending [] would
// create permanent diffs on the set-hashed rule blocks.
func stringArrayOrNil(values []string) pulumi.StringArray {
	if len(values) == 0 {
		return nil
	}
	return pulumi.ToStringArray(values)
}
