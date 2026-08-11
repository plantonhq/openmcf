# GcpSubnetwork Guide

The judgment this guide protects: a subnet is an IP-plan commitment.
Its primary range grows but never shrinks, its secondary ranges are what
GKE clusters stand on, and its destroy is blocked by every address still
in use — plan the CIDR math before the first apply, not after.

## Size for growth — expansion is free, shrinking is an outage

The primary range EXPANDS in place (/20 → /18); shrinking recreates the
subnet and everything addressed in it. A /20 per team-sized workload
tier is a sane default. GCP reserves 4 addresses per primary range;
proxy-only subnets (REGIONAL_MANAGED_PROXY) want /23 or larger.

## Reserved internal ranges are the enterprise path

When a network team pre-plans address space centrally,
`reservedInternalRange` (and its per-secondary-range twin) sources the
CIDR from a Network Connectivity internal range instead of a literal —
the subnet then CONSUMES an allocation instead of claiming one ad hoc.
Exactly one of the literal or the reference per range; the top-level
reference is create-time only, the per-secondary one is mutable.

## Secondary ranges are GKE's ground

Pod and service ranges are selected BY NAME (`rangeName`) from a
cluster's ip_allocation_policy — renaming one is not supported, and
removing one in use breaks the cluster standing on it. The
`sendSecondaryIpRangeIfEmpty` latch exists precisely so a partial
manifest cannot silently wipe them: leave it false unless you are
deliberately clearing all secondary ranges.

## IPv6 is three coupled decisions

`stackType` opts in (IPV4_IPV6 mutable; IPV6_ONLY recreates),
`ipv6AccessType` decides EXTERNAL (internet-routable GUAs) vs INTERNAL
(ULAs — the VPC must have its internal range enabled) and is immutable,
and the prefix pins (`externalIpv6Prefix` / `internalIpv6Prefix`) or
BYOIP (`ipCollection`, from a PublicDelegatedPrefix) only apply to their
matching access type — the spec enforces the pairings before the API
would.

## Special-purpose subnets are infrastructure others depend on

REGIONAL_MANAGED_PROXY subnets back every Envoy-based regional load
balancer in the VPC region (create one before the first regional ALB);
PRIVATE_SERVICE_CONNECT subnets back published services. The `role`
ACTIVE/BACKUP dance on proxy subnets is the drain-and-swap migration
path for proxy address space.

## Teardown discipline

GCP refuses to delete a subnet while any VM, connector, or proxy holds
an address in it — so a failing destroy usually means an inventory gap,
not a policy problem. `PREVENT` suits the subnet whose CIDR peers and
on-prem routes point at; `ABANDON` keeps the address space claimed in
the VPC while dropping management.
