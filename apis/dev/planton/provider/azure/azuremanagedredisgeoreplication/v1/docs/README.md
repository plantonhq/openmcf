# AzureManagedRedisGeoReplication -- Design Research

## The Resource

Active geo-replication group membership for Azure Managed Redis: every
Managed Redis database created with the same geo-replication group name
becomes a member of a multi-primary replica set once linked -- all
members accept writes, with conflict-free merge semantics (CRDT-based).
The component maps onto `azurerm_managed_redis_geo_replication`
(azurerm v4.x,
`internal/services/managedredis/managed_redis_geo_replication_resource.go`),
parity-verified against pulumi-azure v6 (`managedredis.GeoReplication`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `managed_redis_id` | same | `StringValueOrRef` → `AzureManagedRedis.managed_redis_id`; ForceNew |
| `linked_managed_redis_ids` | same | Repeated `StringValueOrRef` (1-4) → the members' `managed_redis_id` outputs |

## Split Verdict

**SPLIT -- the provider's own reasoning.** Linking mutates the
replication state of EVERY member out of band (ARM updates all linked
databases when the group changes), so membership cannot honestly be a
property of any single instance -- the provider models it as its own
resource precisely to absorb that out-of-band churn. The per-instance
half (declaring the group name) lives on AzureManagedRedisSpec; this
kind is the group-wide link/unlink operation.

## Validation Rules

- 1-4 linked members (Azure's group-of-five ceiling with the managing
  member).
- **Documented, not CEL'd** (the members are `StringValueOrRef`s and
  message-level CEL cannot dereference a ref's sub-fields -- the
  protovalidate-java constraint): the managing member must not appear
  in the linked list (Azure rejects self-links at apply time), and IDs
  must be distinct. Both fail fast at plan/apply with the provider's
  own diagnostics.
- **Cross-resource contracts left to link time**: same group name on
  every member, `BALANCED_B3`+ SKUs, no persistence, geo-compatible
  modules only -- the per-instance halves are enforced on
  AzureManagedRedisSpec's own CELs.

## Recorded Skips (with reasons)

- None -- both provider fields are modeled.

## Operational Behavior Worth Knowing

- **ONE resource per group** -- linking is reciprocal; every member
  sees the same state. Creating one resource per member would fight
  itself.
- **Deletion unlinks, it does not delete** -- every member keeps its
  own copy of the data and becomes independent; the E2E verifier is
  therefore state-aware (it reads the managing database's
  `linkedDatabases` rather than probing a 404).
- **Removing one ID force-unlinks just that member** -- the region
  evacuation primitive.
- **The resource's ID is the managing cluster's ARM ID** -- the group
  has no ARM object of its own.
- **Failover is automatic in the multi-primary model** -- applications
  point at their local member; there is no promote step (the contrast
  with classic Redis's linked-server model, where deleting the link WAS
  the failover).
