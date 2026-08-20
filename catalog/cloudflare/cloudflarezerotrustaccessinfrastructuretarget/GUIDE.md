# CloudflareZeroTrustAccessInfrastructureTarget guide

The judgment this guide protects you from: the hostname is not a label -- it is the ACCESS-CONTROL SURFACE. Infrastructure applications select targets by hostname (and IP), so a sloppy naming scheme becomes a sloppy permission scheme.

## Name targets like you mean it

Applications grant access to hostname patterns ("prod-db-*"). A target named casually ("box7") either falls outside every pattern (unreachable) or inside the wrong one (over-granted). Decide the scheme first -- environment-role-index ("prod-db-1") works -- and register every target inside it.

## The default virtual network is a real choice

Omitting `virtual_network_id` places the address in the account's DEFAULT virtual network -- fine for a flat estate, wrong for overlapping CIDRs. If two datacenters both use 10.0.10.0/24, each needs its own `CloudflareZeroTrustTunnelVirtualNetwork` and every target must say which one it means. The provider computes the assignment, so an omitted vnet never drifts the plan -- silence is stable, but it is still a choice.

## A target is inventory, not connectivity

Registering a target does not make it reachable: SSH sessions ride the account's tunnels, and the tunnel must route the target's network (a `CloudflareZeroTrustTunnelRoute` covering the address, in the same virtual network). Register the target, route its network, then point an infrastructure application at it -- three resources, one path.

## First real-world exercise

Cloudflare's own provider ships zero tests for this resource. The catalog's live proof lane is its first systematic real-world exercise anywhere -- treat early live behavior observations as documentation-grade findings.

## Pairs well with

- [CloudflareZeroTrustTunnelVirtualNetwork](../cloudflarezerotrusttunnelvirtualnetwork/README.md) -- the segment the address lives in.
- [CloudflareZeroTrustTunnelRoute](../cloudflarezerotrusttunnelroute/README.md) -- the route that makes it reachable.
- [CloudflareZeroTrustAccessApplication](../cloudflarezerotrustaccessapplication/README.md) -- the application that grants SSH to it.
