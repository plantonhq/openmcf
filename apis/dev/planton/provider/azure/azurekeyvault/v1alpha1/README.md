# AzureKeyVault

## Overview

`AzureKeyVault` provisions an Azure Key Vault: the tenant-scoped container
where an organization's encryption keys, TLS certificates, and application
secrets live behind one security boundary. The vault is where governance is
set -- authorization mode (Azure RBAC vs legacy access policies), network
isolation, deletion safety (soft delete and purge protection), and the
pricing tier that gates HSM-backed keys.

## The Three-Component Shape

The vault is a container; what lives inside it is composed, never bundled:

- **`AzureKeyVaultKey`** -- encryption keys (the customer-managed-key story),
  referencing this vault's `key_vault_id` output
- **`AzureKeyVaultCertificate`** -- TLS certificates the vault enrolls,
  renews, and guards, referencing the same output
- **Secret VALUES are deliberately out of scope** for infrastructure-as-code:
  provision the vault here and manage secret content through a
  secrets-management workflow, so plaintext never enters deployment manifests
  or state

Azure orgs deliberately run FEW vaults with MANY objects inside, because
every vault-level control (network rules, RBAC mode, purge protection, HSM
tier) applies to everything in it -- which vault an object lives in is a
governance decision, and that is exactly why the vault is a first-class,
referenceable node.

## Key Features

- **Azure RBAC by default** -- data-plane grants become ordinary
  `AzureRoleAssignment` resources ("Key Vault Administrator", "Key Vault
  Secrets User", ...) with PIM and access-review support; the legacy
  access-policy mode is fully modeled for orgs that still run it
- **Network isolation spectrum** -- public, public-with-firewall
  (`network_acls`: IP allowlists, VNet service-endpoint subnets, trusted
  Microsoft services bypass), or fully private
  (`public_network_access_enabled: false` + private endpoints)
- **Deletion safety** -- soft delete (7-90 days, fixed at creation) and
  optional purge protection (irreversible; required by many CMK
  integrations)
- **STANDARD / PREMIUM SKUs** -- Premium unlocks HSM-backed key types
  (FIPS 140-2 Level 3) for the keys inside
- **Deployment integrations** -- the three resource-manager switches
  (`enabled_for_deployment` / `_disk_encryption` / `_template_deployment`)
- **User tags** merged over the Planton-derived governance tags

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `region` | string | Yes | -- | Azure region; changing it replaces the vault |
| `resource_group` | StringValueOrRef | Yes | -- | Resource group name (defaults to an AzureResourceGroup reference) |
| `vault_name` | string | Yes | -- | 3-24 chars, letters/digits/hyphens, GLOBALLY unique (becomes `{name}.vault.azure.net`) |
| `sku` | enum | No | STANDARD | `STANDARD` or `PREMIUM` (HSM-backed keys); updatable in place |
| `rbac_authorization_enabled` | bool | No | true | Azure RBAC (recommended) vs legacy access policies |
| `access_policies` | list | No | -- | Legacy grants (object id + permission lists); only honored when RBAC is off |
| `enabled_for_deployment` | bool | No | false | VMs may retrieve certificates stored as secrets |
| `enabled_for_disk_encryption` | bool | No | false | Azure Disk Encryption may unwrap keys (legacy ADE) |
| `enabled_for_template_deployment` | bool | No | false | ARM template deployments may retrieve secrets |
| `public_network_access_enabled` | bool | No | true | false = private-endpoints-only |
| `purge_protection_enabled` | bool | No | false | Irreversible once on; recommended for production and CMK vaults |
| `soft_delete_retention_days` | int32 | No | 90 | 7-90; fixed at creation |
| `network_acls` | message | No | -- | default_action (ALLOW/DENY), bypass (AZURE_SERVICES/NONE), ip_rules, subnet references |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `key_vault_id` | The vault's ARM resource ID -- what `AzureKeyVaultKey`, `AzureKeyVaultCertificate`, and vault-scoped `AzureRoleAssignment` grants reference |
| `key_vault_name` | The vault's name |
| `vault_uri` | The data-plane URI (`https://{name}.vault.azure.net/`) applications call |
| `tenant_id` | The Azure AD tenant the vault authenticates against |
| `resource_group_name` | The resource group the vault was created in |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVault
metadata:
  name: platform-vault
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: security-rg
  vaultName: mycompany-platform-kv
  purgeProtectionEnabled: true
```

Grant an identity data-plane access (RBAC mode):

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRoleAssignment
metadata:
  name: app-reads-secrets
spec:
  scope:
    valueFrom:
      kind: AzureKeyVault
      name: platform-vault
      fieldPath: status.outputs.key_vault_id
  roleDefinitionName: Key Vault Secrets User
  principalId:
    valueFrom:
      name: app-identity
```

## Lifecycle Notes

- **The vault's name is a global DNS name.** A deleted vault's name stays
  reserved for the soft-delete retention window unless the vault is purged
  (the IaC engines purge on destroy by default -- unless purge protection is
  on, which turns destroy into a scheduled deletion at the end of the
  window).
- **Purge protection is a one-way door**: once enabled it cannot be
  disabled.
- `soft_delete_retention_days` can only be set at creation; changing it
  replaces the vault.
- Switching authorization modes on a live vault requires
  Microsoft.Authorization write permission (Owner / User Access
  Administrator) on the vault.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
