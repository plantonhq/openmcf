---
title: "Redis Cache Access Policy Assignment"
description: "Redis Cache Access Policy Assignment deployment documentation"
icon: "package"
order: 100
componentName: "azurerediscacheaccesspolicyassignment"
---

# Azure Redis Cache Access Policy Assignment

Grants a Redis data-plane access policy to a Microsoft Entra identity on
an Azure Cache for Redis -- the Redis analog of a role assignment.
Granted identities connect with Entra tokens instead of shared keys;
combined with disabling the cache's access keys, no secret exists
anywhere in the composition.

## What Gets Created

When you deploy an AzureRedisCacheAccessPolicyAssignment resource,
Planton provisions:

- **Access Policy Assignment** -- an
  `azurerm_redis_cache_access_policy_assignment` binding the policy to
  the principal on the referenced cache

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **An AzureRedisCache** with Entra authentication enabled
  (`redisConfiguration.activeDirectoryAuthenticationEnabled: true`)
- **A principal to grant** -- an AzureUserAssignedIdentity (referenced),
  or a user/group object id (literal)

## Quick Start

Create a file `grant.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRedisCacheAccessPolicyAssignment
metadata:
  name: orders-app-data-reader
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureRedisCacheAccessPolicyAssignment.orders-app-data-reader
spec:
  redisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache
      fieldPath: status.outputs.redis_cache_id
  assignmentName: orders-app-data-reader
  accessPolicyName:
    value: Data Reader
  objectId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: orders-app
      fieldPath: status.outputs.principal_id
  objectIdAlias: orders-app
```

Deploy:

```shell
planton apply -f grant.yaml
```

The granted client connects with username = the object id (or the
alias) and an Entra token as the password. For a managed identity the
object id must be the PRINCIPAL id -- granting the client id fails at
connect time, not deploy time.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `access_policy_assignment_id` | The grant's ARM id |
| `access_policy_assignment_name` | The grant's name within the cache |

## Related Resources

- [Azure Redis Cache](/docs/catalog/azure/cache-for-redis) -- the cache being granted on
- [Azure Redis Cache Access Policy](/docs/catalog/azure/redis-cache-access-policy) -- custom policies beyond the built-ins
- [Azure User Assigned Identity](/docs/catalog/azure/user-assigned-identity) -- the workload principal
