---
title: "Keyless Balanced Cache"
description: "This preset creates the Managed Redis default posture done right: a general-purpose Balanced instance with high availability and NO access keys -- every client authenticates with Microsoft Entra..."
type: "preset"
rank: "01"
presetSlug: "01-keyless-balanced"
componentSlug: "managed-redis"
componentTitle: "Managed Redis"
provider: "azure"
icon: "package"
order: 1
---

# Keyless Balanced Cache

This preset creates the Managed Redis default posture done right: a
general-purpose Balanced instance with high availability and NO access
keys -- every client authenticates with Microsoft Entra tokens under
data-plane grants.

## When to Use

- Application caching, session state, and rate limiting for most
  production workloads
- Teams standardizing on secretless (token-based) data access
- Any new Redis deployment -- Managed Redis is Azure's successor to the
  retiring classic Azure Cache for Redis

## Key Configuration Choices

- **BALANCED_B1** -- 1 GB in the general-purpose family; Azure allows
  many SKU changes in place, so start small and scale
- **Access keys OFF** (the Managed Redis default) -- pair with an
  AzureManagedRedisAccessPolicyAssignment per consuming identity; there
  is no secret to leak or rotate
- **High availability ON** (the default) -- a replica and the
  zone-redundant SLA; disable only for dev/test
- **ALL_KEYS_LRU eviction** -- right for pure caches where every key is
  re-derivable; keep the default VOLATILE_LRU for mixed workloads

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<region>` | Azure region, e.g. `eastus` | Your region strategy (check Managed Redis availability) |
| `<resource-group-resource-name>` | The AzureResourceGroup resource to deploy into | Your resource-group manifest |

## Composition

Grant each consuming identity data-plane access:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedRedisAccessPolicyAssignment
metadata:
  name: my-app-cache-grant
spec:
  managedRedisId:
    valueFrom:
      kind: AzureManagedRedis
      name: my-app-cache
      fieldPath: status.outputs.managed_redis_id
  objectId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: my-app-identity
      fieldPath: status.outputs.principal_id
```
