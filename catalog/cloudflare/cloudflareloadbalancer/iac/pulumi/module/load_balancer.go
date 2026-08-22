package module

import (
	"github.com/pkg/errors"
	cloudflareloadbalancerv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareloadbalancer/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// load_balancer provisions the zone-scoped Cloudflare Load Balancer, wiring it to
// account-scoped pools referenced by ID/reference. Pools and their monitors are
// separate resources (CloudflareLoadBalancerPool / CloudflareLoadBalancerMonitor).
func load_balancer(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) (*cloudflare.LoadBalancer, error) {
	spec := locals.CloudflareLoadBalancer.Spec

	var defaultPools pulumi.StringArray
	for _, p := range spec.DefaultPools {
		defaultPools = append(defaultPools, pulumi.String(p.GetValue()))
	}

	args := &cloudflare.LoadBalancerArgs{
		ZoneId:       pulumi.String(spec.ZoneId.GetValue()),
		Name:         pulumi.String(spec.Hostname),
		DefaultPools: defaultPools,
		FallbackPool: pulumi.String(spec.FallbackPool.GetValue()),
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.Proxied != nil {
		args.Proxied = pulumi.BoolPtr(*spec.Proxied)
	}
	if spec.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(*spec.Enabled)
	}
	if spec.Ttl > 0 {
		args.Ttl = pulumi.Float64Ptr(spec.Ttl)
	}
	if spec.SessionAffinity != cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerSessionAffinity_none {
		args.SessionAffinity = pulumi.StringPtr(spec.SessionAffinity.String())
	}
	if spec.SteeringPolicy != cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerSteeringPolicy_off {
		args.SteeringPolicy = pulumi.StringPtr(spec.SteeringPolicy.String())
	}
	if spec.SessionAffinityTtl > 0 {
		args.SessionAffinityTtl = pulumi.Float64Ptr(spec.SessionAffinityTtl)
	}

	if saa := spec.SessionAffinityAttributes; saa != nil {
		saaArgs := &cloudflare.LoadBalancerSessionAffinityAttributesArgs{
			RequireAllHeaders: pulumi.BoolPtr(saa.RequireAllHeaders),
		}
		if saa.DrainDuration > 0 {
			saaArgs.DrainDuration = pulumi.Float64Ptr(saa.DrainDuration)
		}
		if len(saa.Headers) > 0 {
			saaArgs.Headers = pulumi.ToStringArray(saa.Headers)
		}
		if saa.Samesite != "" {
			saaArgs.Samesite = pulumi.StringPtr(saa.Samesite)
		}
		if saa.Secure != "" {
			saaArgs.Secure = pulumi.StringPtr(saa.Secure)
		}
		if saa.ZeroDowntimeFailover != "" {
			saaArgs.ZeroDowntimeFailover = pulumi.StringPtr(saa.ZeroDowntimeFailover)
		}
		args.SessionAffinityAttributes = saaArgs
	}

	if m := geoPoolMap(spec.RegionPools); len(m) > 0 {
		args.RegionPools = m
	}
	if m := geoPoolMap(spec.CountryPools); len(m) > 0 {
		args.CountryPools = m
	}
	if m := geoPoolMap(spec.PopPools); len(m) > 0 {
		args.PopPools = m
	}

	if ar := spec.AdaptiveRouting; ar != nil {
		args.AdaptiveRouting = &cloudflare.LoadBalancerAdaptiveRoutingArgs{
			FailoverAcrossPools: pulumi.BoolPtr(ar.FailoverAcrossPools),
		}
	}
	if ls := spec.LocationStrategy; ls != nil && (ls.Mode != "" || ls.PreferEcs != "") {
		lsArgs := &cloudflare.LoadBalancerLocationStrategyArgs{}
		if ls.Mode != "" {
			lsArgs.Mode = pulumi.StringPtr(ls.Mode)
		}
		if ls.PreferEcs != "" {
			lsArgs.PreferEcs = pulumi.StringPtr(ls.PreferEcs)
		}
		args.LocationStrategy = lsArgs
	}
	if rs := spec.RandomSteering; rs != nil {
		rsArgs := &cloudflare.LoadBalancerRandomSteeringArgs{}
		if rs.DefaultWeight > 0 {
			rsArgs.DefaultWeight = pulumi.Float64Ptr(rs.DefaultWeight)
		}
		if len(rs.PoolWeights) > 0 {
			weights := pulumi.Float64Map{}
			for k, v := range rs.PoolWeights {
				weights[k] = pulumi.Float64(v)
			}
			rsArgs.PoolWeights = weights
		}
		args.RandomSteering = rsArgs
	}

	if len(spec.Rules) > 0 {
		args.Rules = ruleArray(spec.Rules)
	}
	if len(spec.Networks) > 0 {
		args.Networks = pulumi.ToStringArray(spec.Networks)
	}

	created, err := cloudflare.NewLoadBalancer(
		ctx,
		"load_balancer",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloudflare load balancer")
	}

	ctx.Export(OpLoadBalancerId, created.ID())
	ctx.Export(OpLoadBalancerDnsRecordName, created.Name)
	// The CNAME target for a Cloudflare load balancer is its hostname (clients
	// point their DNS at it); it is not the opaque load-balancer ID.
	ctx.Export(OpLoadBalancerCnameTarget, created.Name)
	ctx.Export(OpZoneId, created.ZoneId)

	return created, nil
}

// ruleArray converts the spec's typed rules list into the provider's rules[]
// shape. Override semantics need real presence: an unset override inherits the
// load balancer's setting, while an explicit value -- INCLUDING "none"/"off"/0
// -- overrides it. The spec carries presence on priority, session_affinity,
// and steering_policy (proto optional), so those are sent exactly when set.
// terminates/disabled are sent only when true: false is the provider default,
// and a fixed_response rule is auto-marked terminating server-side, so an
// explicit false would fight the API's answer.
func ruleArray(rules []*cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerRule) cloudflare.LoadBalancerRuleArray {
	out := cloudflare.LoadBalancerRuleArray{}
	for _, r := range rules {
		ruleArgs := &cloudflare.LoadBalancerRuleArgs{}
		if r.Name != "" {
			ruleArgs.Name = pulumi.StringPtr(r.Name)
		}
		if r.Condition != "" {
			ruleArgs.Condition = pulumi.StringPtr(r.Condition)
		}
		if r.Priority != nil {
			ruleArgs.Priority = pulumi.IntPtr(int(*r.Priority))
		}
		if r.Disabled {
			ruleArgs.Disabled = pulumi.BoolPtr(true)
		}
		if r.Terminates {
			ruleArgs.Terminates = pulumi.BoolPtr(true)
		}
		if fr := r.FixedResponse; fr != nil {
			frArgs := &cloudflare.LoadBalancerRuleFixedResponseArgs{}
			if fr.ContentType != "" {
				frArgs.ContentType = pulumi.StringPtr(fr.ContentType)
			}
			if fr.Location != "" {
				frArgs.Location = pulumi.StringPtr(fr.Location)
			}
			if fr.MessageBody != "" {
				frArgs.MessageBody = pulumi.StringPtr(fr.MessageBody)
			}
			if fr.StatusCode != 0 {
				frArgs.StatusCode = pulumi.IntPtr(int(fr.StatusCode))
			}
			ruleArgs.FixedResponse = frArgs
		}
		if o := r.Overrides; o != nil {
			ruleArgs.Overrides = ruleOverridesArgs(o)
		}
		out = append(out, ruleArgs)
	}
	return out
}

// ruleOverridesArgs converts one rule's steering overrides, mirroring the
// top-level send conventions field by field.
func ruleOverridesArgs(o *cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerRuleOverrides) *cloudflare.LoadBalancerRuleOverridesArgs {
	args := &cloudflare.LoadBalancerRuleOverridesArgs{}
	if ar := o.AdaptiveRouting; ar != nil {
		args.AdaptiveRouting = &cloudflare.LoadBalancerRuleOverridesAdaptiveRoutingArgs{
			FailoverAcrossPools: pulumi.BoolPtr(ar.FailoverAcrossPools),
		}
	}
	if m := geoPoolMap(o.CountryPools); len(m) > 0 {
		args.CountryPools = m
	}
	if len(o.DefaultPools) > 0 {
		var pools pulumi.StringArray
		for _, p := range o.DefaultPools {
			pools = append(pools, pulumi.String(p.GetValue()))
		}
		args.DefaultPools = pools
	}
	if o.FallbackPool.GetValue() != "" {
		args.FallbackPool = pulumi.StringPtr(o.FallbackPool.GetValue())
	}
	if ls := o.LocationStrategy; ls != nil && (ls.Mode != "" || ls.PreferEcs != "") {
		lsArgs := &cloudflare.LoadBalancerRuleOverridesLocationStrategyArgs{}
		if ls.Mode != "" {
			lsArgs.Mode = pulumi.StringPtr(ls.Mode)
		}
		if ls.PreferEcs != "" {
			lsArgs.PreferEcs = pulumi.StringPtr(ls.PreferEcs)
		}
		args.LocationStrategy = lsArgs
	}
	if m := geoPoolMap(o.PopPools); len(m) > 0 {
		args.PopPools = m
	}
	if rs := o.RandomSteering; rs != nil {
		rsArgs := &cloudflare.LoadBalancerRuleOverridesRandomSteeringArgs{}
		if rs.DefaultWeight > 0 {
			rsArgs.DefaultWeight = pulumi.Float64Ptr(rs.DefaultWeight)
		}
		if len(rs.PoolWeights) > 0 {
			weights := pulumi.Float64Map{}
			for k, v := range rs.PoolWeights {
				weights[k] = pulumi.Float64(v)
			}
			rsArgs.PoolWeights = weights
		}
		args.RandomSteering = rsArgs
	}
	if m := geoPoolMap(o.RegionPools); len(m) > 0 {
		args.RegionPools = m
	}
	if o.SessionAffinity != nil {
		args.SessionAffinity = pulumi.StringPtr(o.SessionAffinity.String())
	}
	if saa := o.SessionAffinityAttributes; saa != nil {
		saaArgs := &cloudflare.LoadBalancerRuleOverridesSessionAffinityAttributesArgs{
			RequireAllHeaders: pulumi.BoolPtr(saa.RequireAllHeaders),
		}
		if saa.DrainDuration > 0 {
			saaArgs.DrainDuration = pulumi.Float64Ptr(saa.DrainDuration)
		}
		if len(saa.Headers) > 0 {
			saaArgs.Headers = pulumi.ToStringArray(saa.Headers)
		}
		if saa.Samesite != "" {
			saaArgs.Samesite = pulumi.StringPtr(saa.Samesite)
		}
		if saa.Secure != "" {
			saaArgs.Secure = pulumi.StringPtr(saa.Secure)
		}
		if saa.ZeroDowntimeFailover != "" {
			saaArgs.ZeroDowntimeFailover = pulumi.StringPtr(saa.ZeroDowntimeFailover)
		}
		args.SessionAffinityAttributes = saaArgs
	}
	if o.SessionAffinityTtl > 0 {
		args.SessionAffinityTtl = pulumi.Float64Ptr(o.SessionAffinityTtl)
	}
	if o.SteeringPolicy != nil {
		args.SteeringPolicy = pulumi.StringPtr(o.SteeringPolicy.String())
	}
	if o.Ttl > 0 {
		args.Ttl = pulumi.Float64Ptr(o.Ttl)
	}
	return args
}

// geoPoolMap converts a list of geo-pool mappings into the provider's
// map[code] -> ordered pool IDs shape.
func geoPoolMap(entries []*cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerGeoPools) pulumi.StringArrayMap {
	m := pulumi.StringArrayMap{}
	for _, gp := range entries {
		var ids pulumi.StringArray
		for _, p := range gp.PoolIds {
			ids = append(ids, pulumi.String(p.GetValue()))
		}
		m[gp.Code] = ids
	}
	return m
}
