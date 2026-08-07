# AzureStorageObjectReplication -- Design Research

## The Resource

An object replication policy
(`Microsoft.Storage/storageAccounts/objectReplicationPolicies`) copies
block blobs asynchronously between containers on two accounts. The
component maps onto `azurerm_storage_object_replication` (azurerm v4.x,
`internal/services/storage/storage_object_replication_resource.go`),
parity-verified against pulumi-azure v6 (`storage.ObjectReplication`).
Azure's shape is unusual: ONE logical policy is materialized as ARM
objects on BOTH accounts (the destination copy is authoritative and
assigns rule IDs; the source copy mirrors it), and the provider manages
the pair as one resource -- so does this kind.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `source_storage_account_id` | same | ARM-id FK, ForceNew |
| `destination_storage_account_id` | same | ARM-id FK, ForceNew |
| `rules` | `rules` | Required, 1-1000 (ARM's cap) |
| `rules.source_container_name` / `.destination_container_name` | same | `StringValueOrRef` → `AzureStorageContainer.container_name` -- rules compose against first-class containers |
| `rules.copy_blobs_created_after` | same | `optional` with the platform default `OnlyNewObjects`; azurerm's exact validator (OnlyNewObjects / Everything / RFC 3339) as CEL; the providers translate Everything → ARM's `1601-01-01T00:00:00Z` sentinel identically (shared expand code under the bridge) |
| `rules.filter_out_blobs_with_prefix` | `prefix_match` | RENAMED to ARM's own wire name: these are INCLUDE filters ("replicate only blobs with these prefixes" -- the provider's own docs), and the azurerm attribute name reads as the opposite. Both modules map the spec name onto their engine's attribute and carry the naming note |
| `rules.name` (Computed) | -- | The server-assigned rule ID -- read-back state, not configuration |

## Recorded Skips (with reasons)

- **`metrics_enabled` (azurerm) is NOT modeled**: pulumi-azure v6.38
  (the latest v6) has not bridged it -- `ObjectReplicationArgs` carries
  only the two account ids and rules. A one-engine-only input would
  ship a silent-drop divergence (the release_policy precedent), so the
  field waits on the bridge. Re-enable trigger: a pulumi-azure release
  whose `storage.ObjectReplicationArgs` carries `MetricsEnabled`.
- **The versioning/change-feed prerequisites stay apply-time** -- they
  live on the referenced ACCOUNTS (source: versioning + change feed;
  destination: versioning), and the spec cannot dereference a ref's
  sub-fields; Azure rejects an unprepared account with its own
  diagnostic. Documented on the spec, both modules, and the presets.

## Output Decisions

- **Both per-account policy ARM ids are exported** (the provider's own
  two attributes) plus **`policy_id`** -- the shared GUID parsed from
  the destination-side id identically on both engines; it is what
  `az storage account or-policy` and replication monitoring key on.

## Operational Behavior Worth Knowing

- **Replication is asynchronous with NO default RPO guarantee** --
  most objects land within minutes but Azure commits to nothing without
  the metrics opt-in (currently unmodelable; see the skip).
- **Everything-backfill takes time proportional to container size**;
  check rule progress before relying on the replica.
- **Cross-tenant pairs additionally need
  `cross_tenant_replication_enabled`** on the accounts (an account-spec
  field) -- same-tenant pairs, the overwhelming norm, do not.
- **Destroy removes the policy from BOTH accounts**; already-replicated
  data stays in the destination containers.

## Composition

- `source_storage_account_id` / `destination_storage_account_id` →
  `AzureStorageAccount.status.outputs.storage_account_id`
- `rules[].source_container_name` / `.destination_container_name` →
  `AzureStorageContainer.status.outputs.container_name`
- `policy_id` output ← monitoring/runbooks
  (`az storage account or-policy show --policy-id`)
