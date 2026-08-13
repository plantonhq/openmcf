# GcpVpcNetwork Guide

The judgment this guide protects: the network is the resource everything
else stands on. Its own knobs look small — the expensive mistakes are the
ones made at create time (mode, IPv6, tags) and the teardown that takes a
whole environment's plumbing with it.

## Custom mode, always

`autoCreateSubnetworks: false` (the default) is the production answer.
Auto mode carves a /20 in every region on Google's schedule, burns RFC1918
space you will want later, and cannot be converted back cleanly once
subnets are in use. The only legitimate auto-mode use is a throwaway
sandbox.

## The immutables are the plan

Name, description, ULA IPv6 enablement and range, network profile, and
`deleteDefaultRoutesOnCreate` all recreate the network — and recreating a
network means recreating every subnet, peering, and attachment in it.
Decide them before the first dependent resource exists. MTU is the
exception people expect to be frozen but is not: it is mutable, though
changing it mid-flight disrupts established connections — plan a
maintenance window.

## BGP best-path selection economics

`bgpBestPathSelection` only matters once Cloud Routers carry routes from
multiple sites. `mode: LEGACY` (the default) is fine for single-site
hybrid. Move to `STANDARD` when you need deterministic multi-region
failover: it unlocks `alwaysCompareMed` (compare MED across different
neighbor ASNs — required for cost-based steering between on-prem links)
and `interRegionCost: ADD_COST_TO_MED` (bias toward the closer region).
All three update in place, network-wide — a routing-policy change is one
apply, but it applies to every consumer of the network at once.

When the block is set, the module sends `alwaysCompareMed` explicitly —
turning it off later is an ordinary spec edit (`alwaysCompareMed: false`),
not a remove-the-field dance.

## deleteDefaultRoutesOnCreate is a posture, not a cleanup

Suppressing the automatic 0.0.0.0/0 routes gives a network where NOTHING
egresses until you say so — the right posture for regulated environments.
Pair it with explicit `GcpRoute` resources or Cloud NAT; forget that and
every VM in the network has no path to the internet, including for
package installs during boot.

## Tags bind at create time only

`resourceManagerTags` (`tagKeys/{id}` → `tagValues/{id}`) exist for org
policies and IAM conditions that key on tags. They bind at create time —
adding or changing them later replaces the network, so treat them like
the name: part of the initial design, not an afterthought. Labels (from
platform metadata) remain the mutable, free-form organizational lever.

## Teardown discipline

GCP refuses to delete a network that still contains subnets, peerings, or
attached resources — a failed destroy here usually means the dependency
walk was incomplete, not that anything is wrong with the network. For
shared or long-lived networks set `deletionPolicy: PREVENT`: the network
is the classic victim of an over-broad cleanup, and PREVENT converts that
accident into a visible plan failure. `ABANDON` leaves the network
serving but unmanaged — useful when handing a network over to another
management plane, and otherwise a way to leak infrastructure.

## What is deliberately absent

Subnets, firewall rules, routes, NAT, and private-services access are
first-class sibling kinds (`GcpSubnetwork`, `GcpFirewallRule`, `GcpRoute`,
`GcpRouterNat`, `GcpGlobalAddress` + `GcpServiceNetworkingConnection`) —
compose them in a chart instead of expecting the network to carry them.
