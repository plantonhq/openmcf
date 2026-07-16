---
title: "Dedicated Cluster"
description: "This preset provisions a single-tenant Event Hubs cluster at the entry size (1 capacity unit) -- the top of the capacity ladder, which namespaces join via their `dedicatedClusterId` reference."
type: "preset"
rank: "01"
presetSlug: "01-dedicated-cluster"
componentSlug: "event-hub-cluster"
componentTitle: "Event Hub Cluster"
provider: "azure"
icon: "package"
order: 1
---

# Dedicated Cluster

This preset provisions a single-tenant Event Hubs cluster at the entry
size (1 capacity unit) -- the top of the capacity ladder, which
namespaces join via their `dedicatedClusterId` reference.

## When to Use

- Sustained high-throughput estates that outgrow PREMIUM's processing
  units
- Workloads needing 1024-partition hubs, 90-day retention, or
  customer-managed-key encryption
  (AzureEventHubNamespaceCustomerManagedKey)

## Key Configuration Choices

- **Provision deliberately** -- clusters bill per capacity unit per hour
  at dedicated-tier rates, the most expensive resource in the family
- **The 4-hour deletion moratorium** -- Azure forbids deleting a cluster
  for 4 hours after creation; a destroy inside that window retries for
  hours by the service's own rule
- **`capacityUnits` scales in place** -- start at 1, grow with load

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-streaming-dedicated` | The cluster name (RG-scoped) | Your naming convention |
| `capacityUnits: 1` | Guaranteed single-tenant capacity slices | Your throughput sizing |

## Downstream Wiring

Namespaces join the cluster at creation (ForceNew):

```yaml
# On an AzureEventHubNamespace
dedicatedClusterId:
  valueFrom:
    kind: AzureEventHubCluster
    name: my-streaming-cluster
    fieldPath: status.outputs.cluster_id
```
