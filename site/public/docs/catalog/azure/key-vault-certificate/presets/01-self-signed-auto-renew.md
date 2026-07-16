---
title: "Self-Signed Auto-Renewing Certificate"
description: "This preset is the fully-hands-off shape for internal TLS: the vault generates the key, self-signs a 12-month certificate, and automatically re-enrolls at 80% of lifetime -- no CA account, no renewal..."
type: "preset"
rank: "01"
presetSlug: "01-self-signed-auto-renew"
componentSlug: "key-vault-certificate"
componentTitle: "Key Vault Certificate"
provider: "azure"
icon: "package"
order: 1
---

# Self-Signed Auto-Renewing Certificate

This preset is the fully-hands-off shape for internal TLS: the vault
generates the key, self-signs a 12-month certificate, and automatically
re-enrolls at 80% of lifetime -- no CA account, no renewal calendar, no
human in the loop, ever.

Consumers referencing the versionless outputs
(`versionless_secret_id` for TLS terminators) roll onto each renewal
automatically.

## When to Use

- Internal service-to-service TLS where consumers trust the certificate
  explicitly (pinned thumbprint or private trust store)
- Development and staging environments standing in for a CA-issued
  production certificate
- Application Gateway backends and internal listeners that need *a* valid
  certificate rather than a publicly-trusted one

Public-facing endpoints need a CA-issued certificate instead -- import one
with the imported-bundle preset or configure a CA issuer on the vault.

## Key Configuration Choices

- **`issuerName: Self`** -- zero external dependencies; enrollment
  completes synchronously
- **`lifetimePercentage: 80`** -- the conventional renewal point (~2.4
  months before a 12-month certificate expires)
- **`contentType: PKCS12`** -- what Application Gateway and Windows-world
  consumers expect; switch to `PEM` for the Linux/OpenSSL world
- **TLS-server key usages** -- `DIGITAL_SIGNATURE` + `KEY_ENCIPHERMENT`,
  the conventional pair; modern clients validate SANs, so every hostname
  goes in `dnsNames` with the subject CN as the primary

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | The vault's ARM resource ID (or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<primary-hostname>` | The hostname the certificate secures | Your service's DNS name |

## Downstream Wiring

```yaml
spec:
  sslCertificates:
    - name: internal-tls
      keyVaultSecretId:
        valueFrom:
          name: internal-tls
```
