# CA-Signed via Pending CSR

This preset runs the out-of-band CA flow: the vault generates the private
key (which never leaves it -- `exportable: false`), mints a CSR, and holds
the operation pending until your CA's signed response is merged back. The
30-day `EMAIL_CONTACTS` action notifies the vault's certificate contacts
before expiry -- the only sensible action for `Unknown`-issuer
certificates, whose renewal is inherently manual.

## When to Use

- Publicly-trusted certificates from a CA without a first-party Key Vault
  integration (anything beyond DigiCert/GlobalSign)
- Internal PKI: your enterprise CA signs the CSR while the vault keeps
  sole custody of the private key
- Security postures that forbid private keys ever existing outside the
  vault (note `exportable: false`)

## Key Configuration Choices

- **`issuerName: Unknown`** -- the CSR-pending flow; after deployment,
  fetch the CSR from the vault, have the CA sign it, and merge the signed
  certificate back (portal or CLI)
- **`exportable: false`** -- the private key is generated inside the vault
  and can never be read out; TLS terminators must fetch through Key Vault
  references
- **TLS-server EKU** (`1.3.6.1.5.5.7.3.1`) stamped explicitly -- public
  CAs expect it on server certificates
- **`contentType: PEM`** -- the Linux/OpenSSL world's encoding; switch to
  `PKCS12` for Application Gateway and Windows consumers

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | The vault's ARM resource ID (or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<public-hostname>` | The hostname the certificate secures | Your service's public DNS name |

## Operational Notes

The resource reports created once the CSR exists; the certificate becomes
consumable after the CA's response is merged. Renewals repeat the flow --
the vault re-issues a CSR at the same policy and the contacts are emailed
30 days before expiry.
