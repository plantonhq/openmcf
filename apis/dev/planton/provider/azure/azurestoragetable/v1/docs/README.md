# AzureStorageTable -- Design Research

## The Resource

A Storage table (`Microsoft.Storage/storageAccounts/tableServices/
tables`) is the serverless NoSQL key/value store of Azure storage --
schemaless entities addressed by partition key + row key, petabyte
scale, no capacity provisioning. The component maps onto
`azurerm_storage_table` (azurerm v4.x,
`internal/services/storage/storage_table_resource.go`), parity-verified
against pulumi-azure v6 (`storage.Table`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | The v4/v5-forward ARM-id parent path (the deprecated `storage_account_name` data-plane form is not modeled in the SPEC); single authoritative parent FK, ForceNew. See the engine-addressing exception below |
| `name` | `table_name` | Required, ForceNew, letter-start 3-63 alphanumerics, never the literal `"table"` -- azurerm's exact validator, mirrored as CEL |
| `acl` | `acls` | Stored access policies; start + expiry are REQUIRED on table policies (azurerm's contract); the spec adds Azure's 5-policy limit and the strict `raud` permission-order contract |

## The Engine-Addressing PARITY-EXCEPTION

pulumi-azure v6 has not bridged the table's `storage_account_id` input
(verified against v6.38, the LATEST v6 -- no SDK bump closes it; the
bridged `storage.Table` carries only the deprecated-direction
`storageAccountName`). So:

- The **Terraform module** addresses the table by `storage_account_id`
  (the resource-manager path).
- The **Pulumi module** parses the account NAME from the same resolved
  ARM id and passes `storageAccountName`.

The created table is identical and ALL stack outputs match
byte-for-byte -- both engines natively export the same
`resource_manager_id`, which is what the `table_id` output carries.
Only the provider's internal addressing differs. The exception is
documented in both modules and dissolves when a bridge release carries
`storageAccountId` on `storage.Table`.

## Recorded Skips (with reasons)

- **Table entities** (`azurerm_storage_table_entity`) -- table CONTENT
  (a single row), not infrastructure; applications own their rows.
- **No `url` output** -- the table's data-plane URL is the ACCOUNT's
  table endpoint plus the table name, and only the account knows its
  real endpoint (partitioned-DNS accounts differ). Compose client URLs
  from `AzureStorageAccount.primary_table_endpoint` + `table_name`.

## Operational Behavior Worth Knowing

- **Creation and ACLs ride the DATA PLANE with shared-key auth** on
  both engines (unlike shares/queues/containers, whose management rides
  ARM): the parent account must keep `shared_access_key_enabled` true
  (Azure's default) for deploys to work -- documented on the spec.
  VERIFICATION still rides ARM: the management API exposes a
  first-class table read proxy.
- **Partition design is the scalability lever**: one partition is one
  server's throughput; entities sharing a partition key support atomic
  batch transactions, entities across partitions do not.
- **A table cannot move between accounts** -- the parent reference is
  ForceNew.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `table_id` output ← data-plane role-assignment scopes
  (AzureRoleAssignment with Storage Table Data roles)
- `table_name` + the account's `primary_table_endpoint` ← SDK clients,
  Functions table bindings
