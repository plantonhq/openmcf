# Cross-Region Replica with VPC Placement

This preset creates a replica in a DIFFERENT region from its primary, joined to that region's VPC for private access, with custom storage and production tags. European readers get local latency; the primary stays where it is.

## When to Use

- Serving reads close to users in another geography
- A geographically redundant warm copy (promotable manually in DigitalOcean's console)
- Private-network access for workloads in the replica's region

## Key Configuration Choices

- **`region: ams3` with `vpc`** -- the replica joins the REPLICA region's VPC (the primary's VPC is in another region); both are create-only.
- **Bigger size + custom storage** -- replica size and storage must stay >= the primary's; both grow in place, never shrink.
- **Tags set at birth** -- replica tags are create-only upstream: a later retag REPLACES the replica (a fresh seed from the primary), so the full tag set is declared now.

## What You Get

A read-only endpoint in the second region with private (`private_host`/`private_uri`) and public connectivity, following the primary across the ocean -- billed like a cluster node of its slug.
