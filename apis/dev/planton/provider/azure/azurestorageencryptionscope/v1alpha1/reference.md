# AzureStorageEncryptionScope

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureStorageEncryptionScopeSpec** defines the configuration for
creating an encryption scope inside an Azure Storage Account: a named
encryption boundary that lets different data in ONE account encrypt
under DIFFERENT keys. The account's own encryption settings cover
everything by default; a scope overrides that for the blobs and
containers that opt into it -- the multi-tenant pattern (one account,
per-tenant keys) and the mixed-sensitivity pattern (platform-managed
keys for ordinary data, a customer-managed key for the regulated
subset) without splitting into per-tenant accounts.

Scopes are many-per-account with independent lifecycles and are
referenced by name from containers (default_encryption_scope) and
ADLS filesystems, which is why they are a first-class kind rather
than a list folded into the account's spec. The parent is fixed at
creation: a scope cannot move between accounts.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageEncryptionScope
metadata:
  name: test-storage-encryption-scope
spec:
  storageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackstorage
  scopeName: hackTenantScope
  # Exercises the Key Vault source enum mapping AND the key pairing --
  # the deeper of the two source paths.
  source: MICROSOFT_KEY_VAULT
  keyVaultKeyId:
    value: https://plantonhackvault.vault.azure.net/keys/hack-tenant-key
  # Exercises the double-encryption flag (sent only when true).
  infrastructureEncryptionRequired: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.scopeName` | `string` | yes |  |  |
| `spec.source` | `enum` |  |  |  |
| `spec.keyVaultKeyId` | `string \| valueFrom` |  |  | AzureKeyVaultKey (`status.outputs.versionless_id`) |
| `spec.infrastructureEncryptionRequired` | `bool` |  |  |  |

## Field Details

### spec.storageAccountId

`string | valueFrom` · required

The storage account the scope lives in, by ARM ID. References an
AzureStorageAccount's storage_account_id output so the account and
its scopes compose in one manifest set. Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.scopeName

`string` · required

The scope's name: 4-63 letters and digits (hyphens are NOT allowed
-- stricter than most storage names). Unique within the account;
containers and blobs reference the scope by this name. Changing the
name replaces the scope.

- rule: scope_name must be 4-63 letters and digits (no hyphens)
- rule: {"required":true}

### spec.source

`enum`

Who owns the scope's encryption key. MICROSOFT_STORAGE uses a
platform-managed key Azure creates and rotates; MICROSOFT_KEY_VAULT
encrypts under YOUR Key Vault key (customer-managed) -- which then
requires key_vault_key_id, and the account must carry an identity
with wrap/unwrap access on the key's vault (the same plumbing as
the account-level customer_managed_key).

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_storage_encryption_scope_source_unspecified` -- Not specified -- invalid; choose an explicit key ownership model.
- `MICROSOFT_STORAGE` -- A platform-managed key Azure creates and rotates -- zero key management, no Key Vault involved.
- `MICROSOFT_KEY_VAULT` -- A customer-managed key in YOUR Key Vault -- you control rotation, revocation, and access; requires key_vault_key_id and an account identity with wrap/unwrap access on the vault.

### spec.keyVaultKeyId

`string | valueFrom`

The Key Vault key the scope encrypts under -- required when (and
only meaningful when) source is MICROSOFT_KEY_VAULT. References an
AzureKeyVaultKey's VERSIONLESS ID so key rotation propagates to the
scope with zero intervention; pin a versioned key URI only when a
compliance regime demands a frozen version.

- references: AzureKeyVaultKey (`status.outputs.versionless_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

### spec.infrastructureEncryptionRequired

`bool`

Whether data under this scope is DOUBLE-encrypted: a second,
independent layer with platform-managed keys underneath the scope's
key -- for regimes that require defense against a single-algorithm
or single-key compromise. Independent of the ACCOUNT's
infrastructure-encryption switch, so a scope can add the second
layer for just the regulated subset. Fixed at creation.

## Validation Rules

- `storage_encryption_scope_key_vault_key_required`: key_vault_key_id is required when source is MICROSOFT_KEY_VAULT

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageEncryptionScope, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.encryption_scope_id` | `string` | The Azure Resource Manager ID of the encryption scope. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{account}/encryptionScopes/{name} |
| `status.outputs.encryption_scope_name` | `string` | The scope's name -- what containers (default_encryption_scope), ADLS filesystems, and per-blob upload options reference within the account. |
| `status.outputs.storage_account_name` | `string` | The name of the storage account the scope lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/scope pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.keyVaultKeyId` | AzureKeyVaultKey | `status.outputs.versionless_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureStorageContainer | `spec.defaultEncryptionScope` | `status.outputs.encryption_scope_name` |
| AzureStorageDataLakeGen2Filesystem | `spec.defaultEncryptionScope` | `status.outputs.encryption_scope_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
