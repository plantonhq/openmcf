# AzureKeyVaultKey -- Design Research

## The Resource

A Key Vault key is a data-plane object inside a vault
(`https://{vault}.vault.azure.net/keys/{name}`) with a read-only ARM proxy
(`Microsoft.KeyVault/vaults/keys`). The component maps onto
`azurerm_key_vault_key` (azurerm v4.x,
`internal/services/keyvault/key_vault_key_resource.go`), parity-verified
against pulumi-azure v6 (`keyvault.Key`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew, 1-127 alphanumeric/hyphens |
| `key_vault_id` | `key_vault_id` | FK → AzureKeyVault `key_vault_id` output; ForceNew |
| `key_type` | `key_type` enum | RSA / RSA-HSM / EC / EC-HSM (hyphenated on the wire); ForceNew |
| `key_size` | `key_size` | 2048/3072/4096; RSA-only (CEL-paired); ForceNew |
| `curve` | `curve` enum | P-256 / P-256K / P-384 / P-521; EC-only (CEL-paired); unset lets Azure default P-256; ForceNew |
| `key_opts` | `key_opts` enum list | Required, min 1; camelCase on the wire (case-sensitive) |
| `not_before_date` / `expiration_date` | same | RFC 3339, pattern-validated |
| `rotation_policy` | `rotation_policy` | expire_after/notify_before_expiry paired (azurerm's RequiredWith) + automatic trigger with at-least-one (azurerm's AtLeastOneOf), all as CELs |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `version` / `versionless_id` / `resource_id` / `resource_versionless_id` / `public_key_pem` / `public_key_openssh` (computed) | outputs | The composition surface |

## Recorded Skips (with reasons)

- **`release_policy` (secure key release)** -- azurerm models it (JSON
  policy + immutable flag, HSM-only, confidential-computing scenarios), but
  the pinned pulumi-azure v6 SDK carries no `releasePolicy` on
  `keyvault.Key` **as of the latest v6 release** -- no SDK bump can close
  the gap. Modeling a field only Terraform can deploy would break the
  100% behavioral-parity invariant (silent drop or deploy-time error on
  Pulumi). The field is narrow (requires RSA-HSM/EC-HSM on a Premium
  vault and a trusted-execution-environment consumer); it joins the
  adoption backlog, enabled the moment pulumi-azure gains `releasePolicy`.

## Design Decisions

- **`key_size` XOR `curve` enforced as CELs mirroring azurerm's contract**:
  RSA requires a size (Azure has no RSA default), EC forbids one; curve is
  EC-only and optional (Azure defaults P-256). Catching the mismatch at
  validation time beats a mid-apply data-plane error.
- **`key_opts` is a closed enum list with min 1** -- azurerm marks it
  Required; the six operations are Azure's complete, stable vocabulary and
  the list is the key's capability boundary.
- **Rotation-policy pairings as CELs**: a policy needs `expire_after`
  and/or `automatic` (an empty policy configures nothing);
  `expire_after`/`notify_before_expiry` go together (azurerm's
  RequiredWith); the automatic block needs at least one trigger. ISO 8601
  duration patterns validated at the spec level.
- **Outputs export BOTH the versioned and versionless identities** at
  data-plane and ARM levels: `versionless_id` is the CMK grain (rotation
  propagates); `key_id` pins a version (AKS KMS pins versions by design);
  the ARM proxy ids serve control-plane integrations and the E2E verifier.
- **Wire-format vocabularies are exhaustive maps in both engines** (proto
  enum name → Azure's case-sensitive string): RSA_HSM → `RSA-HSM`,
  P_256K → `P-256K`, UNWRAP_KEY → `unwrapKey`. A missing entry would
  silently drop a capability, so the maps cover every enum value by
  construction.

## Operational Behavior Worth Knowing

- **Data-plane creation**: the provider talks to the vault endpoint, not
  ARM -- the deploying credential needs key permissions on the vault
  (Key Vault Administrator / Crypto Officer, or an access policy);
  subscription Owner alone 403s.
- **Key material is immutable**: type/size/curve changes replace the key;
  every CMK consumer then re-wraps through the new key.
- **Rotation creates versions**: consumers referencing `versionless_id`
  follow automatically; ACR-style consumers re-read on their next unwrap.
- **Soft delete + purge**: a deleted key's name stays reserved in the
  vault for the retention window; the providers purge soft-deleted keys on
  destroy by default so re-creates work.
- **Expiry is sticky**: once set, `expiration_date` cannot be fully unset
  even across delete/recreate -- Azure restores the purged name's state.

## Composition

- `key_vault_id` → `AzureKeyVault.status.outputs.key_vault_id`
- `versionless_id` output is consumed by:
  - `AzureContainerRegistry.encryption.key_vault_key_id` (CMK; rotation
    propagates)
- `key_id` output (versioned) is consumed by:
  - `AzureAksCluster.key_management_service.key_vault_key_id` (AKS KMS
    pins key versions; rotation = update to the new version's id)
- The deployer's data-plane grant composes as `AzureRoleAssignment`
  ("Key Vault Crypto Officer") scoped at the vault's `key_vault_id`
