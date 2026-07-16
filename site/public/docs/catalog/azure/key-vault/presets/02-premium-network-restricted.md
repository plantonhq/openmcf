---
title: "Premium Network-Restricted CMK Vault"
description: "This preset creates the production posture for a vault holding customer-managed keys: Premium SKU (unlocks the HSM-backed key types for the `AzureKeyVaultKey` resources inside), purge protection ON..."
type: "preset"
rank: "02"
presetSlug: "02-premium-network-restricted"
componentSlug: "key-vault"
componentTitle: "Key Vault"
provider: "azure"
icon: "package"
order: 2
---

# Premium Network-Restricted CMK Vault

This preset creates the production posture for a vault holding
customer-managed keys: Premium SKU (unlocks the HSM-backed key types for
the `AzureKeyVaultKey` resources inside), purge protection ON (many CMK
integrations refuse to enroll against a vault without it), and a DENY
network firewall with explicit IP and subnet allowlists.

Trusted Microsoft services (Azure Backup, Disk Encryption, Azure Monitor)
keep access through the default bypass, so first-party integrations work
even behind the DENY posture.

## When to Use

- Vaults whose keys encrypt other resources (storage accounts, disk
  encryption sets, container registries, database CMK)
- Compliance regimes that mandate hardware key protection (PCI-DSS,
  FedRAMP High) or network isolation
- Any vault where accidental permanent deletion is unacceptable

## Key Configuration Choices

- **`sku: PREMIUM`** -- required before any key inside can use the
  RSA_HSM / EC_HSM types; the SKU itself can be changed in place later
- **`purgeProtectionEnabled: true`** -- a ONE-WAY door: once on it cannot
  be disabled, and destroying the vault schedules deletion for the end of
  the retention window instead of purging (the name stays reserved)
- **`defaultAction: DENY` + allowlists** -- the middle ground between
  "open to the internet" and "private endpoints only"; keep
  `publicNetworkAccessEnabled: true` for this shape
- **Subnets need the `Microsoft.KeyVault` service endpoint** enabled on
  the referenced `AzureSubnet`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the vault in | The resource group's `status.outputs.resource_group_name` |
| `myorg-prod-cmk-vault` | 3-24 chars, globally unique | Your naming convention |
| `<office-or-ci-cidr>` | Public IPv4 CIDR allowed through the firewall | Your office/VPN/CI egress range |
| `<app-subnet-arm-id>` | The subnet whose workloads reach the vault (or a `valueFrom` reference to an AzureSubnet's `subnet_id` output) | The subnet's `status.outputs.subnet_id` |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Put an HSM-backed CMK inside and wire it to a consumer:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVaultKey
metadata:
  name: registry-cmk
spec:
  name: registry-cmk
  keyVaultId:
    valueFrom:
      name: prod-cmk-vault
  keyType: RSA_HSM
  keySize: 2048
  keyOpts:
    - WRAP_KEY
    - UNWRAP_KEY
```
