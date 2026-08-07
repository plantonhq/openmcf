# AzureRedisLinkedServer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureRedisLinkedServerSpec** links two PREMIUM Azure Cache for Redis
instances into a geo-replication pair: the primary serves reads and
writes while continuously replicating to the secondary in another
region, which serves as the warm disaster-recovery target.

The link is a first-class resource because DELETING it IS the failover
operation: to promote the secondary during a regional outage, the link
is removed (unlinking makes the secondary writable), traffic moves to
the promoted cache, and a new link is created in the opposite direction
once the region recovers. Applications that want a stable DNS name
across failovers point at the geo_replicated_primary_host_name output
instead of either cache's own hostname.

Requirements Azure enforces at link time: both caches must be PREMIUM
tier, in different regions, and the secondary must be the same size or
larger than the primary. The secondary rejects writes while linked.

Note: ARM has begun rejecting NEW Premium cache creations region by
region as classic Azure Cache for Redis retires in favor of Azure
Managed Redis (live-verified) -- for NEW geo-replicated deployments,
prefer AzureManagedRedis with its native geo-replication
(AzureManagedRedisGeoReplication); this kind links caches that already
exist.

Every field is fixed at creation -- changing anything replaces the link
(which is safe: replacing a link re-establishes replication; it does not
touch cached data on the primary).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRedisLinkedServer
metadata:
  name: test-redis-link
spec:
  # The PRIMARY cache -- the link's parent; its name and resource group
  # derive from this ARM id.
  targetRedisCacheId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-east/providers/Microsoft.Cache/redis/planton-hack-redis-east
  # The SECONDARY (disaster-recovery) cache in another region.
  linkedRedisCacheId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-west/providers/Microsoft.Cache/redis/planton-hack-redis-west
  linkedRedisCacheLocation:
    value: westus2
  # Exercises the role enum mapping.
  serverRole: SECONDARY
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.targetRedisCacheId` | `string \| valueFrom` | yes |  | AzureRedisCache (`status.outputs.redis_cache_id`) |
| `spec.linkedRedisCacheId` | `string \| valueFrom` | yes |  | AzureRedisCache (`status.outputs.redis_cache_id`) |
| `spec.linkedRedisCacheLocation` | `string \| valueFrom` | yes |  | AzureRedisCache (`status.outputs.region`) |
| `spec.serverRole` | `enum` |  |  |  |

## Field Details

### spec.targetRedisCacheId

`string | valueFrom` · required

The PRIMARY cache -- the one serving writes -- by ARM ID. References
an AzureRedisCache's redis_cache_id output. The link is created as a
child of this cache; its resource group and name are derived from
this ID, never spelled twice.

- references: AzureRedisCache (`status.outputs.redis_cache_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCache, name: <that resource's name>, fieldPath: status.outputs.redis_cache_id}} -- a bare string does not parse

### spec.linkedRedisCacheId

`string | valueFrom` · required

The SECONDARY cache -- the disaster-recovery replica in another
region -- by ARM ID. References an AzureRedisCache's redis_cache_id
output. Azure names the link after this cache.

- references: AzureRedisCache (`status.outputs.redis_cache_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCache, name: <that resource's name>, fieldPath: status.outputs.redis_cache_id}} -- a bare string does not parse

### spec.linkedRedisCacheLocation

`string | valueFrom` · required

The SECONDARY cache's region. References the same AzureRedisCache as
linked_redis_cache_id -- its region output -- so the location is
derived from the one referenced cache rather than hand-repeated
(a literal value is accepted when the linked cache is not managed
in the same manifest set).

- references: AzureRedisCache (`status.outputs.region`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureRedisCache, name: <that resource's name>, fieldPath: status.outputs.region}} -- a bare string does not parse

### spec.serverRole

`enum`

The role the LINKED cache plays in the pair. SECONDARY is the normal
choice: the target cache stays primary and the linked cache becomes
its read-only replica. (PRIMARY inverts the pair -- rarely used
outside re-linking after a failover.)

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_redis_linked_server_role_unspecified` -- Not specified -- invalid; choose an explicit role.
- `PRIMARY` -- The linked cache serves writes and replicates to the target cache. Rarely used outside re-linking after a failover.
- `SECONDARY` -- The linked cache is the read-only disaster-recovery replica of the target cache -- the normal geo-replication shape.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureRedisLinkedServer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.linked_server_id` | `string` | The Azure Resource Manager ID of the linked-server resource. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cache/redis/{primary}/linkedServers/{secondary} |
| `status.outputs.linked_server_name` | `string` | The link's name -- Azure names it after the LINKED (secondary) cache, so this equals the secondary cache's name. |
| `status.outputs.geo_replicated_primary_host_name` | `string` | The geo-replication DNS hostname that always resolves to the CURRENT primary of the pair. Applications that point here instead of a specific cache's hostname keep working across failovers without a connection-string change. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.targetRedisCacheId` | AzureRedisCache | `status.outputs.redis_cache_id` |
| `spec.linkedRedisCacheId` | AzureRedisCache | `status.outputs.redis_cache_id` |
| `spec.linkedRedisCacheLocation` | AzureRedisCache | `status.outputs.region` |

## See Also

- [Overview](../README.md)
