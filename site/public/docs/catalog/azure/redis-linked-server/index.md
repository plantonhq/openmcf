---
title: "Redis Linked Server"
description: "Redis Linked Server deployment documentation"
icon: "package"
order: 100
componentName: "azureredislinkedserver"
---

# Azure Redis Linked Server

Creates a geo-replication link between two Premium Azure Cache for Redis
instances -- the primary continuously replicates to a warm secondary in
another region, and deleting the link is itself the failover operation
that promotes the secondary.

## What Gets Created

When you deploy an AzureRedisLinkedServer resource, Planton provisions:

- **Linked Server** -- an `azurerm_redis_linked_server` on the primary
  cache, pairing it with the secondary; Azure names the link after the
  secondary cache

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **Two PREMIUM AzureRedisCache instances** in different regions, the
  secondary at least as large as the primary

## Quick Start

Create a file `geo-link.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRedisLinkedServer
metadata:
  name: app-cache-geo-link
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureRedisLinkedServer.app-cache-geo-link
spec:
  targetRedisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache-east
      fieldPath: status.outputs.redis_cache_id
  linkedRedisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache-west
      fieldPath: status.outputs.redis_cache_id
  linkedRedisCacheLocation:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache-west
      fieldPath: status.outputs.region
  serverRole: SECONDARY
```

Deploy:

```shell
planton apply -f geo-link.yaml
```

Point applications at the `geo_replicated_primary_host_name` output
instead of either cache's own hostname -- it follows the current primary
across failovers, so a failover needs no connection-string change.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `linked_server_id` | The link's ARM id |
| `linked_server_name` | The link's name (equals the secondary cache's name) |
| `geo_replicated_primary_host_name` | The failover-stable endpoint applications should use |

## Related Resources

- [Azure Redis Cache](/docs/catalog/azure/cache-for-redis) -- both ends of the pair
