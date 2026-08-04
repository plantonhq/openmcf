# AzureManagedRedisGeoReplication

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureManagedRedisGeoReplicationSpec** links Azure Managed Redis
instances into an ACTIVE geo-replication group: every member accepts
writes in its own region and Azure merges the datasets with
conflict-free semantics -- multi-primary, not the classic
primary/warm-standby model. Applications write locally everywhere and
read their own region's instance.

Group membership is a first-class resource because Azure mutates the
replication state of EVERY member when the group changes -- linking and
unlinking are group-wide operations performed through one member, not a
property of any single instance. Each instance joins the group by
declaring the same geo_replication_group_name on its default database;
this resource then links the declared members together. Deleting it
unlinks the members (each keeps its own copy of the data and becomes an
independent instance again) -- removing a member from
linked_managed_redis_ids force-unlinks just that member, which is also
how a region is evacuated.

Requirements Azure enforces at link time (cross-resource contracts that
cannot be checked statically): every member must carry the SAME
geo_replication_group_name, be BALANCED_B3 or larger, have no
persistence enabled, and use only the RediSearch/RedisJSON modules --
the per-instance halves of these are enforced on AzureManagedRedisSpec.
The managing instance must not appear in linked_managed_redis_ids (it
is always a member implicitly); Azure rejects self-links at apply time.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedRedisGeoReplication
metadata:
  name: test-managed-redis-geo
spec:
  # The managing member; the group is reciprocal so ONE resource manages
  # the whole group.
  managedRedisId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cache/redisEnterprise/planton-hack-redis-east
  # The other members (the managing member is implicit and never
  # repeated here).
  linkedManagedRedisIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg-west/providers/Microsoft.Cache/redisEnterprise/planton-hack-redis-west
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg-europe/providers/Microsoft.Cache/redisEnterprise/planton-hack-redis-europe
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.managedRedisId` | `string \| valueFrom` | yes |  | AzureManagedRedis (`status.outputs.managed_redis_id`) |
| `spec.linkedManagedRedisIds` | `[]string \| valueFrom` | yes |  | AzureManagedRedis (`status.outputs.managed_redis_id`) |

## Field Details

### spec.managedRedisId

`string | valueFrom` · required

The Managed Redis instance through which the group is managed, by
ARM ID. References an AzureManagedRedis's managed_redis_id output.
Linking is reciprocal -- every member sees the same group state, so
ONE geo-replication resource manages the whole group (do not create
one per member). Fixed at creation.

- references: AzureManagedRedis (`status.outputs.managed_redis_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureManagedRedis, name: <that resource's name>, fieldPath: status.outputs.managed_redis_id}} -- a bare string does not parse

### spec.linkedManagedRedisIds

`[]string | valueFrom` · required

The OTHER members of the group, by ARM ID -- 1 to 4 of them, making
a group of up to 5 with the managing instance (which is always a
member and must not be repeated here). Each references an
AzureManagedRedis's managed_redis_id output. Adding an ID links
that instance into the group; removing one force-unlinks it (it
keeps its data and becomes independent).

- references: AzureManagedRedis (`status.outputs.managed_redis_id`)
- rule: {"repeated":{"minItems":"1","maxItems":"4"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureManagedRedis, name: <that resource's name>, fieldPath: status.outputs.managed_redis_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureManagedRedisGeoReplication, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.geo_replication_id` | `string` | The group's resource ID -- the ARM ID of the managing Managed Redis cluster (the group has no ARM object of its own; membership lives on every member's default database). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.managedRedisId` | AzureManagedRedis | `status.outputs.managed_redis_id` |
| `spec.linkedManagedRedisIds` | AzureManagedRedis | `status.outputs.managed_redis_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
