# AzureStorageLocalUser -- Design Research

## The Resource

A local user (`Microsoft.Storage/storageAccounts/localUsers`) is the
credential identity a storage account's SFTP endpoint authenticates --
Azure's answer to "partners who only speak SFTP need to land files in
blob storage." The component maps onto
`azurerm_storage_account_local_user` (azurerm v4.x,
`internal/services/storage/storage_account_local_user_resource.go`),
parity-verified against pulumi-azure v6 (`storage.LocalUser`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `storage_account_id` | `storage_account_id` | ARM-id parent FK, ForceNew |
| `name` | `user_name` | Required, ForceNew; azurerm's exact `^[a-z0-9]{3,64}$` validator as CEL |
| `ssh_key_enabled` / `ssh_password_enabled` | same | Plain bools (provider default false); the provider's at-least-one contract as a message CEL |
| `ssh_authorized_key` | `ssh_authorized_keys` | key (an OpenSSH PUBLIC key -- `sensitive_exempt_reason`, never `(sensitive)`; format CEL covers the key types Azure SFTP accepts) + description; the provider's create-path keys-iff-key-auth check as a message CEL |
| `home_directory` | same | Sent only when set |
| `permission_scope` | `permission_scopes` | service (blob/file) as a closed enum; `resource_name` a `StringValueOrRef` → `AzureStorageContainer.container_name` (file shares via explicit valueFrom → `AzureStorageShare.share_name`); the provider's one-item `permissions` wrapper list FLATTENED to the five booleans directly on the scope (a one-choice nesting is TF ergonomics -- the security-policy flatten precedent); the API's `rwdlc` wire string stays provider-internal on both engines |
| `sid` / `password` (attributes) | outputs | Both provider-sensitive; prose-documented secret-bearing OUTPUTS (the storage-account access-keys precedent -- `(sensitive)` and secret-coverage govern SPEC fields) |

## Output Decisions

- **`sftp_username` is exported as the composed login**
  (`{account}.{user}`) -- both engines derive it from the same parsed
  account name, and every SFTP client needs exactly this string.
- **The password regenerates when `ssh_password_enabled` flips false →
  true** (provider CustomizeDiff behavior) -- documented on the spec
  field; the old password stops working.

## Recorded Skips (with reasons)

- **The SFTP-connectivity contract stays apply-time-adjacent**: Azure
  ACCEPTS a local user on an account without SFTP (it just cannot log
  in), so requiring `sftp_enabled` here would be stricter than Azure;
  the pairing is documented on the spec instead. The account spec
  already CELs `sftp_enabled requires is_hns_enabled`.

## Operational Behavior Worth Knowing

- **The password is returned EXACTLY ONCE** (at the creation/update
  that enabled password auth); Azure's GET never returns it again --
  losing it means regenerating it.
- **SSH public keys are write-only in Azure's API** (the GET omits
  them); the provider tracks them from its own state, so out-of-band
  key edits are invisible to a plan.
- **Deleting the user severs the partner's access immediately** -- the
  offboarding story is `planton delete`, nothing else to clean up.

## Composition

- `storage_account_id` → `AzureStorageAccount.status.outputs.storage_account_id`
- `permission_scopes[].resource_name` → `AzureStorageContainer.status.outputs.container_name`
  (default) or `AzureStorageShare.status.outputs.share_name` (explicit
  valueFrom, service FILE)
- `sftp_username` / `password` outputs ← the partner's connection
  credentials
