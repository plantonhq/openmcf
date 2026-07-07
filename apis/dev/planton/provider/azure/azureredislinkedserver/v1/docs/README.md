# AzureRedisLinkedServer -- Design Research

## The Resource

A linked server (`Microsoft.Cache/redis/linkedServers`) joins two
PREMIUM caches into a geo-replication pair: the primary replicates
continuously to a warm secondary in another region. The component maps
onto `azurerm_redis_linked_server` (azurerm v4.x,
`internal/services/redis/redis_linked_server_resource.go`),
parity-verified against pulumi-azure v6 (`redis.LinkedServer`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `target_redis_cache_name` + `resource_group_name` | `target_redis_cache_id` | ONE parent reference; name and RG parse from the ARM id on both engines (the parent-derivation precedent) -- azurerm's two-string form is its legacy shape |
| `linked_redis_cache_id` | same | `StringValueOrRef` → the partner cache's `redis_cache_id` |
| `linked_redis_cache_location` | same | `StringValueOrRef` → the SAME partner cache's `region` output -- azurerm requires the caller to hand-repeat the region; referencing the cache's own output makes disagreement unrepresentable (a literal remains legal for externally-managed caches) |
| `server_role` | enum | PRIMARY/SECONDARY, closed, required |
| `name` (computed) | `linked_server_name` output | Azure names the link after the LINKED cache; there is no name argument |
| `geo_replicated_primary_host_name` (computed) | output | The DNS name that follows the CURRENT primary across failovers |

## Decomposition Decision

**SPLIT.** The link has an independent lifecycle BY DESIGN: deleting it
IS the failover operation (unlinking makes the secondary writable).
Folding geo-replication into the cache spec would turn a DR runbook
step into a spec edit on the surviving cache -- the wrong shape for the
one operation that happens during an outage.

## Recorded Skips (with reasons)

- None -- the azurerm surface is fully modeled. The two-string parent
  addressing is replaced by the ARM-id parent, not skipped.

## Operational Behavior Worth Knowing

- **Azure's link-time requirements**: both caches PREMIUM, different
  regions, secondary at least as large as the primary.
- **The classic-Redis retirement constrains NEW pairs**: ARM has begun
  rejecting new Premium cache creations region by region ("Azure Cache
  for Redis is retiring, create Azure Managed Redis instance instead"),
  so standing up a brand-new geo pair may fail in affected regions.
  Links on EXISTING Premium caches keep working; Azure Managed Redis
  carries its own geo-replication model for new deployments.
- **The secondary rejects writes while linked** -- it serves reads in
  its region.
- **Everything is ForceNew** -- replacing the link re-establishes
  replication; it does not touch cached data on the primary.
- **Failover sequence**: delete the link (secondary becomes writable) →
  repoint traffic (or rely on the geo hostname) → once the failed
  region recovers, create a new link in the opposite direction (the
  recovered cache's data is overwritten by replication).
- **Deletes poll to full detachment** -- azurerm waits until the link is
  genuinely gone (about a minute) so a re-link never races the unlink.

## Composition

- `target_redis_cache_id` / `linked_redis_cache_id` →
  `AzureRedisCache.redis_cache_id`
- `linked_redis_cache_location` → `AzureRedisCache.region`
- `geo_replicated_primary_host_name` output ← application configuration
  (the failover-stable endpoint)
