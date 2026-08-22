# Cross-Region Link

This preset peers VPCs in two different regions -- private connectivity between, say, a New York fleet and a San Francisco fleet, with the traffic staying on DigitalOcean's backbone at no charge.

## When to Use

- Multi-region deployments whose services replicate or coordinate privately (database replicas, cache warming, internal APIs)
- Keeping cross-region traffic off the public internet without running a VPN

## Key Configuration Choices

- **Cross-region works out of the box** -- peering is region-agnostic; only the CIDRs must be disjoint.
- **Expect real latency** -- private does not mean local; design replication and timeouts for the inter-region round trip.
- **Pairwise only** -- peering is non-transitive; a third region needs its own peerings to each network it must reach.

## What You Get

A free private inter-region link, with per-host access control staying on each droplet's firewall.
