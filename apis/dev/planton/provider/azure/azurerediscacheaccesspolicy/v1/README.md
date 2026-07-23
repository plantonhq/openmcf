# AzureRedisCacheAccessPolicy

A custom data-plane access policy on an Azure Cache for Redis: a named
permission set in Redis's own ACL syntax, granted to Microsoft Entra
identities through AzureRedisCacheAccessPolicyAssignment. The Redis
analog of a custom role definition.

## When to Use

Use AzureRedisCacheAccessPolicy when the built-in policies ("Data
Owner", "Data Contributor", "Data Reader") are too coarse:

- **Prefix-scoped access** -- each application confined to its own key
  namespace on a shared cache
- **Command-scoped access** -- a queue worker that can push and pop one
  key and do nothing else
- **Write-without-admin** -- full data access minus Redis's dangerous
  command category

## Key Configuration

- `redis_cache_id` -- the cache the policy is defined on (fixed at
  creation)
- `policy_name` -- what assignments reference; never a built-in name
  (the spec rejects shadowing them)
- `permissions` -- Redis ACL syntax: `+@read +@connection ~app:*`

## Composition

```yaml
redisCacheId:
  valueFrom:
    kind: AzureRedisCache
    name: app-cache
    fieldPath: status.outputs.redis_cache_id
```

Assignments grant this policy through its `access_policy_name` output --
the three-kind composition: cache → policy → grant.

## Documentation

- [Design research](docs/README.md) -- field mapping, the split verdict
- [Presets](presets/) -- prefix read-only, writer-without-admin, queue worker
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
