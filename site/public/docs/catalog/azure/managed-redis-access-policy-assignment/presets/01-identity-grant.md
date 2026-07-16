---
title: "Workload Identity Grant"
description: "This preset grants a user-assigned managed identity data-plane access to a Managed Redis instance -- the secretless workload pattern: the application connects with its object ID as the Redis username..."
type: "preset"
rank: "01"
presetSlug: "01-identity-grant"
componentSlug: "managed-redis-access-policy-assignment"
componentTitle: "Managed Redis Access Policy Assignment"
provider: "azure"
icon: "package"
order: 1
---

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
