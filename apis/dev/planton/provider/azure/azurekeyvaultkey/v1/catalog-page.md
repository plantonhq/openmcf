# Azure Key Vault Key

Creates a cryptographic key inside an Azure Key Vault -- the customer-managed-key (CMK) building block. The private part never leaves the vault; Storage accounts, disk encryption sets, container registries, and databases encrypt their data with a data-encryption key that this key wraps.

## What Gets Created

When you deploy an AzureKeyVaultKey resource, Planton provisions:

- **Key Vault key** — an `azurerm_key_vault_key` (data-plane object) with your chosen algorithm family, capability list, and optional automatic-rotation policy

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A Key Vault** to create the key in (an `AzureKeyVault` in composed environments); the `_HSM` key types need the vault's `PREMIUM` SKU
- **Data-plane key permissions on the vault** for the deploying credential: the "Key Vault Administrator" or "Key Vault Crypto Officer" RBAC role (or key permissions in a legacy access policy). Subscription Owner alone is NOT enough -- keys are data-plane objects.

## Quick Start

Create a file `key-vault-key.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVaultKey
metadata:
  name: storage-cmk
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureKeyVaultKey.storage-cmk
spec:
  name: storage-cmk
  keyVaultId:
    valueFrom:
      name: platform-vault
  keyType: RSA
  keySize: 2048
  keyOpts:
    - WRAP_KEY
    - UNWRAP_KEY
```

Deploy:

```shell
planton apply -f key-vault-key.yaml
```

After deployment, read `status.outputs.versionless_id` -- the reference CMK consumers should use so rotation propagates automatically.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | The key's name within the vault. | Required, 1-127 letters/digits/hyphens; fixed at creation |
| `keyVaultId` | `StringValueOrRef` | The vault, by ARM ID. Defaults to referencing an `AzureKeyVault`'s `key_vault_id` output. | Required; fixed at creation |
| `keyType` | `enum` | `RSA`, `RSA_HSM`, `EC`, or `EC_HSM`. RSA is the universal CMK choice; EC keys sign/verify. `_HSM` variants require the vault's PREMIUM SKU. | Required; fixed at creation |
| `keyOpts` | `enum[]` | The operations the key may perform: `DECRYPT`, `ENCRYPT`, `SIGN`, `UNWRAP_KEY`, `VERIFY`, `WRAP_KEY`. CMK needs WRAP_KEY + UNWRAP_KEY. | Required, at least one |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `keySize` | `int32` | -- | RSA modulus bits: `2048`, `3072`, or `4096`. Required for RSA keys, forbidden for EC. Fixed at creation. |
| `curve` | `enum` | `P_256` | EC curve: `P_256`, `P_256K`, `P_384`, `P_521`. EC keys only. Fixed at creation. |
| `notBeforeDate` | `string` | -- | RFC 3339 UTC instant before which the key must not be used. |
| `expirationDate` | `string` | -- | RFC 3339 UTC expiry. Prefer a rotation policy for keys that encrypt live data -- an expired CMK takes its dependents down with it. |
| `rotationPolicy` | `object` | -- | `expireAfter` + `notifyBeforeExpiry` (ISO 8601 durations, set together) and/or an `automatic` block with `timeAfterCreation` or `timeBeforeExpiry`. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Auto-Rotating CMK

Rotate every ~11 months with 2-year per-version expiry and 30-day notice:

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
  rotationPolicy:
    expireAfter: P2Y
    notifyBeforeExpiry: P30D
    automatic:
      timeAfterCreation: P335D
```

### EC Signing Key

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVaultKey
metadata:
  name: jwt-signing
spec:
  name: jwt-signing
  keyVaultId:
    valueFrom:
      name: platform-vault
  keyType: EC
  curve: P_384
  keyOpts:
    - SIGN
    - VERIFY
```

### Wire to a CMK Consumer

Container-registry encryption follows rotation through the versionless reference:

```yaml
spec:
  encryption:
    identityClientId: <identity-client-id>
    keyVaultKeyId:
      valueFrom:
        name: registry-cmk
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `key_id` | `string` | Versioned data-plane ID -- pins consumers to this version |
| `versionless_id` | `string` | Versionless data-plane ID -- the CMK reference; rotation propagates automatically |
| `key_name` | `string` | The key's name within the vault |
| `version` | `string` | The current version identifier |
| `resource_id` | `string` | Versioned ARM resource ID (control-plane identity) |
| `resource_versionless_id` | `string` | Versionless ARM resource ID |
| `public_key_pem` | `string` | The public half in PEM form |
| `public_key_openssh` | `string` | The public half in OpenSSH form |

## Related Components

- [AzureKeyVault](/docs/catalog/azure/key-vault) — the vault the key lives in (Premium SKU for HSM types)
- [AzureKeyVaultCertificate](/docs/catalog/azure/key-vault-certificate) — the TLS sibling
- [AzureContainerRegistry](/docs/catalog/azure/container-registry) — CMK consumer via `encryption.keyVaultKeyId`
- [AzureAksCluster](/docs/catalog/azure/aks-cluster) — KMS etcd encryption via `keyManagementService.keyVaultKeyId` (pins versions)
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants the deployer/consumers data-plane key permissions
