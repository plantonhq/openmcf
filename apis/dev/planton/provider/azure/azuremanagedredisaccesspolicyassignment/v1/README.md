# AzureManagedRedisAccessPolicyAssignment

A Managed Redis data-plane grant: assigns the built-in "default" access
policy to a Microsoft Entra identity on an Azure Managed Redis instance.
The Redis analog of a role assignment, and the grant half of Managed
Redis's keyless-by-default story -- with access keys off (the default),
granted identities are how clients connect at all.

## When to Use

Use AzureManagedRedisAccessPolicyAssignment when you need:

- **Secretless workload access** -- a managed identity connects with
  Entra tokens; no key exists to leak or rotate
- **Personal, auditable human access** -- on-call engineers connect as
  themselves, not through a shared key in a team vault
- **Group-based access** -- grant an Entra group once and manage
  membership in Entra, never in manifests

## Key Configuration

- `managed_redis_id` -- the instance the grant applies to (the grant
  lives on its default database)
- `object_id` -- the principal being granted; defaults to an
  AzureUserAssignedIdentity's PRINCIPAL id (never the client id)

Azure names the assignment after the object ID, so an identity is
granted at most once per database -- there is nothing else to name or
configure.

## Composition

```yaml
managedRedisId:
  valueFrom:
    kind: AzureManagedRedis
    name: app-cache
    fieldPath: status.outputs.managed_redis_id
objectId:
  valueFrom:
    kind: AzureUserAssignedIdentity
    name: app-identity
    fieldPath: status.outputs.principal_id
```

The fully keyless posture: the instance keeps access keys off (the
Managed Redis default) and every consumer connects under one of these
grants.

## Documentation

- [Design research](docs/README.md) -- field mapping, the grant-class verdict
- [Presets](presets/) -- workload identity grant, human operator grant
