---
title: "Geo-Replicated Group Member"
description: "This preset creates one member of an ACTIVE geo-replication group -- multi-primary Redis where every region accepts writes and Azure merges the datasets with conflict-free semantics. Deploy one per..."
type: "preset"
rank: "03"
presetSlug: "03-geo-replicated"
componentSlug: "managed-redis"
componentTitle: "Managed Redis"
provider: "azure"
icon: "package"
order: 3
---

# Geo-Replicated Group Member

This preset creates one member of an ACTIVE geo-replication group --
multi-primary Redis where every region accepts writes and Azure merges
the datasets with conflict-free semantics. Deploy one per region with
the same group name, then link them with AzureManagedRedisGeoReplication.

## When to Use

- Globally distributed applications that write from every region
  (leaderboards, session stores, feature flags, presence)
- Active-active disaster recovery -- losing a region loses no writes
  accepted elsewhere
- Read-local/write-local architectures that classic Redis's
  primary/warm-standby geo-replication could not serve

## Key Configuration Choices

- **BALANCED_B3** -- the geo-replication SKU floor (B0/B1 do not
  support it); every member should be sized the same
- **The SAME `geoReplicationGroupName` on every member** -- Azure
  requires it at link time; the name is the group's identity
- **No persistence** -- Azure forbids it on geo-replicated databases;
  the cross-region replicas are the durability story
- **Module discipline** -- only RediSearch and RedisJSON are allowed on
  geo-replicated databases

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | This member's Azure region, e.g. `eastus` | One region per member |
| `<resource-group-resource-name>` | The AzureResourceGroup resource to deploy into | Your resource-group manifest |

## Composition

After deploying a sibling (e.g. `my-global-cache-west` in another
region with the same group name), link the group ONCE:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedRedisGeoReplication
metadata:
  name: my-global-cache-group
spec:
  managedRedisId:
    valueFrom:
      kind: AzureManagedRedis
      name: my-global-cache-east
      fieldPath: status.outputs.managed_redis_id
  linkedManagedRedisIds:
    - valueFrom:
        kind: AzureManagedRedis
        name: my-global-cache-west
        fieldPath: status.outputs.managed_redis_id
```
