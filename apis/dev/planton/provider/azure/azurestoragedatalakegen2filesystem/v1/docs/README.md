# AzureStorageDataLakeGen2Filesystem -- Design Research

## The Resource

A Data Lake Gen2 filesystem is the namespace unit of hierarchical-
namespace (HNS) blob storage -- the "container" the dfs endpoint and
the abfss:// driver address. The component maps onto
`azurerm_storage_data_lake_gen2_filesystem` (azurerm v4.x,
`internal/services/storage/storage_data_lake_gen2_filesystem_resource.go`),
parity-verified against pulumi-azure v6
(`storage.DataLakeGen2Filesystem`). Unlike most storage children this
is a DATA-PLANE resource: the provider creates it through the account's
dfs endpoint (shared-key auth by default), and its provider-native id
is a data-plane URL, not an ARM id.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | ARM-id parent FK, ForceNew |
| `name` | `filesystem_name` | Required, ForceNew; azurerm's exact validator mirrored as CEL: `$root` or lowercase alnum/hyphen 3-63 not starting with a hyphen |
| `default_encryption_scope` | `default_encryption_scope` | `StringValueOrRef` → `AzureStorageEncryptionScope.encryption_scope_name` (the container-kind precedent), ForceNew; sent only when set (Computed on the provider -- Azure returns a scope value even when unspecified) |
| `owner` / `group` | same | UUID or `$superuser` (azurerm's exact `validation.Any` contract as CEL); Computed -- sent only when set so Azure's defaults stand |
| `ace` | `aces` | scope (access/default) and type (user/group/mask/other) as closed enums; `id` renamed `object_id` (it IS an Entra object ID -- the catalog's established Entra vocabulary; a bare `id` inside an ACL entry reads as a resource id); permissions `[r-][w-][x-]` (giovanni's `ValidateACEPermissions` regex as CEL); the mask/other-take-no-qualifier contract (giovanni's `ACE.Validate`) as a message CEL |
| `properties` | `properties` | Metadata map; Azure requires base64-encoded VALUES (documented on the field) |

## Decomposition Decisions

- **First-class kind, not a fold**: filesystems are many-per-account
  with independent lifecycles, and they are the DATA-PLANE RBAC grant
  boundary (role assignments scope to the container-proxy id) -- the
  same split test the container kind passed.
- **`storage_data_lake_gen2_path` (directories with ACLs) is NOT
  modeled** -- adoption-backlog verdict. Data-lake zones are separated
  as per-zone FILESYSTEMS (each with its own root ACL and RBAC scope),
  which this kind covers; managing directory trees inside a filesystem
  shades into content management (the share-directories / table-entities
  class). The kind lands if adoption demands declarative intra-filesystem
  directory ACLs.

## Output Decisions

- **`filesystem_id` carries the ARM container-proxy id**
  (`{account}/blobServices/default/containers/{name}`), CONSTRUCTED
  identically on both engines -- ADLS filesystems surface in ARM as
  blob containers, and the proxy id is what role assignments scope to.
  The provider's own resource id (a dfs URL) is useless to
  management-plane consumers and differs per engine's readback shape.
- **No URL output** -- URLs compose from the account's `primary_dfs_endpoint`
  output plus `filesystem_name` (the container-kind precedent).

## Recorded Skips (with reasons)

- **The ACL-requires-HNS contract stays apply-time** -- the account
  arrives as a reference, and message-level CEL cannot dereference a
  ref's sub-fields; azurerm rejects ACLs on flat accounts at create AND
  update with its own diagnostic. Documented on the spec and both
  modules.

## Operational Behavior Worth Knowing

- **Data-plane create path**: the deploy principal reaches the dfs
  endpoint with the account's shared key (obtained via the management
  plane); an account data-plane firewall that blocks the runner blocks
  the create.
- **ACL reads require HNS** -- the provider only reads back
  owner/group/aces on HNS accounts; on flat accounts those fields stay
  empty (and setting them fails).
- **Renaming replaces the filesystem** and everything stored in it --
  the name is ForceNew on a data container.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `default_encryption_scope` → `AzureStorageEncryptionScope.status.outputs.encryption_scope_name`
- `filesystem_id` output ← `AzureRoleAssignment.scope` (data-plane
  grants at filesystem granularity)
