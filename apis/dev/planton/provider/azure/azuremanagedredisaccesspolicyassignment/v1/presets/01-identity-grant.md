# Workload Identity Grant

This preset grants a user-assigned managed identity data-plane access
to a Managed Redis instance -- the secretless workload pattern: the
application connects with its object ID as the Redis username and an
Entra token as the password.

## When to Use

- Every application consuming a keyless Managed Redis (access keys are
  off by default -- grants are how clients connect at all)
- AKS workloads using workload identity, App Service / Container Apps
  with managed identities, VMs with assigned identities

## Key Configuration Choices

- **`objectId` references the identity's PRINCIPAL id** -- never the
  client id; granting the client id fails at connect time, not at
  deploy time
- **One grant per identity per instance** -- Azure names the assignment
  after the object ID
- **Nothing else to configure** -- Managed Redis grants ride the
  built-in "default" access policy (full data access)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<managed-redis-resource-name>` | The AzureManagedRedis being granted on | Your cache manifest |
| `<identity-resource-name>` | The AzureUserAssignedIdentity being granted | Your identity manifest |
