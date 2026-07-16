# AzureKeyVault -- Design Research

## The Resource

An Azure Key Vault (`Microsoft.KeyVault/vaults`) is the tenant-scoped
container for an organization's keys, certificates, and secrets. The
component maps onto `azurerm_key_vault` (azurerm v4.x,
`internal/services/keyvault/key_vault_resource.go`), parity-verified against
pulumi-azure v6 (`keyvault.KeyVault`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `vault_name` | Required, ForceNew, globally unique (DNS name) |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `sku_name` | `sku` enum | standard / premium; updatable |
| `tenant_id` | -- (derived) | Read from the deploying credential's client config on both engines; a vault cannot be managed cross-tenant, so modeling it would only invite contradiction. The `tenant_id` OUTPUT reports it. |
| `rbac_authorization_enabled` | `rbac_authorization_enabled` | azurerm v4 canonical name; spec default TRUE (Azure's recommendation) while ARM's create default is false -- deliberate, documented in the field comment |
| `access_policy` | `access_policies` | Full legacy mode: object id (FK → UAI principal id), tenant (defaults to the vault's own), application id, four closed permission enums |
| `enabled_for_deployment` / `_disk_encryption` / `_template_deployment` | same | Plain bools, Azure defaults false |
| `public_network_access_enabled` | same | optional bool, default true |
| `purge_protection_enabled` | same | Plain bool -- azurerm's real default is FALSE (the prior spec's invented always-true default removed) |
| `soft_delete_retention_days` | same | 7-90, default 90, ForceNew |
| `network_acls` | `network_acls` | default_action required-in-block; `bypass` a closed enum (AZURE_SERVICES / NONE) with unspecified = Azure's default (AzureServices); subnet ids are FK references → AzureSubnet |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `vault_uri` (computed) | `vault_uri` output | The data-plane endpoint |

## Decomposition Decisions

- **Keys and certificates are first-class children**
  (`AzureKeyVaultKey` / `AzureKeyVaultCertificate`), each FK-referencing
  `key_vault_id`: independent lifecycle, many-per-vault, and FK-referenced
  by CMK/TLS consumers. Hiding the vault behind them was considered and
  rejected -- which vault an object lives in is a governance decision
  (network isolation, RBAC mode, purge protection, HSM tier), the vault is
  itself FK-referenced, and container-plus-children is the established
  catalog pattern.
- **No secret-creation surface.** Secret values are runtime configuration,
  not infrastructure -- IaC-managed secret entries put plaintext (or empty
  placeholders that drift) into manifests and state. The vault is
  provisioned here; secret content belongs to a secrets-management workflow.
- **A standalone access-policy kind was folded away**: RBAC is the modern
  grain (grants are `AzureRoleAssignment` resources already), and azurerm
  itself warns that mixing inline and standalone policies on one vault
  causes perpetual drift -- the vault owns its complete legacy grant list
  inline.

## Recorded Skips (with reasons)

- **`contact`** -- deprecated in azurerm v4 (update-only; new vaults with it
  error) and removed in v5; certificate contacts belong to the
  data-plane contacts resource, evaluated below.
- **`azurerm_key_vault_certificate_contacts` / `_certificate_issuer` /
  `_managed_storage_account` (+ SAS definitions)** -- niche data-plane
  side-resources (expiry-notification contacts, integrated-CA issuer
  configs, the legacy storage-key-rotation feature). None is FK-referenced
  or chart-demanded; they join the adoption backlog rather than shipping as
  half-used kinds.
- **`azurerm_key_vault_access_policy` (standalone)** -- folded inline (see
  above).

## Design Decisions

- **RBAC default true vs ARM's false.** Azure's own guidance, the Planton
  composition story (`AzureRoleAssignment`), and Microsoft's portal default
  for new vaults all point to RBAC; the spec default encodes the
  recommended posture while the field stays explicit for legacy-mode orgs.
  The divergence from ARM's create default is called out in the field
  comment.
- **`purge_protection_enabled` conformed to azurerm's real default
  (false).** The prior spec defaulted it true -- an invented default that
  made every dev vault un-purgeable for the retention window. Presets teach
  enabling it for production and CMK vaults instead.
- **`bypass` promoted from bool to a closed enum** (AZURE_SERVICES / NONE)
  -- provider-authentic vocabulary instead of a bool that hid the ARM
  value names.
- **Access-policy `object_id` is a StringValueOrRef** defaulting to an
  `AzureUserAssignedIdentity`'s `principal_id` output -- the same
  composition seam the role-assignment kind established, so both
  authorization modes compose.
- **Permission lists are closed enums** (20 key / 8 secret / 16 certificate
  / 14 storage values) with exhaustive vocabulary maps in both engines -- a
  missing entry would silently drop a grant, so the maps are complete by
  construction.

## Operational Behavior Worth Knowing

- **Soft delete is always on** (ARM removed the opt-out): destroy makes the
  vault name unavailable for the retention window unless purged. Both
  engines' provider defaults purge on destroy; with purge protection ON,
  destroy becomes a scheduled deletion and the name stays reserved until the
  window ends.
- **A name-colliding create against a soft-deleted vault auto-recovers it**
  (provider default `recover_soft_deleted_key_vaults = true`) -- re-running
  a deployment after a destroy inside the retention window resurrects the
  old vault rather than failing.
- **Many CMK integrations refuse vaults without purge protection** (disk
  encryption sets, several database CMK flows) -- the preset for CMK vaults
  turns it on.
- **RBAC data-plane grants take time to propagate** (typically under two
  minutes, occasionally longer) -- a key/certificate create immediately
  after granting the deployer can 403 transiently.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `network_acls.virtual_network_subnet_ids[]` →
  `AzureSubnet.status.outputs.subnet_id` (subnets need the
  `Microsoft.KeyVault` service endpoint)
- `access_policies[].object_id` →
  `AzureUserAssignedIdentity.status.outputs.principal_id`
- `key_vault_id` output is consumed by:
  - `AzureKeyVaultKey.key_vault_id` / `AzureKeyVaultCertificate.key_vault_id`
  - `AzureAksCluster.service_mesh_profile.certificate_authority.key_vault_id`
  - `AzureVirtualMachine` / `AzureVirtualMachineScaleSet` os-profile secrets
    (`key_vault_id`, `source_vault_id`)
  - vault-scoped `AzureRoleAssignment.scope`
- `vault_uri` output is the endpoint applications and Config Manager's
  Azure secret backend call
