# AzureStorageEncryptionScope -- Design Research

## The Resource

An encryption scope (`Microsoft.Storage/storageAccounts/
encryptionScopes`) is a named encryption boundary inside a storage
account: blobs and containers that opt into a scope encrypt under the
scope's key instead of the account's, enabling per-tenant and
mixed-sensitivity key isolation without per-tenant accounts. The
component maps onto `azurerm_storage_encryption_scope` (azurerm v4.x,
`internal/services/storage/storage_encryption_scope_resource.go`),
parity-verified against pulumi-azure v6 (`storage.EncryptionScope`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | ARM-id parent path (the resource has NO legacy name form -- pure management plane); single authoritative parent FK, ForceNew |
| `name` | `scope_name` | Required, ForceNew, 4-63 alphanumerics (hyphens forbidden -- stricter than other storage names); azurerm's exact validator, mirrored as CEL |
| `source` | enum | Microsoft.Storage / Microsoft.KeyVault; required, closed |
| `key_vault_key_id` | `key_vault_key_id` | `StringValueOrRef` → `AzureKeyVaultKey.versionless_id` (rotation propagates -- the account-CMK precedent); azurerm enforces required-when-KeyVault at apply, the spec enforces it as a message CEL presence check |
| `infrastructure_encryption_required` | same | Plain bool (Azure default false), ForceNew; sent only when true on both engines (the one-way-flag convention) |

## Decomposition Decisions

- **First-class kind, not a fold**: scopes are many-per-account with
  independent lifecycles and are REFERENCED BY NAME from containers
  (`default_encryption_scope`), ADLS filesystems, and per-blob upload
  options -- the FK-referenced test. The container's former plain-string
  scope field now defaults to this kind's `encryption_scope_name`
  output.

## Recorded Skips (with reasons)

- **The reverse source/key restriction is not invented**: azurerm only
  enforces key-required-when-KeyVault; a key reference alongside a
  Microsoft.Storage source is left to ARM's own handling rather than a
  stricter-than-Azure spec rule.

## Operational Behavior Worth Knowing

- **Delete is a SOFT-DISABLE**: ARM has no true delete for scopes --
  destroy flips `state` to Disabled, the object stays GETtable, and the
  name stays reserved within the account (recreating the same name
  re-enables it). The providers treat a Disabled scope as gone; so does
  the E2E verifier (state-aware, not 404-based).
- **The Key Vault source needs account-identity plumbing**: the ACCOUNT
  must carry a managed identity with wrap/unwrap access on the key's
  vault -- the same requirement as the account-level
  `customer_managed_key`; the scope itself has no identity.
- **A scope cannot move between accounts** -- the parent reference is
  ForceNew.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `key_vault_key_id` → `AzureKeyVaultKey.status.outputs.versionless_id`
- `encryption_scope_name` output ←
  `AzureStorageContainer.default_encryption_scope` (and future ADLS
  filesystem kinds)
