# AzureRedisCache -- Design Research

## The Resource

Azure Cache for Redis (`Microsoft.Cache/redis`) is the managed
open-source-Redis service: caching, session state, leaderboards, and
pub/sub with sub-millisecond latency. The component maps onto
`azurerm_redis_cache` plus `azurerm_redis_firewall_rule` (azurerm v4.x,
`internal/services/redis/`), parity-verified field-by-field against
pulumi-azure v6 (`redis.Cache`, `redis.FirewallRule`).

Redis ENTERPRISE (`Microsoft.Cache/redisEnterprise`) is a different ARM
family with different SKUs, modules, and CMK support -- deliberately NOT
this kind; it is a separate breadth evaluation.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `cache_name` | Required, ForceNew, globally unique DNS label; azurerm ships no validator -- the CEL encodes ARM's actual rule (1-63 alnum+hyphen, no edge/consecutive hyphens) |
| `sku_name` | enum | BASIC/STANDARD/PREMIUM; unspecified deploys STANDARD; downgrade-forces-recreate documented (azurerm CustomizeDiff) |
| `family` | -- derived | Fully determined by the tier ("C"/"P"); modeling it would spell one fact twice -- both modules derive it |
| `capacity` | `capacity` | 0-6 field bound + the Premium 1-5 matrix as message CEL (azurerm's CustomizeDiff table) |
| `redis_version` | same | "4"/"6" string, matching azurerm's contract; major-only |
| `zones` | same | ForceNew; items constrained to 1-3 |
| `subnet_id` | same | `StringValueOrRef` → `AzureSubnet.subnet_id`; Premium gate as CEL |
| `private_static_ip_address` | same | ipv4-validated; requires-subnet CEL |
| `shard_count` | same | 1-10; Premium gate + mutual exclusion with replicas as CELs |
| `replicas_per_primary` | same | 1-3; Premium gate |
| `replicas_per_master` | -- skipped | ARM's legacy alias for the same setting (azurerm carries both with an equal-if-both-set rule only because it must round-trip legacy state); modeling both is contradictable redundant state |
| `non_ssl_port_enabled` | same | Plain bool, Azure default false |
| `public_network_access_enabled` | same | optional bool, default true |
| `access_keys_authentication_enabled` | same | optional bool, default true; the keys-off-requires-Entra CustomizeDiff mirrored as message CEL |
| `redis_configuration.*` | `redis_configuration` message | Full block: Entra auth, eviction policy vocabulary (Redis's own), memory dials (non-Basic CEL), keyspace events, auth toggle (VNet-only CEL), persistence auth enum, RDB/AOF (Premium CELs; conn-string-or-managed-identity CEL); all three storage connection strings `(sensitive)` |
| `redis_configuration.maxclients` | -- skipped | Computed-only in azurerm (ARM derives it from size); not an input |
| `identity` | `identity` message | SystemAssigned/UserAssigned with UAI FKs and the ids-match-type CEL (the storage-account precedent) |
| `tenant_settings` | same | Raw map passthrough |
| `patch_schedule` | `patch_schedules` | day_of_week as a closed 7-day enum (azurerm validates IsDayOfTheWeek -- ARM's Everyday/Weekend aliases are not accepted by the provider); ISO-8601 window CEL |
| firewall rules (own resource) | `firewall_rules` | Folded (the Postgres/MySQL/MSSQL verdict): pure IP filters, never FK-referenced, no life without the cache; ipv4 validation added |
| `minimum_tls_version` | -- skipped | Azure retired TLS <1.2 platform-wide (Aug 2025); azurerm v5 accepts only "1.2" -- a one-legal-value knob is dead surface, the provider default applies |

## Decomposition Decisions

- **Firewall rules and patch schedules FOLD** -- rules are pure filters;
  the patch schedule is ARM's singleton per-cache document (the storage
  lifecycle-policy class).
- **The linked server SPLITS** (`AzureRedisLinkedServer`) -- deleting
  the link IS the failover operation; a fold would make DR a spec edit
  on the surviving cache.
- **Access policies and assignments SPLIT**
  (`AzureRedisCacheAccessPolicy`/`...Assignment`) -- the assignment is
  the grant class (a module must never own grants -- the role-assignment
  precedent), and policies are many-per-cache, referenced by name.

## Recorded Skips (with reasons)

- `replicas_per_master`, `minimum_tls_version`, `maxclients` -- above.
- Redis Enterprise (CMK, RediSearch modules, active geo-replication) --
  a different ARM resource family, recorded as a breadth evaluation.

## Operational Behavior Worth Knowing

- **Provisioning is the slowest in the Azure catalog**: 15-40 minutes
  typical; azurerm's own timeouts are 3 hours.
- **RDB and AOF persistence are alternatives in practice** -- enable one;
  ARM accepts both flags but the engine persists via one pipeline.
- **Managed-identity persistence** needs the identity granted "Storage
  Blob Data Contributor" on the persistence account -- the exported
  `identity_principal_id` is the grant target.
- **A just-deleted cache name is held by Azure** -- sequential
  delete-recreate of the same name fails; the E2E scenarios use
  scenario-local fixtures for exactly this reason.

## Composition

- `resource_group` → `AzureResourceGroup.resource_group_name`
- `subnet_id` → `AzureSubnet.subnet_id`
- `identity.user_assigned_identity_ids` → `AzureUserAssignedIdentity.identity_id`
- `redis_cache_id` output ← `AzureRedisLinkedServer`,
  `AzureRedisCacheAccessPolicy`, `AzureRedisCacheAccessPolicyAssignment`,
  `AzurePrivateEndpoint`
- `region` output ← `AzureRedisLinkedServer.linked_redis_cache_location`
