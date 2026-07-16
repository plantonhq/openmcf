# Shared VPC with Connection Cap

A guarded policy for Shared VPC topologies: the platform team owns the
policy in the host project, labels it for cost attribution, and caps how
many managed-service instances can attach before the cap is deliberately
raised.

## When to use

Shared networks where many teams create Memorystore instances and the
network owner wants an explicit ceiling on PSC endpoint sprawl — the
`limit` turns unbounded self-service into reviewed growth.

## What to customize

- `projectId` — the Shared VPC host project (the network's owner).
- `pscConfig.limit` — the connection ceiling; raising it later is an
  in-place update.
- `labels` / `description` — ownership and cost-attribution metadata.

## Composes with

`GcpVpcNetwork` (the shared network), `GcpSubnetwork` (a shared-services
subnet), and every `GcpMemorystoreInstance` teams deploy against the
shared network.
