# AzureKeyVaultCertificate

## Overview

`AzureKeyVaultCertificate` provisions an X.509 certificate inside an Azure
Key Vault -- the TLS building block the vault manages end to end: private
key custody, enrollment or import, renewal, and expiry notification.

A vault certificate is really three linked objects sharing one name: the
certificate (public part + policy), its private KEY, and a SECRET whose
value is the full bundle. That third face is the composition seam -- Azure
services that terminate TLS (Application Gateway listeners, App Service
custom domains) consume the certificate THROUGH its secret identifier,
which is why the `secret_id` / `versionless_secret_id` outputs exist and
are what downstream kinds reference.

## Two Ways In (Combinable)

- **Generate** (`certificate_policy` only): the vault creates the key and
  either self-signs (issuer `Self`) or forwards a CSR to a configured CA.
  Self-signed + auto-renew is the fully-hands-off shape for internal TLS.
- **Import** (`certificate`): bring an existing PFX/PEM bundle; a policy
  may accompany it to govern the imported material. Without an explicit
  policy Azure derives one from the bundle.

## Key Features

- **Self-signed, CSR-pending, and integrated-CA issuance** -- issuer
  `Self`, `Unknown`, or the name of a CA issuer configured on the vault
  (DigiCert / GlobalSign first-party integrations)
- **Automatic renewal** -- lifetime actions fire at a days-before-expiry
  or percentage-of-lifetime trigger: `AUTO_RENEW` re-enrolls,
  `EMAIL_CONTACTS` notifies the vault's certificate contacts
- **Full X.509 shape** -- subject, DNS/email/UPN SANs, key-usage and
  extended-key-usage extensions, validity in months
- **RSA and EC keys, exportable or vault-locked** -- HSM variants on
  Premium vaults; `exportable: false` keeps the private key inside the
  vault for consumers that fetch through Key Vault references
- **PKCS12 or PEM secret face** -- what TLS terminators reading the secret
  get

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | 1-127 letters/digits/hyphens, unique among the vault's certificates; ForceNew |
| `key_vault_id` | StringValueOrRef | Yes | -- | The vault, by ARM ID (defaults to an AzureKeyVault reference); ForceNew |
| `certificate` | message | One of | -- | Import: base64 `contents` (sensitive) + optional `password` (sensitive) |
| `certificate_policy` | message | One of | -- | Generate/govern: issuer, key properties, lifetime actions, secret content type, X.509 properties |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

At least one of `certificate` / `certificate_policy` is required; a
generated certificate additionally requires
`certificate_policy.x509_certificate_properties`.

## Outputs

| Output | Description |
|--------|-------------|
| `certificate_id` / `versionless_id` | The certificate's data-plane IDs (versioned / renewal-following) |
| `secret_id` / `versionless_secret_id` | The secret face's IDs -- what Application Gateway and App Service consume |
| `certificate_name` / `version` | Name and current version |
| `thumbprint` | SHA-1 fingerprint integrations match on |
| `resource_manager_id` / `resource_manager_versionless_id` | ARM (control-plane) IDs |

## Quick Example

Self-signed, auto-renewing internal TLS:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureKeyVaultCertificate
metadata:
  name: internal-tls
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

Wire it to an Application Gateway TLS listener:

```yaml
spec:
  sslCertificates:
    - name: internal-tls
      keyVaultSecretId:
        valueFrom:
          name: internal-tls
```

## Lifecycle Notes

- **The deploying credential needs data-plane certificate permissions on
  the vault** -- subscription Owner alone is not enough; the full deployer
  grant set is cataloged in [`iac/permissions.yaml`](iac/permissions.yaml)
- **Renewals and re-imports create new VERSIONS** -- consumers on the
  versionless references follow automatically
- Changing any policy part except `lifetime_actions` creates a new
  certificate version; lifetime actions update in place
- A deleted certificate's name stays reserved for the vault's soft-delete
  retention window unless purged (the IaC engines purge on destroy by
  default)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
