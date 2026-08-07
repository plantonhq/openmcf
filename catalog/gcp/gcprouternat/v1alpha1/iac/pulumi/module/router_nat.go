package module

import (
	"github.com/pkg/errors"
	gcprouternatv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcprouternat/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// routerNat provisions a Cloud Router with a NAT gateway — the managed
// egress path that lets instances without external IPs reach the internet
// (public NAT) or other private networks (private NAT).
//
// Manual NAT IPs are REFERENCED GcpAddress reservations, never created here
// — the reservation is its own composable node with its own lifecycle, and
// its literal IP is read from that node's outputs.
//
// Lifecycle notes the API enforces:
//   - router_name, nat_name, region, the network, endpoint_types, and type
//     are immutable — changing them replaces the resource.
//   - Everything else (IP allocation, subnetwork scoping, port tuning,
//     timeouts, rules, logging) updates in place — NAT IP rotation and
//     fleet-wide egress tuning are zero-downtime operations.
func routerNat(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
) (*compute.RouterNat, error) {
	spec := locals.GcpRouterNat.Spec

	// Enable the Compute Engine API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one gateway
	// must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"routernat-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	// The Cloud Router carrying the NAT configuration. NAT-only routers need
	// no BGP surface at all — the asn/keepalive knobs matter only when the
	// router will also terminate BGP sessions (Interconnect/VPN).
	routerArgs := &compute.RouterArgs{
		Name:    pulumi.String(spec.RouterName),
		Region:  pulumi.String(spec.Region),
		Network: pulumi.String(spec.VpcSelfLink.GetValue()),
	}
	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		routerArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.RouterAsn > 0 || spec.RouterKeepaliveInterval > 0 {
		bgpArgs := &compute.RouterBgpArgs{}
		if spec.RouterAsn > 0 {
			bgpArgs.Asn = pulumi.Int(int(spec.RouterAsn))
		}
		if spec.RouterKeepaliveInterval > 0 {
			bgpArgs.KeepaliveInterval = pulumi.IntPtr(int(spec.RouterKeepaliveInterval))
		}
		routerArgs.Bgp = bgpArgs
	}

	createdRouter, err := compute.NewRouter(ctx, "router", routerArgs,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create router")
	}

	natArgs := &compute.RouterNatArgs{
		Name:   pulumi.String(spec.NatName),
		Router: createdRouter.Name,
		Region: pulumi.String(spec.Region),
		// Which subnetworks (and which of their ranges) route through this
		// NAT (see the derivation notes below).
		SourceSubnetworkIpRangesToNat: pulumi.String(sourceSubnetworkIpRangesToNat(spec)),
	}
	if spec.ProjectId.GetValue() != "" {
		natArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// PUBLIC (internet egress) or PRIVATE (NCC spoke-to-spoke); omitted
	// means PUBLIC. Private NAT draws addresses from subnetwork ranges, so
	// the manual/auto IP machinery stays empty for it (spec CEL enforces).
	if spec.Type != "" {
		natArgs.Type = pulumi.StringPtr(spec.Type)
	}

	// Allocation is derived, not declared: referencing reservations selects
	// MANUAL_ONLY, an empty list selects AUTO_ONLY. Drained IPs keep serving
	// established connections but accept no new ones — the zero-downtime
	// path for rotating an egress IP out of service.
	if spec.Type != "PRIVATE" {
		if len(spec.NatIps) > 0 {
			natArgs.NatIpAllocateOption = pulumi.StringPtr("MANUAL_ONLY")
		} else {
			natArgs.NatIpAllocateOption = pulumi.StringPtr("AUTO_ONLY")
		}
	}
	if len(spec.NatIps) > 0 {
		natIps := pulumi.StringArray{}
		for _, natIp := range spec.NatIps {
			natIps = append(natIps, pulumi.String(natIp.GetValue()))
		}
		natArgs.NatIps = natIps
	}
	if len(spec.DrainNatIps) > 0 {
		drainIps := pulumi.StringArray{}
		for _, drainIp := range spec.DrainNatIps {
			drainIps = append(drainIps, pulumi.String(drainIp.GetValue()))
		}
		natArgs.DrainNatIps = drainIps
	}
	if spec.AutoNetworkTier != "" {
		natArgs.AutoNetworkTier = pulumi.StringPtr(spec.AutoNetworkTier)
	}

	if len(spec.Subnetworks) > 0 {
		subnetworks := compute.RouterNatSubnetworkArray{}
		for _, subnetwork := range spec.Subnetworks {
			subnetworkArgs := &compute.RouterNatSubnetworkArgs{
				Name: pulumi.String(subnetwork.Subnetwork.GetValue()),
			}
			// An empty list means everything: primary + all secondary ranges.
			if len(subnetwork.SourceIpRangesToNat) > 0 {
				subnetworkArgs.SourceIpRangesToNats = pulumi.ToStringArray(subnetwork.SourceIpRangesToNat)
			} else {
				subnetworkArgs.SourceIpRangesToNats = pulumi.StringArray{pulumi.String("ALL_IP_RANGES")}
			}
			if len(subnetwork.SecondaryIpRangeNames) > 0 {
				subnetworkArgs.SecondaryIpRangeNames = pulumi.ToStringArray(subnetwork.SecondaryIpRangeNames)
			}
			subnetworks = append(subnetworks, subnetworkArgs)
		}
		natArgs.Subnetworks = subnetworks
	}

	// Port allocation: unset fields defer to GCP's defaults (64 static / 32
	// dynamic min ports). Dynamic allocation grows a busy VM's ports toward
	// the max; it cannot coexist with endpoint-independent mapping (spec CEL
	// enforces the conflict pre-deploy, matching the API).
	if spec.MinPortsPerVm > 0 {
		natArgs.MinPortsPerVm = pulumi.IntPtr(int(spec.MinPortsPerVm))
	}
	if spec.MaxPortsPerVm > 0 {
		natArgs.MaxPortsPerVm = pulumi.IntPtr(int(spec.MaxPortsPerVm))
	}
	natArgs.EnableDynamicPortAllocation = pulumi.BoolPtr(spec.EnableDynamicPortAllocation)
	natArgs.EnableEndpointIndependentMapping = pulumi.BoolPtr(spec.EnableEndpointIndependentMapping)

	// Which resource class this NAT serves (VM instances by default).
	if len(spec.EndpointTypes) > 0 {
		natArgs.EndpointTypes = pulumi.ToStringArray(spec.EndpointTypes)
	}

	// Connection timeouts: unset fields defer to GCP's defaults (30s UDP/
	// ICMP/transitory-TCP, 1200s established-TCP, 120s TIME_WAIT).
	if spec.UdpIdleTimeoutSec > 0 {
		natArgs.UdpIdleTimeoutSec = pulumi.IntPtr(int(spec.UdpIdleTimeoutSec))
	}
	if spec.IcmpIdleTimeoutSec > 0 {
		natArgs.IcmpIdleTimeoutSec = pulumi.IntPtr(int(spec.IcmpIdleTimeoutSec))
	}
	if spec.TcpEstablishedIdleTimeoutSec > 0 {
		natArgs.TcpEstablishedIdleTimeoutSec = pulumi.IntPtr(int(spec.TcpEstablishedIdleTimeoutSec))
	}
	if spec.TcpTransitoryIdleTimeoutSec > 0 {
		natArgs.TcpTransitoryIdleTimeoutSec = pulumi.IntPtr(int(spec.TcpTransitoryIdleTimeoutSec))
	}
	if spec.TcpTimeWaitTimeoutSec > 0 {
		natArgs.TcpTimeWaitTimeoutSec = pulumi.IntPtr(int(spec.TcpTimeWaitTimeoutSec))
	}

	// NAT rules: dedicated NAT IPs/ranges for traffic matching a CEL
	// expression — e.g. a stable, separately allowlistable source IP for one
	// partner's endpoints. Lower rule numbers win.
	if len(spec.Rules) > 0 {
		rules := compute.RouterNatRuleArray{}
		for _, rule := range spec.Rules {
			ruleArgs := &compute.RouterNatRuleArgs{
				RuleNumber: pulumi.Int(int(rule.RuleNumber)),
				Match:      pulumi.String(rule.Match),
			}
			if rule.Description != "" {
				ruleArgs.Description = pulumi.StringPtr(rule.Description)
			}
			if rule.Action != nil {
				ruleArgs.Action = buildRuleAction(rule.Action)
			}
			rules = append(rules, ruleArgs)
		}
		natArgs.Rules = rules
	}

	// Logging: the DISABLED sentinel turns logging off; every other filter
	// value enables it. The filter must still be a valid value when
	// disabled, so ERRORS_ONLY is sent as a placeholder the API ignores.
	logFilter := "ERRORS_ONLY"
	if spec.LogFilter != nil {
		logFilter = spec.LogFilter.String()
	}
	if logFilter == "DISABLED" {
		natArgs.LogConfig = &compute.RouterNatLogConfigArgs{
			Enable: pulumi.Bool(false),
			Filter: pulumi.String("ERRORS_ONLY"),
		}
	} else {
		natArgs.LogConfig = &compute.RouterNatLogConfigArgs{
			Enable: pulumi.Bool(true),
			Filter: pulumi.String(logFilter),
		}
	}

	createdRouterNat, err := compute.NewRouterNat(ctx, "router-nat", natArgs,
		pulumi.Provider(gcpProvider), pulumi.Parent(createdRouter))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create router nat")
	}

	ctx.Export(OpName, createdRouterNat.Name)
	ctx.Export(OpRouterSelfLink, createdRouter.SelfLink)
	ctx.Export(OpNatIps, createdRouterNat.NatIps)

	return createdRouterNat, nil
}

// sourceSubnetworkIpRangesToNat derives the scoping mode: an explicit value
// wins; otherwise listing subnetworks implies LIST_OF_SUBNETWORKS and an
// empty list means everything in the region (primary + secondary ranges).
func sourceSubnetworkIpRangesToNat(spec *gcprouternatv1alpha1.GcpRouterNatSpec) string {
	if spec.SourceSubnetworkIpRangesToNat != "" {
		return spec.SourceSubnetworkIpRangesToNat
	}
	if len(spec.Subnetworks) > 0 {
		return "LIST_OF_SUBNETWORKS"
	}
	return "ALL_SUBNETWORKS_ALL_IP_RANGES"
}

// buildRuleAction translates a rule's action block. Public NAT rules carry
// address references; private NAT rules carry subnetwork-range references
// (the spec's CEL rules keep the two from mixing).
func buildRuleAction(action *gcprouternatv1alpha1.GcpRouterNatRuleAction) *compute.RouterNatRuleActionArgs {
	actionArgs := &compute.RouterNatRuleActionArgs{}
	if len(action.SourceNatActiveIps) > 0 {
		activeIps := pulumi.StringArray{}
		for _, activeIp := range action.SourceNatActiveIps {
			activeIps = append(activeIps, pulumi.String(activeIp.GetValue()))
		}
		actionArgs.SourceNatActiveIps = activeIps
	}
	if len(action.SourceNatDrainIps) > 0 {
		drainIps := pulumi.StringArray{}
		for _, drainIp := range action.SourceNatDrainIps {
			drainIps = append(drainIps, pulumi.String(drainIp.GetValue()))
		}
		actionArgs.SourceNatDrainIps = drainIps
	}
	if len(action.SourceNatActiveRanges) > 0 {
		activeRanges := pulumi.StringArray{}
		for _, activeRange := range action.SourceNatActiveRanges {
			activeRanges = append(activeRanges, pulumi.String(activeRange.GetValue()))
		}
		actionArgs.SourceNatActiveRanges = activeRanges
	}
	if len(action.SourceNatDrainRanges) > 0 {
		drainRanges := pulumi.StringArray{}
		for _, drainRange := range action.SourceNatDrainRanges {
			drainRanges = append(drainRanges, pulumi.String(drainRange.GetValue()))
		}
		actionArgs.SourceNatDrainRanges = drainRanges
	}
	return actionArgs
}
