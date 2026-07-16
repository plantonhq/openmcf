---
title: "Workload Identity Data Reader Grant"
description: "This preset grants the built-in read-only policy to a user-assigned managed identity -- the standard first grant on the road to a keyless cache."
type: "preset"
rank: "01"
presetSlug: "01-identity-data-reader"
componentSlug: "redis-cache-access-policy-assignment"
componentTitle: "Redis Cache Access Policy Assignment"
provider: "azure"
icon: "package"
order: 1
---

# Workload Identity Data Reader Grant

This preset grants the built-in read-only policy to a user-assigned
managed identity -- the standard first grant on the road to a keyless
cache.

## When to Use

- Read-only workloads (dashboards, report generators, cache warmers'
  verification side) running under a managed identity
- The template for every workload-identity grant: swap the policy name
  for "Data Contributor" or a custom policy as the workload demands

## Key Configuration Choices

- **Built-in policies are literals** -- "Data Owner", "Data Contributor",
  and "Data Reader" exist on every cache; only custom policies are
  referenced through AzureRedisCacheAccessPolicy
- **`objectId` references the PRINCIPAL id** -- the default FK wiring
  does this for you; granting the client id instead is the classic
  mistake that fails at connect time
- **Entra auth must be on** -- the cache needs
  `redisConfiguration.activeDirectoryAuthenticationEnabled: true` for
  any assignment to matter
- **Connecting** -- the client authenticates with the object id (or the
  alias set here) as the Redis username and an Entra token as the
  password

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<assignmentName>` | The grant's name, unique within the cache | Convention: `<identity>-<policy>`, e.g. `orders-app-data-reader` |
| `<identity-resource-name>` | The AzureUserAssignedIdentity's Planton resource name | Your identity composition |
