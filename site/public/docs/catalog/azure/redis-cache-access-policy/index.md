---
title: "Redis Cache Access Policy"
description: "Redis Cache Access Policy deployment documentation"
icon: "package"
order: 100
componentName: "azurerediscacheaccesspolicy"
---

# Azure Redis Cache Access Policy

Defines a CUSTOM data-plane access policy on an Azure Cache for Redis -- a named permission set, written in Redis's own ACL syntax, that Microsoft Entra identities are granted through AzureRedisCacheAccessPolicyAssignment. This is the Redis data-plane analog of a custom role definition: the policy says WHAT is allowed (commands, command categories, key patterns); the assignment says WHO gets it.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Access Policy** -- a named permission set on the referenced cache, expressed in Redis ACL syntax and updatable in place

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Azure Subscription

- **An Azure Cache for Redis** with Microsoft Entra token authentication enabled (`redisConfiguration.activeDirectoryAuthenticationEnabled: true`) -- policies gate token-authenticated clients, not access-key clients. Reference the AzureRedisCache Cloud Resource via ValueFromRef.

### When You Do NOT Need This Kind

Azure ships three built-in policies -- **Data Owner** (full access including admin commands), **Data Contributor** (read-write), and **Data Reader** (read-only) -- which assignments reference by literal name without any policy resource existing. Create a custom policy only when the built-ins are too coarse: read-write on one key prefix, no admin commands, a three-command worker.

## Deploy

### Console

Open the deployment store, find **Azure Redis Cache Access Policy**, and click **Deploy**. Start from the **Read-Only Prefix** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRedisCacheAccessPolicy
metadata:
  name: session-worker-policy
  org: acme-corp
  env: prod
spec:
  redisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache
      fieldPath: status.outputs.redis_cache_id
  policyName: session-worker
  permissions: "+get +set +del ~session:*"
```

```shell
planton apply -f access-policy.yaml
```

## Key Configuration

**Permissions (Redis ACL syntax)** -- Command/category grants followed by key patterns. Building blocks: `+@read` allows a command category, `+get` a single command, `-@dangerous` carves one out, and `~pattern` scopes which keys the grants apply to (`~*` for all keys). Semantics are ALLOW-only with carve-outs -- there is no deny-first mode. Updatable in place: an edit changes what every assignment holding this policy can do, immediately.

**Policy name** -- What assignments reference, unique within the cache. Must not collide with the built-in names. Changing the name replaces the policy and strands assignments holding the old one.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureRedisCache** | `redisCacheId` | `status.outputs.redis_cache_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `access_policy_id` | Azure Resource Manager ID of the policy | Diagnostics, audit trails |
| `access_policy_name` | The policy's name | AzureRedisCacheAccessPolicyAssignment's `accessPolicyName` |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Read-only on a prefix** -- `+@read +@connection ~metrics:*` for dashboards and debuggers. Start from the **Read-Only Prefix** preset.

**Application writer, no admin** -- `+@all -@dangerous ~app1:*` for the everyday application identity. Start from the **App Writer** preset.

**Least-privilege worker** -- exact commands on one prefix, e.g. `+get +set +del ~session:*`. Start from the **Queue Worker** preset.

## Works With

- [**Azure Redis Cache**](/cloud-catalog/azure-redis-cache) -- the cache the policy is defined on
- [**Azure Redis Cache Access Policy Assignment**](/cloud-catalog/azure-redis-cache-access-policy-assignment) -- grants this policy (by name) to a Microsoft Entra identity
