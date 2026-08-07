# AzureManagedRedisGeoReplication

Links Azure Managed Redis instances into an ACTIVE geo-replication
group: every member accepts writes in its own region and Azure merges
the datasets with conflict-free semantics -- multi-primary, not the
classic primary/warm-standby model.

## When to Use

Use AzureManagedRedisGeoReplication when you need:

- **Write-anywhere global caches** -- leaderboards, presence, session
  stores, feature flags served locally on every continent
- **Active-active disaster recovery** -- losing a region loses no
  writes accepted elsewhere
- **Region evacuation as a list edit** -- removing a member's ID
  force-unlinks just that member, which keeps its data

## Key Configuration

- `managed_redis_id` -- the member the group is managed through
  (linking is reciprocal; ONE resource manages the whole group)
- `linked_managed_redis_ids` -- the other members, 1-4 of them (groups
  of up to 5)

Every member must declare the SAME `geo_replication_group_name` on its
default database, be `BALANCED_B3` or larger, have no persistence, and
use only the RediSearch/RedisJSON modules -- Azure enforces these at
link time.

## Composition

```yaml
managedRedisId:
  valueFrom:
    kind: AzureManagedRedis
    name: global-cache-east
    fieldPath: status.outputs.managed_redis_id
linkedManagedRedisIds:
  - valueFrom:
      kind: AzureManagedRedis
      name: global-cache-west
      fieldPath: status.outputs.managed_redis_id
```

Deleting the resource unlinks the members -- each keeps its own copy of
the data and becomes an independent instance again.

## Documentation

- [Presets](presets/) -- two-region pair, global mesh

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
