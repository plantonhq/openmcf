# AzureStorageShare -- Design Research

## The Resource

An Azure Files share (`Microsoft.Storage/storageAccounts/fileServices/
shares`) is the SMB/NFS file system unit of Azure storage -- the thing
workloads mount for shared POSIX-style state, and the level at which
Azure bills, throttles, tiers, and snapshots file storage. The component
maps onto `azurerm_storage_share` (azurerm v4.x,
`internal/services/storage/storage_share_resource.go`), parity-verified
against pulumi-azure v6 (`storage.Share`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | The v4/v5-forward ARM-id parent path (the deprecated `storage_account_name` data-plane form is not modeled); single authoritative parent FK, ForceNew |
| `name` | `share_name` | Required, ForceNew, 3-63 lowercase/digits/hyphens, no consecutive hyphens (Azure's service contract; azurerm's validator is looser and would defer the rejection to apply time) |
| `quota` | `quota_gb` | Required, 1-102400; standard accounts cap at 5120 without `large_file_share_enabled`, premium FileStorage floors at 100 -- both account-kind-dependent bounds Azure enforces at apply (cross-resource, so not spec CELs) |
| `enabled_protocol` | enum | SMB (default) / NFS; ForceNew; NFS requires a FileStorage account -- cross-resource, enforced by Azure at apply |
| `access_tier` | enum | TransactionOptimized/Hot/Cool/Premium; sent only when chosen so Azure's per-account-kind default applies |
| `acl` | `acls` | Stored access policies; the spec adds Azure's 5-policy limit and the strict `rwdl` permission-order contract azurerm leaves to the service |
| `metadata` | `metadata` | Share metadata is NOT Azure tags -- ARM does not support tags on shares, so the platform's identity tags live on the account |

## Decomposition Decisions

- **First-class kind, not a fold**: shares are many-per-account with
  independent lifecycles, quotas, and billing; mount consumers (VMs,
  AKS CSI volumes, container apps) reference them individually.

## Recorded Skips (with reasons)

- **`resource_manager_id` attribute** -- a 4.x-only compat shim
  deprecated in favor of the resource id (removed in azurerm v5); the
  `share_id` output IS the resource-manager id.
- **Directories inside the share**
  (`azurerm_storage_share_directory`) -- filesystem CONTENT, not
  infrastructure; applications create their own directory trees.
- **No `url` output** -- the share's data-plane URL is the ACCOUNT's
  file endpoint plus the share name, and only the account knows its
  real endpoint (partitioned-DNS accounts differ). Compose mount paths
  from `AzureStorageAccount.primary_file_endpoint` + `share_name`.

## Operational Behavior Worth Knowing

- **The RBAC scope is NOT the management id**: Azure Files role
  assignments target `.../fileServices/default/fileshares/{name}` (note
  the `fileshares` segment) while the management id uses `/shares/` --
  which is exactly why the module exports BOTH `share_id` and
  `rbac_scope_id`.
- **Quota grows in place**; shrinking below used capacity fails.
- **NFS shares are reachable only over private network paths** -- plan a
  private endpoint or VNet rules on the account.
- **A share cannot move between accounts** -- the parent reference is
  ForceNew.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `rbac_scope_id` output ← data-plane role-assignment scopes
  (AzureRoleAssignment with Storage File Data SMB Share roles)
- `share_name` + the account's `primary_file_endpoint` ← mount commands,
  CSI volume definitions
