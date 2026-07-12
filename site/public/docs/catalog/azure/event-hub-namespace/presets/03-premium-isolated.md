---
title: "Premium Isolated Namespace"
description: "This preset creates a PREMIUM-tier namespace: reserved processing units with predictable latency, dynamic partition scale-up, and extended retention -- the isolation tier below a whole dedicated..."
type: "preset"
rank: "03"
presetSlug: "03-premium-isolated"
componentSlug: "event-hub-namespace"
componentTitle: "Event Hub Namespace"
provider: "azure"
icon: "package"
order: 3
---

# Premium Isolated Namespace

This preset creates a PREMIUM-tier namespace: reserved processing units
with predictable latency, dynamic partition scale-up, and extended
retention -- the isolation tier below a whole dedicated cluster.

## When to Use

- Latency-sensitive production streams that must not share throughput
  with other tenants
- Workloads needing retention beyond STANDARD's 7 days or partition
  counts beyond 32
- The stepping stone before a dedicated cluster (AzureEventHubCluster)

## Key Configuration Choices

- **`sku: PREMIUM` is ForceNew across the boundary** -- Azure cannot
  convert between the reserved and multi-tenant tiers in place, so
  entering or leaving PREMIUM replaces the namespace and everything in
  it
- **`capacity: 1` processing unit** -- reserved compute, not throughput
  units; scale up in place as load grows (Azure sells 1/2/4/8/16)
- **System-assigned identity** -- the namespace's own principal, ready
  for identity-based capture auth and CMK grants

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-premium-hubs` | Globally unique namespace name | Your naming convention |
| `capacity: 1` | Processing units (1, 2, 4, 8, 16) | Your throughput sizing |

## Downstream Wiring

The identity's principal shows up as an output for grants:

```yaml
# Grant the namespace's identity Storage Blob Data Contributor for capture
principalId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: my-premium-hubs
    fieldPath: status.outputs.identity_principal_id
```
