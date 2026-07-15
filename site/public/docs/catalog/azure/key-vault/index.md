---
title: "Key Vault"
description: "Key Vault deployment documentation"
icon: "package"
order: 100
componentName: "azurekeyvault"
---

# Azure Key Vault

Creates an Azure Key Vault -- the tenant-scoped container where encryption keys, TLS certificates, and application secrets live behind one security boundary. The vault sets the governance every object inside inherits: authorization mode, network isolation, deletion safety, and the pricing tier that gates HSM-backed keys.

## What Gets Created

When you deploy an AzureKeyVault resource, Planton provisions:

- **Key Vault** — an `azurerm_key_vault` in the specified region and resource group, with your chosen authorization mode, network rules, deletion-safety posture, and tags

What lives inside the vault is deliberately composed, never bundled: encryption keys are `AzureKeyVaultKey` resources, TLS certificates are `AzureKeyVaultCertificate` resources, and secret values belong to a secrets-management workflow rather than infrastructure manifests.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the vault in (an `AzureResourceGroup` in composed environments)
- **Key Vault write rights**: `Microsoft.KeyVault/vaults/write` (Key Vault Contributor, Contributor, or Owner)

## Quick Start

Create a file `key-vault.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVault
metadata:
  name: platform-vault
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureKeyVault.platform-vault
spec:
  region: eastus
  resourceGroup:
    value: security-rg
  vaultName: myorg-platform-kv
```

Deploy:

```shell
planton apply -f key-vault.yaml
```

After deployment, read `status.outputs.key_vault_id` for downstream wiring and `status.outputs.vault_uri` for the endpoint applications call.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region. Changing it replaces the vault. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `vaultName` | `string` | The vault's name -- GLOBALLY unique across Azure (it becomes `{name}.vault.azure.net`). | Required, 3-24 chars, letter start, letter/digit end, no consecutive hyphens |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `sku` | `enum` | `STANDARD` | `STANDARD` or `PREMIUM` (HSM-backed key types for the keys inside; FIPS 140-2 Level 3). Updatable in place. |
| `rbacAuthorizationEnabled` | `bool` | `true` | Azure RBAC (recommended -- grants compose as `AzureRoleAssignment`) vs legacy access policies. |
| `accessPolicies` | `list` | `[]` | Legacy grants: principal object id (defaults to referencing a user-assigned identity's `principal_id`), optional tenant/application ids, and explicit key/secret/certificate/storage permission lists. Only honored when RBAC is off. |
| `enabledForDeployment` | `bool` | `false` | Azure VMs may retrieve certificates stored as secrets. |
| `enabledForDiskEncryption` | `bool` | `false` | Azure Disk Encryption may retrieve secrets and unwrap keys (legacy in-guest ADE). |
| `enabledForTemplateDeployment` | `bool` | `false` | ARM template deployments may retrieve secrets. |
| `publicNetworkAccessEnabled` | `bool` | `true` | `false` takes the vault fully private (private endpoints only). |
| `purgeProtectionEnabled` | `bool` | `false` | Deleted objects can only be recovered, never purged, until retention expires. Irreversible once on. Turn on for production and CMK vaults. |
| `softDeleteRetentionDays` | `int32` | `90` | 7-90 days; fixed at creation. |
| `networkAcls` | `object` | -- | Public-endpoint firewall: `defaultAction` (ALLOW/DENY), `bypass` (AZURE_SERVICES/NONE), `ipRules`, `virtualNetworkSubnetIds` (references to `AzureSubnet` outputs; subnets need the Microsoft.KeyVault service endpoint). |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Production Vault with Purge Protection and Network Rules

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVault
metadata:
  name: prod-vault
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureKeyVault.prod-vault
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: security-rg
  vaultName: myorg-prod-kv
  purgeProtectionEnabled: true
  networkAcls:
    defaultAction: DENY
    ipRules:
      - "203.0.113.0/24"
    virtualNetworkSubnetIds:
      - valueFrom:
          name: app-subnet
  tags:
    cost-center: "1234"
```

### Legacy Access-Policy Vault

For orgs that have not yet moved to RBAC:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVault
metadata:
  name: legacy-vault
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: security-rg
  vaultName: myorg-legacy-kv
  rbacAuthorizationEnabled: false
  accessPolicies:
    - objectId:
        valueFrom:
          name: app-identity
      keyPermissions:
        - KEY_GET
        - KEY_UNWRAP_KEY
        - KEY_WRAP_KEY
      secretPermissions:
        - SECRET_GET
```

### Grant Access in RBAC Mode

Grants are ordinary role assignments scoped at the vault:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: app-reads-secrets
spec:
  scope:
    valueFrom:
      kind: AzureKeyVault
      name: prod-vault
      fieldPath: status.outputs.key_vault_id
  roleDefinitionName: Key Vault Secrets User
  principalId:
    valueFrom:
      name: app-identity
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `key_vault_id` | `string` | The vault's ARM resource ID -- referenced by `AzureKeyVaultKey`, `AzureKeyVaultCertificate`, VM/VMSS secret blocks, and vault-scoped role assignments |
| `key_vault_name` | `string` | The vault's name |
| `vault_uri` | `string` | The data-plane URI (`https://{name}.vault.azure.net/`) applications call |
| `tenant_id` | `string` | The Azure AD tenant the vault authenticates against |
| `resource_group_name` | `string` | The resource group the vault was created in |

## Related Components

- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — encryption keys inside the vault (the customer-managed-key story)
- [AzureKeyVaultCertificate](/docs/catalog/azure/key-vault-certificate) — TLS certificates the vault enrolls, renews, and guards
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — data-plane grants in RBAC mode
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the workload identities grants target
- [AzureSubnet](/docs/catalog/azure/subnet) — service-endpoint subnets for `networkAcls`
- [AzureResourceGroup](/docs/catalog/azure/resource-group) — provides the resource group for vault placement
