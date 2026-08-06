---
title: "Key Vault Certificate"
description: "Key Vault Certificate deployment documentation"
icon: "package"
order: 100
componentName: "azurekeyvaultcertificate"
---

# Azure Key Vault Certificate

Creates an X.509 certificate inside an Azure Key Vault -- enrolled by the vault (self-signed or via a CA) or imported from an existing PFX/PEM bundle, with the vault owning private-key custody, renewal, and expiry notification.

## What Gets Created

When you deploy an AzureKeyVaultCertificate resource, Planton provisions:

- **Key Vault certificate** — an `azurerm_key_vault_certificate` (data-plane object): the certificate, its private key, and its secret face -- the identifier TLS terminators like Application Gateway consume

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A Key Vault** to create the certificate in (an `AzureKeyVault` in composed environments)
- **Data-plane certificate permissions on the vault** for the deploying credential: the "Key Vault Administrator" or "Key Vault Certificates Officer" RBAC role (or certificate permissions in a legacy access policy). Subscription Owner alone is NOT enough -- certificates are data-plane objects.

## Quick Start

Create a file `key-vault-certificate.yaml` (self-signed, auto-renewing internal TLS):

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVaultCertificate
metadata:
  name: internal-tls
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureKeyVaultCertificate.internal-tls
spec:
  name: internal-tls
  keyVaultId:
    valueFrom:
      name: platform-vault
  certificatePolicy:
    issuerName: Self
    keyProperties:
      exportable: true
      keyType: RSA
      keySize: 2048
    lifetimeActions:
      - actionType: AUTO_RENEW
        trigger:
          lifetimePercentage: 80
    secretProperties:
      contentType: PKCS12
    x509CertificateProperties:
      subject: CN=internal.example.com
      subjectAlternativeNames:
        dnsNames:
          - internal.example.com
      keyUsage:
        - DIGITAL_SIGNATURE
        - KEY_ENCIPHERMENT
      validityInMonths: 12
```

Deploy:

```shell
planton apply -f key-vault-certificate.yaml
```

After deployment, read `status.outputs.versionless_secret_id` -- the reference TLS terminators consume, following renewals automatically.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | The certificate's name within the vault. | Required, 1-127 letters/digits/hyphens; fixed at creation |
| `keyVaultId` | `StringValueOrRef` | The vault, by ARM ID. Defaults to referencing an `AzureKeyVault`'s `key_vault_id` output. | Required; fixed at creation |

At least one of `certificate` (import) / `certificatePolicy` (generate) is required. A generated certificate additionally requires `certificatePolicy.x509CertificateProperties`.

### The Import Block (`certificate`)

| Field | Type | Description |
|-------|------|-------------|
| `contents` | `string` (sensitive) | Base64 PFX or PEM bundle containing BOTH the chain and the private key. Changing it imports a new version. |
| `password` | `string` (sensitive) | The bundle's password; omit for unprotected PEMs and passwordless PFX bundles. |

### The Policy Block (`certificatePolicy`)

| Field | Type | Description |
|-------|------|-------------|
| `issuerName` | `string` | `Self` (self-signed), `Unknown` (out-of-band CA via CSR), or the name of a CA issuer configured on the vault. |
| `keyProperties` | `object` | `exportable`, `keyType` (RSA/RSA_HSM/EC/EC_HSM/OCT), `keySize`, `curve`, `reuseKey` (renewals re-use the private key). |
| `lifetimeActions` | `list` | `actionType` (AUTO_RENEW / EMAIL_CONTACTS) + `trigger` with exactly one of `daysBeforeExpiry` / `lifetimePercentage`. |
| `secretProperties.contentType` | `enum` | `PKCS12` (Application Gateway, Windows world) or `PEM` (Linux/OpenSSL world). |
| `x509CertificateProperties` | `object` | `subject` (DN), `subjectAlternativeNames` (dnsNames/emails/upns), `keyUsage` (min 1), `extendedKeyUsage` OIDs, `validityInMonths`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Import an Existing PFX

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVaultCertificate
metadata:
  name: purchased-wildcard
spec:
  name: purchased-wildcard
  keyVaultId:
    valueFrom:
      name: platform-vault
  certificate:
    contents: <base64-pfx-bytes>
    password: <bundle-password>
```

### Terminate TLS on Application Gateway

```yaml
spec:
  sslCertificates:
    - name: internal-tls
      keyVaultSecretId:
        valueFrom:
          name: internal-tls
```

The gateway's managed identity needs secret-read on the vault (the "Key Vault Secrets User" role in RBAC mode).

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `certificate_id` | `string` | Versioned data-plane ID |
| `versionless_id` | `string` | Versionless data-plane ID -- follows renewals |
| `secret_id` | `string` | Versioned ID of the certificate's secret face |
| `versionless_secret_id` | `string` | Versionless secret-face ID -- what TLS terminators should reference |
| `certificate_name` | `string` | The certificate's name within the vault |
| `version` | `string` | The current version identifier |
| `thumbprint` | `string` | SHA-1 fingerprint, hex-encoded |
| `resource_manager_id` | `string` | Versioned ARM resource ID |
| `resource_manager_versionless_id` | `string` | Versionless ARM resource ID |

## Related Components

- [AzureKeyVault](/docs/catalog/azure/key-vault) — the vault the certificate lives in
- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the encryption-key sibling
- [AzureApplicationGateway](/docs/catalog/azure/application-gateway) — consumes the secret face for TLS listeners
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants the deployer/consumers data-plane certificate permissions
