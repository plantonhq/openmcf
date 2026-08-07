# Azure Redis Cache Access Policy

Creates a custom data-plane access policy on an Azure Cache for Redis --
a named permission set written in Redis's own ACL syntax, which
assignments then grant to Microsoft Entra identities. Use it when the
built-in "Data Owner" / "Data Contributor" / "Data Reader" policies are
too coarse.

## What Gets Created

When you deploy an AzureRedisCacheAccessPolicy resource, Planton
provisions:

- **Access Policy** -- an `azurerm_redis_cache_access_policy` on the
  referenced cache, carrying the ACL permission string

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **An AzureRedisCache** to define the policy on; enable Entra
  authentication (`redisConfiguration.activeDirectoryAuthenticationEnabled`)
  for policies to have any effect

## Quick Start

Create a file `policy.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRedisCacheAccessPolicy
metadata:
  name: orders-read-only
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureRedisCacheAccessPolicy.orders-read-only
spec:
  redisCacheId:
    valueFrom:
      kind: AzureRedisCache
      name: app-cache
      fieldPath: status.outputs.redis_cache_id
  policyName: orders-read-only
  permissions: "+@read +@connection ~orders:*"
```

Deploy:

```shell
planton apply -f policy.yaml
```

The permission string uses Redis's ACL building blocks: `+@category`
grants a command category, `+command` a single command, `-@category`
carves one out, and `~pattern` scopes which keys the grants touch.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `access_policy_id` | The policy's ARM id |
| `access_policy_name` | What AzureRedisCacheAccessPolicyAssignment references to grant the policy |

## Related Resources

- [Azure Redis Cache](/docs/catalog/azure/azurerediscache) -- the parent cache
- [Azure Redis Cache Access Policy Assignment](/docs/catalog/azure/azurerediscacheaccesspolicyassignment) -- grants this policy to an identity
