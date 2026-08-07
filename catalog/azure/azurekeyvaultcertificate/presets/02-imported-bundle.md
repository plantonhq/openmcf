# Imported Certificate Bundle

This preset imports an existing certificate -- a wildcard or EV certificate
purchased from a public CA, delivered as a PFX (PKCS#12) or PEM bundle --
into the vault, which takes over private-key custody from that point on.

Azure derives the certificate's policy (key shape, content type) from the
bundle itself; add an explicit `certificatePolicy` alongside the import
only when you need to govern exportability or renewal actions beyond what
the bundle implies.

## When to Use

- Publicly-trusted certificates purchased outside Azure (wildcards, EV)
- Migrating certificates from another secret store into Key Vault custody
- Any certificate whose issuance flow lives outside the vault

## Key Configuration Choices

- **`contents` carries the private key** -- it is a `(sensitive)` field:
  the platform encrypts it at rest and never echoes it in outputs; supply
  it through your secrets workflow, never a committed manifest
- **Re-import = new version** -- pushing a renewed bundle into the same
  resource creates a new certificate version; consumers on
  `versionless_secret_id` follow automatically, making renewal a
  one-field update
- **Renewal stays YOUR responsibility** -- the vault notifies (add an
  `EMAIL_CONTACTS` lifetime action via an explicit policy) but cannot
  re-enroll a certificate it did not issue

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | The vault's ARM resource ID (or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<base64-pfx-or-pem-bundle>` | The bundle, base64-encoded (`base64 -i cert.pfx`) | Your CA's delivery |
| `<bundle-password>` | The bundle's password (omit the field for passwordless bundles) | Your CA's delivery |

## Downstream Wiring

Same as any vault certificate -- TLS terminators consume the secret face:

```yaml
spec:
  sslCertificates:
    - name: purchased-tls
      keyVaultSecretId:
        valueFrom:
          name: purchased-tls
```
