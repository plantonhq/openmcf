# Basic Cache

This preset deploys a single-node BASIC tier Redis instance — the smallest, cheapest Memorystore shape — peered to your VPC over direct peering.

## When to Use

- Development and staging environments
- Ephemeral caches where losing the data on a restart is acceptable (BASIC has no replication and no SLA)
- Rate limiters, short-lived session data, and computed-value caches that rebuild themselves

## Key Configuration Choices

- **BASIC tier, 1 GiB** — one node, no failover; a restart or maintenance event flushes the cache
- **VPC by reference** — `authorizedNetwork` resolves the `GcpVpc` node's self link, so the cache lands on the same network as its consumers
- **Direct peering** (the default `connectMode`) — the simplest connectivity; GCP picks an unused /29 automatically

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `my-gcp-project-123` | GCP project ID | GCP Console or `GcpProject` outputs |
| `my-app-vpc` | Your `GcpVpc` resource name | Your VPC manifest |

## Related Presets

- **02-ha-production** — when the cache becoming unavailable is an incident
- **03-private-services-access** — Shared VPC / PSA connectivity with read replicas and CMEK

## Related Components

- [GcpVpc](/docs/catalog/gcp/gcpvpc) — the network the cache attaches to
