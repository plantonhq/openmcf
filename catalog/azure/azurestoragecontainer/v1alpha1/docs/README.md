# AzureStorageContainer -- Design Research

## The Resource

A blob container (`Microsoft.Storage/storageAccounts/blobServices/
containers`) is the namespace unit of Azure blob storage -- objects live
in containers the way files live in top-level directories, and Azure
scopes public access, data-plane RBAC, encryption scopes, and lifecycle
prefixes at the container level. The component maps onto
`azurerm_storage_container` (azurerm v4.x,
`internal/services/storage/storage_container_resource.go`),
parity-verified against pulumi-azure v6 (`storage.Container`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | The v4/v5-forward ARM-id parent path (the deprecated `storage_account_name` data-plane form is not modeled); single authoritative parent FK, ForceNew |
| `name` | `container_name` | Required, ForceNew, 3-63 lowercase/digits/hyphens, no consecutive hyphens |
| `container_access_type` | enum | private (default) / blob / container; azurerm's lowercase wire values |
| `default_encryption_scope` | same | `StringValueOrRef` -> `AzureStorageEncryptionScope.encryption_scope_name` (the scope must live on the same account); ForceNew |
| `encryption_scope_override_enabled` | same | Optional bool (Azure defaults true when a scope is set); paired-with-scope CEL; ForceNew |
| `metadata` | `metadata` | Container metadata is NOT Azure tags -- ARM does not support tags on containers, so the platform's identity tags live on the account |

## Decomposition Decisions

- **First-class kind, not a fold**: containers are many-per-account with
  independent lifecycles, and consumers (function bindings, deployment
  packages, inventory destinations, immutability policies) reference
  them individually. The account's old bundled `containers` list is
  dissolved into this kind.

## Recorded Skips (with reasons)

- **Container-level immutability policy + legal hold**
  (`azurerm_storage_container_immutability_policy`) -- a genuine
  standalone resource, but its one-way `locked` state blocks deletion of
  the policy, the container, AND the account; an adoption-backlog kind,
  not a fold.
- **No `url` output** -- the container's data-plane URL is the ACCOUNT's
  blob endpoint plus the container name, and only the account knows its
  real endpoint (partitioned-DNS accounts differ). Compose it from
  `AzureStorageAccount.primary_blob_endpoint` + `container_name` instead
  of this kind exporting a value it would have to guess.

## Operational Behavior Worth Knowing

- **Containers create in seconds** and are addressed via the ARM control
  plane -- creation and verification are unaffected by the account's
  data-plane firewall.
- **Anonymous access is doubly gated**: the container's access type AND
  the account's `allow_nested_items_to_be_public` must both permit it.
- **A container cannot move between accounts** -- the parent reference is
  ForceNew.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `default_encryption_scope` → `AzureStorageEncryptionScope.status.outputs.encryption_scope_name`
  (the scope must live on the same account)
- `container_id` output ← data-plane role-assignment scopes
  (AzureRoleAssignment with Storage Blob Data roles)
- `container_name` + the account's `primary_blob_endpoint` ← SDK clients,
  function bindings, blob URLs
