# HSM Auto-Rotating Production CMK

This preset creates the production-hardened CMK: HSM-protected key material
(FIPS 140-2 Level 3, never leaves the hardware module) with a fully
automatic rotation policy -- Azure mints a new version every ~11 months,
stamps each version with a 2-year expiry, and raises an Event Grid
notification 30 days before any version expires.

## When to Use

- Production CMKs under compliance regimes that mandate hardware
  protection (PCI-DSS, HIPAA, FedRAMP High) or scheduled rotation
- Keys encrypting data whose exposure window must be bounded -- rotation
  limits how much data any single key version protects

## Key Configuration Choices

- **`keyType: RSA_HSM`** -- requires the vault's PREMIUM SKU; the
  software-RSA preset works on Standard vaults
- **`timeAfterCreation: P335D`** -- rotates ~30 days before the notify
  window so a fresh version is always live before anyone is paged
- **`expireAfter: P2Y`** -- old versions stay decryptable for consumers
  holding data wrapped by them, then expire; pair with consumers on
  `versionless_id` so new writes always use the current version
- **Rotation is invisible to versionless consumers** -- no manifest
  change, no redeployment

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | A PREMIUM-SKU vault's ARM resource ID (HSM types need it; or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<purpose>` | What this key encrypts (tag value) | Your tagging convention |

## Operational Notes

The vault holding a production CMK should have purge protection enabled --
many CMK integrations refuse to enroll without it, and it guarantees the
key can be recovered for the full retention window.
