# Memorystore for Valkey Policy

The policy that unlocks Memorystore for Valkey on a network: authorizes the
`gcp-memorystore` service class in one region and gives the connectivity
automation a subnet to place PSC endpoints in.

## When to use

Deploy this once per (network, region) before the first
`GcpMemorystoreInstance` there. Without it, instance creation fails with a
connectivity error; with it, endpoints appear automatically and the
instance's discovery IP surfaces on its outputs.

## What to customize

- `location` — match the region your Memorystore instances will live in.
- `network` / `pscConfig.subnetworks` — point at your VPC and a subnet in
  the same region. Regular-purpose subnets work.
- Add `pscConfig.limit` to cap how many instances can attach through this
  policy in a shared network.

## Composes with

`GcpVpcNetwork` and `GcpSubnetwork` upstream; `GcpMemorystoreInstance`
downstream (deploy the instance after this policy — keep the policy alive
as long as any instance depends on it).
