# AzureRedisCacheAccessPolicyAssignment

A Redis data-plane grant: assigns an access policy -- built-in or an
AzureRedisCacheAccessPolicy -- to a Microsoft Entra identity on an
Azure Cache for Redis. The Redis analog of a role assignment, and the
grant half of the keyless (token-only) cache story.

## When to Use

Use AzureRedisCacheAccessPolicyAssignment when you need:

- **Secretless workload access** -- a managed identity connects with
  Entra tokens instead of shared keys
- **Personal, auditable human access** -- on-call engineers get "Data
  Owner" without a key in a team vault
- **Least-privilege data access** -- pair with custom policies for
  prefix- or command-scoped grants

## Key Configuration

- `redis_cache_id` -- the cache the grant applies to
- `access_policy_name` -- "Data Owner" / "Data Contributor" /
  "Data Reader" as literals, or a custom policy by reference
- `object_id` -- the principal being granted; defaults to an
  AzureUserAssignedIdentity's PRINCIPAL id (never the client id)
- `object_id_alias` -- a readable label that doubles as an alternative
  Redis username

## Composition

```yaml
accessPolicyName:
  valueFrom:
    kind: AzureRedisCacheAccessPolicy
    name: orders-read-only
    fieldPath: status.outputs.access_policy_name
objectId:
  valueFrom:
    kind: AzureUserAssignedIdentity
    name: orders-app
    fieldPath: status.outputs.principal_id
```

The fully keyless posture: the cache enables Entra auth and disables
access keys, and every consumer connects under one of these grants.

## Documentation

- [Presets](presets/) -- identity reader, custom-policy grant, human operator
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
