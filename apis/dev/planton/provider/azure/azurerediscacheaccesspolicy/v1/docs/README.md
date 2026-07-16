# AzureRedisCacheAccessPolicy -- Design Research

## The Resource

A cache access policy (`Microsoft.Cache/redis/accessPolicies`) is a
named, custom data-plane permission set in Redis's own ACL syntax,
granted to Microsoft Entra identities through access policy
assignments. The component maps onto
`azurerm_redis_cache_access_policy` (azurerm v4.x,
`internal/services/redis/redis_cache_access_policy_resource.go`),
parity-verified against pulumi-azure v6 (`redis.CacheAccessPolicy`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `redis_cache_id` | same | `StringValueOrRef` → the cache's `redis_cache_id`; ForceNew |
| `name` | `policy_name` | Required, ForceNew; a CEL rejects the three built-in names ("Data Owner"/"Data Contributor"/"Data Reader") -- built-ins exist on every cache and are referenced directly by assignments, so shadowing one is always a mistake |
| `permissions` | same | Raw Redis ACL syntax, updatable in place; the spec documents the building blocks (`+@category`, `+command`, `-@carveout`, `~pattern`) rather than inventing an abstraction over Redis's own vocabulary |

## Decomposition Decision

**SPLIT.** Policies are many-per-cache with independent lifecycles and
are referenced BY NAME from assignments -- the FK-referenced test. The
policy/assignment pair is the Redis data-plane analog of
role-definition/role-assignment, and the same reasoning applies: the
definition is a reusable, referenceable object, never a fold.

## Recorded Skips (with reasons)

- None -- the azurerm surface is three fields and all three are modeled.

## Operational Behavior Worth Knowing

- **Policies only gate token-authenticated clients** -- Entra auth must
  be on at the cache for them to matter; access-key clients bypass the
  ACL system entirely (which is why the keyless posture turns keys off).
- **Permissions update in place** -- widening or narrowing a policy
  never recreates the assignments granting it.
- **ARM serializes policy operations per cache** -- azurerm locks on the
  cache id; parallel policy deploys against one cache queue up.

## Composition

- `redis_cache_id` → `AzureRedisCache.redis_cache_id`
- `access_policy_name` output ←
  `AzureRedisCacheAccessPolicyAssignment.access_policy_name`
