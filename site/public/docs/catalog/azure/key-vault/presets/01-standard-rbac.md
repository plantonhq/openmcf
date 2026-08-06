---
title: "Standard RBAC Vault"
description: "This preset creates the baseline modern vault: Standard SKU, Azure RBAC authorization (the spec default), public endpoint, and Azure's own deletion-safety defaults (90-day soft delete, purge..."
type: "preset"
rank: "01"
presetSlug: "01-standard-rbac"
componentSlug: "key-vault"
componentTitle: "Key Vault"
provider: "azure"
icon: "package"
order: 1
---

# Standard RBAC Vault

This preset creates the baseline modern vault: Standard SKU, Azure RBAC
authorization (the spec default), public endpoint, and Azure's own
deletion-safety defaults (90-day soft delete, purge protection off so
non-production vaults can be cleanly destroyed and their names reused).

Data-plane access is granted by composing `AzureRoleAssignment` resources
scoped at the vault -- nothing about who-can-do-what lives inside the vault
spec in RBAC mode.

## When to Use

- The starting point for almost every new vault -- development, staging,
  and any production vault that does not yet hold customer-managed keys
- Teams standardizing on Azure RBAC for data-plane authorization (PIM,
  access reviews, fine-grained scopes)

## Key Configuration Choices

- **RBAC on (default)** -- grants compose as `AzureRoleAssignment` with
  roles like "Key Vault Administrator", "Key Vault Secrets User", "Key
  Vault Crypto Officer"
- **Purge protection off** -- destroy purges the vault so the globally
  unique name frees up immediately; move to the production preset before
  the vault holds CMKs
- **`vault_name` is a global DNS name** -- it becomes
  `{name}.vault.azure.net`, so pick something org-prefixed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the vault in | The resource group's `status.outputs.resource_group_name` |
| `myorg-platform-vault` | 3-24 chars, letters/digits/hyphens, globally unique | Your naming convention (org prefix recommended) |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Grant an identity read access to secrets:

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
