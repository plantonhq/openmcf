# AzureRedisCacheAccessPolicyAssignment -- Design Research

## The Resource

An access policy assignment
(`Microsoft.Cache/redis/accessPolicyAssignments`) grants a data-plane
access policy -- built-in or custom -- to a Microsoft Entra identity on
a cache: the Redis analog of a role assignment, and the grant half of
the keyless cache story. The component maps onto
`azurerm_redis_cache_access_policy_assignment` (azurerm v4.x,
`internal/services/redis/redis_cache_access_policy_assignment_resource.go`),
parity-verified against pulumi-azure v6
(`redis.CacheAccessPolicyAssignment`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `redis_cache_id` | same | `StringValueOrRef` → the cache's `redis_cache_id`; ForceNew |
| `name` | `assignment_name` | Required, ForceNew; a label for the grant (does not affect authentication) |
| `access_policy_name` | same | `StringValueOrRef` → `AzureRedisCacheAccessPolicy.access_policy_name`; built-ins ("Data Owner"/"Data Contributor"/"Data Reader") pass as literals -- they exist on every cache with no resource to reference |
| `object_id` | same | `StringValueOrRef` → `AzureUserAssignedIdentity.principal_id` (the workload-identity default); literal GUID for users/groups. The PRINCIPAL-vs-client-id trap is called out in the field comment -- the wrong one fails at connect time, not deploy time |
| `object_id_alias` | same | Required; doubles as an alternative Redis username at connect time |

## Decomposition Decision

**SPLIT.** This is the grant class: a module must never own grants (the
principle that extracted role assignments from the identity kind).
Grants are many-per-cache, span two other resources (cache + principal),
and their lifecycle -- revoke, re-grant, audit -- is independent of both.

## Recorded Skips (with reasons)

- None -- all five azurerm fields are modeled.

## Operational Behavior Worth Knowing

- **Connecting under a grant**: username = the object id (or the alias),
  password = an Entra token
  (`az account get-access-token --scope https://redis.azure.com/.default`
  for humans; the identity SDK for workloads). Entra auth must be on at
  the cache.
- **Everything is ForceNew** -- replacing an assignment momentarily
  revokes and re-grants; safe for the grant class.
- **ARM serializes assignment operations per cache** -- azurerm locks on
  the cache id; parallel grant deploys against one cache queue up.
- **The keyless posture**: cache with
  `active_directory_authentication_enabled: true` +
  `access_keys_authentication_enabled: false` + assignments = no secret
  exists anywhere in the composition.

## Composition

- `redis_cache_id` → `AzureRedisCache.redis_cache_id`
- `access_policy_name` → `AzureRedisCacheAccessPolicy.access_policy_name`
  (or a built-in literal)
- `object_id` → `AzureUserAssignedIdentity.principal_id`
